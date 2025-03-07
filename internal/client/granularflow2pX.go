// This source file is part of the Gel open source project.
//
// Copyright Gel Data Inc. and the Gel authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gel

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/geldata/gel-go/internal/buff"
	"github.com/geldata/gel-go/internal/codecs"
	"github.com/geldata/gel-go/internal/descriptor"
	"github.com/geldata/gel-go/internal/gelerr"
	"github.com/geldata/gel-go/internal/state"
)

func (c *protocolConnection) execGranularFlow2pX(
	r *buff.Reader,
	q *query,
) error {
	var cdcs *codecPair
	if q.parse {
		ids, ok := c.getCachedTypeIDs(q)
		if !ok {
			return c.pesimistic2pX(r, q)
		}

		var err error
		cdcs, err = c.codecsFromIDsV2(ids, q)
		if err != nil {
			return err
		} else if cdcs == nil {
			return c.pesimistic2pX(r, q)
		}
	} else {
		cdcs = &codecPair{in: codecs.NoOpEncoder, out: codecs.NoOpDecoder}
	}

	return c.execute2pX(r, q, cdcs)
}

func (c *protocolConnection) pesimistic2pX(r *buff.Reader, q *query) error {
	desc, err := c.parse2pX(r, q)
	if err != nil {
		return err
	}

	cdcs, err := c.codecsFromDescriptors2pX(q, desc)
	if err != nil {
		return err
	}

	return c.execute2pX(r, q, cdcs)
}

func (c *protocolConnection) writeAnnotations(w *buff.Writer, q *query) error {
	n := len(q.cfg.Annotations)
	if c.protocolVersion.GTE(protocolVersion3p0) && n > 0 {
		if n >= 1<<16 {
			return gelerr.NewInvalidArgumentError("too many annotations", nil)
		}
		w.PushUint16(uint16(n))
		for key, val := range q.cfg.Annotations {
			w.PushString(key)
			w.PushString(val)
		}
	} else {
		w.PushUint16(0)
	}

	return nil
}

func (c *protocolConnection) parse2pX(
	r *buff.Reader,
	q *query,
) (*CommandDescriptionV2, error) {
	w := buff.NewWriter(c.writeMemory[:0])
	w.BeginMessage(uint8(Parse))
	if err := c.writeAnnotations(w, q); err != nil {
		return nil, err
	}
	w.PushUint64(q.getCapabilities())
	w.PushUint64(0) // no compilation_flags
	w.PushUint64(q.cfg.QueryOptions.ImplicitLimit())
	if c.protocolVersion.GTE(protocolVersion3p0) {
		w.PushUint8(uint8(q.lang))
	}
	w.PushUint8(uint8(q.fmt))
	w.PushUint8(uint8(q.expCard))
	w.PushString(q.cmd)

	w.PushUUID(c.stateCodec.DescriptorID())
	err := c.stateCodec.Encode(w, q.state, codecs.Path("state"), false)
	if err != nil {
		return nil, gelerr.NewBinaryProtocolError("", fmt.Errorf(
			"invalid connection state: %w", err))
	}
	w.EndMessage()

	w.BeginMessage(uint8(Sync))
	w.EndMessage()

	if e := c.soc.WriteAll(w.Unwrap()); e != nil {
		return nil, gelerr.NewClientConnectionClosedError("", e)
	}

	var desc *CommandDescriptionV2
	done := buff.NewSignal()

	for r.Next(done.Chan) {
		switch Message(r.MsgType) {
		case StateDataDescription:
			if e := c.decodeStateDataDescription(r); e != nil {
				err = wrapAll(err, e)
			}
		case CommandDataDescription:
			var e error
			desc, e = c.decodeCommandDataDescriptionMsg2pX(r, q)
			err = wrapAll(err, e)
		case ReadyForCommand:
			decodeReadyForCommandMsg(r)
			done.Signal()
		case ErrorResponse:
			err = wrapAll(err, decodeErrorResponseMsg(r, q.cmd))
		default:
			if e := c.fallThrough(r); e != nil {
				// the connection will not be usable after this x_x
				return nil, e
			}
		}
	}

	if r.Err != nil || err != nil {
		return nil, wrapAll(r.Err, err)
	}

	return desc, nil
}

