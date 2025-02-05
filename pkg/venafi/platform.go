/*
 * Copyright 2023 Venafi, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package venafi

import (
	"strings"

	"gopkg.in/yaml.v3"
)

type Platform int

const (
	// Undefined represents an invalid Platform
	Undefined Platform = iota
	// Fake is a fake platform for tests
	Fake
	// TLSPCloud represents the TLS Protect Cloud platform type
	TLSPCloud
	// TPP represents the TPP platform type
	TPP
	// Firefly represents the Firefly platform type
	Firefly

	// String representations of the Platform types
	strPlatformFake    = "FAKE"
	strPlatformFirefly = "FIREFLY"
	strPlatformTPP     = "TPP"
	strPlatformVaaS    = "VAAS"
	strPlatformUnknown = "Unknown"

	// alias for TPP
	strPlatformTLSPDC = "TLSPDC"
	// alias for VaaS
	strPlatformTLSPC = "TLSPC"
)

// String returns a string representation of this object
func (p Platform) String() string {
	switch p {
	case Fake:
		return strPlatformFake
	case Firefly:
		return strPlatformFirefly
	case TPP:
		return strPlatformTPP
	case TLSPCloud:
		return strPlatformVaaS
	default:
		return strPlatformUnknown
	}
}

// MarshalYAML customizes the behavior of Platform when being marshaled into a YAML document.
// The returned value is marshaled in place of the original value implementing Marshaller
func (p Platform) MarshalYAML() (interface{}, error) {
	return p.String(), nil
}

// UnmarshalYAML customizes the behavior when being unmarshalled from a YAML document
func (p *Platform) UnmarshalYAML(value *yaml.Node) error {
	var strValue string
	err := value.Decode(&strValue)
	if err != nil {
		return err
	}
	*p = GetPlatformType(strValue)
	return nil
}

func GetPlatformType(platformString string) Platform {
	switch strings.ToUpper(platformString) {
	case strPlatformFake:
		return Fake
	case strPlatformFirefly:
		return Firefly
	case strPlatformTPP, strPlatformTLSPDC:
		return TPP
	case strPlatformVaaS, strPlatformTLSPC:
		return TLSPCloud
	default:
		return Undefined
	}
}
