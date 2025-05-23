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
	"github.com/geldata/gel-go/internal"
	"github.com/geldata/gel-go/internal/cache"
	"github.com/geldata/gel-go/internal/snc"
)

var (
	descCache = cache.New(1_000)
	rnd       = snc.NewRand()

	// ProtocolVersionMax is the highest version that the client supports.
	ProtocolVersionMax  = protocolVersion3p0
	protocolVersionMin  = protocolVersion0p13
	protocolVersion0p13 = internal.ProtocolVersion{Major: 0, Minor: 13}
	// ProtocolVersion1p0 is the version 1.0 protocol version.
	ProtocolVersion1p0 = internal.ProtocolVersion{Major: 1, Minor: 0}
	protocolVersion2p0 = internal.ProtocolVersion{Major: 2, Minor: 0}
	protocolVersion3p0 = internal.ProtocolVersion{Major: 3, Minor: 0}

	capabilitiesModifications uint64 = 0x1
	capabilitiesSessionConfig uint64 = 0x2
	capabilitiesTransaction   uint64 = 0x4
	capabilitiesDDL           uint64 = 0x8
	capabilitiesAll           uint64 = 0xffffffffffffffff

	txCapabilities   = capabilitiesAll ^ capabilitiesSessionConfig
	userCapabilities = capabilitiesAll ^
		(capabilitiesSessionConfig | capabilitiesTransaction)
)