func (c *protocolConnection) decodeCommandDataDescriptionMsg2pX(
	r *buff.Reader,
	q *query,
) (*CommandDescriptionV2, error) {
	_, err := decodeHeaders2pX(r, q.cmd, q.cfg.WarningHandler)
	if err != nil {
		return nil, err
	}

	c.cacheCapabilities1pX(q, r.PopUint64())

	var descs CommandDescriptionV2
	descs.Card = Cardinality(r.PopUint8())
	id := r.PopUUID()
	descs.In, err = descriptor.PopV2(
		r.PopSlice(r.PopUint32()),
		c.protocolVersion,
	)
	if err != nil {
		return nil, err
	} else if descs.In.ID != id {
		return nil, gelerr.NewClientError(fmt.Sprintf(
			"unexpected in descriptor id: %v",
			descs.In.ID,
		), nil)
	}

	id = r.PopUUID()
	descs.Out, err = descriptor.PopV2(
		r.PopSlice(r.PopUint32()),
		c.protocolVersion,
	)
	if err != nil {
		return nil, err
	} else if descs.Out.ID != id {
		return nil, gelerr.NewClientError(fmt.Sprintf(
			"unexpected out descriptor id: got %v but expected %v",
			descs.Out.ID,
			id,
		), nil)
	}

	if q.expCard == AtMostOne && descs.Card == Many {
		return nil, gelerr.NewResultCardinalityMismatchError(fmt.Sprintf(
			"the query has cardinality %v "+
				"which does not match the expected cardinality %v",
			descs.Card,
			q.expCard), nil)
	}

	c.cacheTypeIDs(q, idPair{in: descs.In.ID, out: descs.Out.ID})
	descCache.Put(descs.In.ID, descs.In)
	descCache.Put(descs.Out.ID, descs.Out)
	return &descs, nil
}

func (c *protocolConnection) execute2pX(
	r *buff.Reader,
	q *query,
	cdcs *codecPair,
) error {
	w := buff.NewWriter(c.writeMemory[:0])
	w.BeginMessage(uint8(Execute))
	if err := c.writeAnnotations(w, q); err != nil {
		return err
	}
	w.PushUint64(q.getCapabilities())
	w.PushUint64(0) // no compilation_flags
	w.PushUint64(q.cfg.QueryOptions.ImplicitLimit())
	if c.protocolVersion.GTE(protocolVersion3p0) {
		w.PushUint8(uint8(q.lang))
	}
	w.PushUint8(uint8(q.fmt))
	w.PushUint8(uint8(q.expCard))
	w.PushString(q.cmd)
	w.PushUUID(c.stateCodec.DescriptorID())
	err := c.stateCodec.Encode(w, q.state, codecs.Path("state"), false)
	if err != nil {
		return gelerr.NewBinaryProtocolError("", fmt.Errorf(
			"invalid connection state: %w", err))
	}

	w.PushUUID(cdcs.in.DescriptorID())
	w.PushUUID(cdcs.out.DescriptorID())
	if e := cdcs.in.Encode(w, q.args, codecs.Path("args"), true); e != nil {
		return gelerr.NewInvalidArgumentError(e.Error(), nil)
	}
	w.EndMessage()

	w.BeginMessage(uint8(Sync))
	w.EndMessage()

	if e := c.soc.WriteAll(w.Unwrap()); e != nil {
		return gelerr.NewClientConnectionClosedError("", e)
	}

	tmp := q.out
	if q.expCard == AtMostOne {
		err = ErrZeroResults
	}
	done := buff.NewSignal()

	for r.Next(done.Chan) {
		switch Message(r.MsgType) {
		case StateDataDescription:
			if e := c.decodeStateDataDescription(r); e != nil {
				err = wrapAll(err, e)
			}
		case CommandDataDescription:
			descs, e := c.decodeCommandDataDescriptionMsg2pX(r, q)
			err = wrapAll(err, e)
			cdcs, e = c.codecsFromDescriptors2pX(q, descs)
			err = wrapAll(err, e)
		case Data:
			val, ok, e := decodeDataMsg(r, q, cdcs)
			if e != nil {
				if err == ErrZeroResults {
					err = e
				} else {
					err = wrapAll(err, e)
				}
			}
			if ok {
				tmp = reflect.Append(tmp, val)
			}

			if err == ErrZeroResults {
				err = nil
			}
		case CommandComplete:
			if e := c.decodeCommandCompleteMsg2pX(q, r); e != nil {
				err = wrapAll(err, e)
			}
		case ReadyForCommand:
			decodeReadyForCommandMsg(r)
			done.Signal()
		case ErrorResponse:
			if err == ErrZeroResults {
				err = nil
			}

			err = wrapAll(err, decodeErrorResponseMsg(r, q.cmd))
		default:
			if e := c.fallThrough(r); e != nil {
				// the connection will not be usable after this x_x
				return e
			}
		}
	}

	if r.Err != nil {
		return wrapAll(err, r.Err)
	}

	if !q.flat() && q.fmt != Null {
		q.out.Set(tmp)
	}

	return err
}

