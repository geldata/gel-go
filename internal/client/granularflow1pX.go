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

	"github.com/geldata/gel-go/internal/buff"
	"github.com/geldata/gel-go/internal/codecs"
	"github.com/geldata/gel-go/internal/descriptor"
	"github.com/geldata/gel-go/internal/gelerr"
	"github.com/geldata/gel-go/internal/state"
)

func (c *protocolConnection) execGranularFlow1pX(
	r *buff.Reader,
	q *query,
) error {
	ids, ok := c.getCachedTypeIDs(q)
	if !ok {
		return c.pesimistic1pX(r, q)
	}

	cdcs, err := c.codecsFromIDs(ids, q)
	if err != nil {
		return err
	} else if cdcs == nil {
		return c.pesimistic1pX(r, q)
	}

	return c.execute1pX(r, q, cdcs)
}

func (c *protocolConnection) pesimistic1pX(r *buff.Reader, q *query) error {
	desc, err := c.parse1pX(r, q)
	if err != nil {
		return err
	}

	cdcs, err := c.codecsFromDescriptors1pX(q, desc)
	if err != nil {
		return err
	}

	return c.execute1pX(r, q, cdcs)
}

func (c *protocolConnection) parse1pX(
	r *buff.Reader,
	q *query,
) (*CommandDescription, error) {
	w := buff.NewWriter(c.writeMemory[:0])
	w.BeginMessage(uint8(Parse))
	w.PushUint16(0) // no headers
	w.PushUint64(q.getCapabilities())
	w.PushUint64(0) // no compilation_flags
	w.PushUint64(q.cfg.QueryOptions.ImplicitLimit())
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

	var desc *CommandDescription
	done := buff.NewSignal()

	for r.Next(done.Chan) {
		switch Message(r.MsgType) {
		case StateDataDescription:
			if e := c.decodeStateDataDescription(r); e != nil {
				err = wrapAll(err, e)
			}
		case CommandDataDescription:
			var e error
			desc, e = c.decodeCommandDataDescriptionMsg1pX(r, q)
			err = wrapAll(err, e)
		case ReadyForCommand:
			decodeReadyForCommandMsg(r)
			done.Signal()
		case ErrorResponse:
			err = wrapAll(err, decodeErrorResponseMsg(r, q.cmd, q.filename))
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

func (c *protocolConnection) decodeCommandDataDescriptionMsg1pX(
	r *buff.Reader,
	q *query,
) (*CommandDescription, error) {
	_, err := decodeHeaders1pX(r, q.cmd, q.filename, q.cfg.WarningHandler)
	if err != nil {
		return nil, err
	}

	c.cacheCapabilities1pX(q, r.PopUint64())

	var descs CommandDescription
	descs.Card = Cardinality(r.PopUint8())
	id := r.PopUUID()
	descs.In, err = descriptor.Pop(
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
	descs.Out, err = descriptor.Pop(
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

func (c *protocolConnection) execute1pX(
	r *buff.Reader,
	q *query,
	cdcs *codecPair,
) error {
	w := buff.NewWriter(c.writeMemory[:0])
	w.BeginMessage(uint8(Execute))
	w.PushUint16(0) // no headers
	w.PushUint64(q.getCapabilities())
	w.PushUint64(0) // no compilation_flags
	w.PushUint64(q.cfg.QueryOptions.ImplicitLimit())
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
			descs, e := c.decodeCommandDataDescriptionMsg1pX(r, q)
			err = wrapAll(err, e)
			cdcs, e = c.codecsFromDescriptors1pX(q, descs)
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
			if e := c.decodeCommandCompleteMsg1pX(q, r); e != nil {
				err = wrapAll(err, e)
			}
		case ReadyForCommand:
			decodeReadyForCommandMsg(r)
			done.Signal()
		case ErrorResponse:
			if err == ErrZeroResults {
				err = nil
			}

			err = wrapAll(err, decodeErrorResponseMsg(r, q.cmd, q.filename))
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

func (c *protocolConnection) codecsFromDescriptors1pX(
	q *query,
	descs *CommandDescription,
) (*codecPair, error) {
	var cdcs codecPair
	var err error
	cdcs.in, err = codecs.BuildEncoder(descs.In, c.protocolVersion)
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

		cdcs.out, err = codecs.BuildDecoder(descs.Out, q.outType, path)
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

func (c *protocolConnection) decodeCommandCompleteMsg1pX(
	q *query,
	r *buff.Reader,
) error {
	discardHeaders1pX(r)
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

func (c *protocolConnection) decodeStateDataDescription(r *buff.Reader) error {
	if c.protocolVersion.GTE(protocolVersion2p0) {
		return c.decodeStateDataDescription2pX(r)
	}

	id := r.PopUUID()
	desc, err := descriptor.Pop(
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

	codec, err := state.BuildEncoder(desc, codecs.Path("state"))
	if err != nil {
		return gelerr.NewBinaryProtocolError("", fmt.Errorf(
			"building decoder from ParameterStatus state_description: %w",
			err))
	}

	c.stateCodec = codec
	return nil
}

func (c *protocolConnection) codecsFromIDs(
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

		in, err = codecs.BuildEncoder(
			desc.(descriptor.Descriptor),
			c.protocolVersion,
		)
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

		d := desc.(descriptor.Descriptor)
		path := codecs.Path(q.outType.String())
		out, err = codecs.BuildDecoder(d, q.outType, path)
		if err != nil {
			return nil, gelerr.NewInvalidArgumentError(fmt.Sprintf(
				"the \"out\" argument does not match query schema: %v",
				err,
			), nil)
		}
	}

	return &codecPair{in: in.(codecs.Encoder), out: out.(codecs.Decoder)}, nil
}