func (c *protocolConnection) codecsFromIDsV2(
	ids *idPair,
	q *query,
) (*codecPair, error) {
	var err error

	in, ok := c.inCodecCache.Get(ids.in)
	if !ok {
		desc, OK := descCache.Get(ids.in)
		if !OK {
			return nil, nil
		}

		d := desc.(descriptor.V2)
		in, err = codecs.BuildEncoderV2(&d, c.protocolVersion)
		if err != nil {
			return nil, gelerr.NewInvalidArgumentError(err.Error(), nil)
		}
	}

	out, ok := c.outCodecCache.Get(codecKey{ID: ids.out, Type: q.outType})
	if !ok {
		desc, OK := descCache.Get(ids.out)
		if !OK {
			return nil, nil
		}

		d := desc.(descriptor.V2)
		path := codecs.Path(q.outType.String())
		out, err = codecs.BuildDecoderV2(&d, q.outType, path)
		if err != nil {
			return nil, gelerr.NewInvalidArgumentError(fmt.Sprintf(
				"the \"out\" argument does not match query schema: %v",
				err,
			), nil)
		}
	}

	return &codecPair{in: in.(codecs.Encoder), out: out.(codecs.Decoder)}, nil
}

func (c *protocolConnection) codecsFromDescriptors2pX(
	q *query,
	descs *CommandDescriptionV2,
) (*codecPair, error) {
	var cdcs codecPair
	var err error
	cdcs.in, err = codecs.BuildEncoderV2(&descs.In, c.protocolVersion)
	if err != nil {
		return nil, gelerr.NewInvalidArgumentError(err.Error(), nil)
	}

	if q.fmt == JSON {
		cdcs.out = codecs.JSONBytes
	} else {
		var path codecs.Path
		if q.fmt == Null {
			// There is no outType value for Null output format queries.
			path = "null"
		} else {
			path = codecs.Path(q.outType.String())
		}

		cdcs.out, err = codecs.BuildDecoderV2(&descs.Out, q.outType, path)
		if err != nil {
			err = fmt.Errorf(
				"the \"out\" argument does not match query schema: %v",
				err,
			)
			return nil, gelerr.NewInvalidArgumentError(err.Error(), nil)
		}
	}

	c.inCodecCache.Put(cdcs.in.DescriptorID(), cdcs.in)
	c.outCodecCache.Put(
		codecKey{ID: cdcs.out.DescriptorID(), Type: q.outType},
		cdcs.out,
	)

	return &cdcs, nil
}

func (c *protocolConnection) decodeCommandCompleteMsg2pX(
	q *query,
	r *buff.Reader,
) error {
	discardHeaders2pX(r)
	c.cacheCapabilities1pX(q, r.PopUint64())
	r.Discard(int(r.PopUint32())) // discard command status
	if r.PopUUID() == descriptor.IDZero {
		// empty state data
		r.Discard(4)
		return nil
	}

	r.Discard(int(r.PopUint32())) // state data
	return nil
}

func (c *protocolConnection) decodeStateDataDescription2pX(
	r *buff.Reader,
) error {
	id := r.PopUUID()
	desc, err := descriptor.PopV2(
		r.PopSlice(r.PopUint32()),
		c.protocolVersion,
	)
	if err != nil {
		return gelerr.NewBinaryProtocolError("", fmt.Errorf(
			"decoding ParameterStatus state_description: %w", err))
	} else if desc.ID != id {
		return gelerr.NewBinaryProtocolError("", fmt.Errorf(
			"state_description ids don't match: %v != %v", id, desc.ID))
	}

	codec, err := state.BuildEncoderV2(&desc, codecs.Path("state"))
	if err != nil {
		return gelerr.NewBinaryProtocolError("", fmt.Errorf(
			"building decoder from ParameterStatus state_description: %w",
			err))
	}

	c.stateCodec = codec
	return nil
}

func decodeReadyForCommandMsg(r *buff.Reader) {
	ignoreHeaders(r)
	r.Discard(1) // transaction state
}

func decodeDataMsg(
	r *buff.Reader,
	q *query,
	cdcs *codecPair,
) (reflect.Value, bool, error) {
	elmCount := r.PopUint16()
	if elmCount != 1 {
		return reflect.Value{}, false, fmt.Errorf(
			"unexpected number of elements: expected 1, got %v", elmCount)
	}
	elmLen := r.PopUint32()

	if !q.flat() {
		val := reflect.New(q.outType).Elem()
		err := cdcs.out.Decode(
			r.PopSlice(elmLen),
			unsafe.Pointer(val.UnsafeAddr()),
		)
		if err != nil {
			return reflect.Value{}, false, err
		}
		return val, true, nil
	}

	err := cdcs.out.Decode(
		r.PopSlice(elmLen),
		unsafe.Pointer(q.out.UnsafeAddr()),
	)
	if err != nil {
		return reflect.Value{}, false, err
	}

	return reflect.Value{}, false, nil
}
