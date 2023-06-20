// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package attr

import (
	"fmt"

	"github.com/open-component-model/ocm/pkg/contexts/datacontext"
	"github.com/open-component-model/ocm/pkg/errors"
	"github.com/open-component-model/ocm/pkg/runtime"
)

const (
	ATTR_KEY   = "ocm.software/signing/sigstore"
	ATTR_SHORT = "sigstore"
)

func init() {
	datacontext.RegisterAttributeType(ATTR_KEY, AttributeType{}, ATTR_SHORT)
}

type AttributeType struct{}

func (a AttributeType) Name() string {
	return ATTR_KEY
}

func (a AttributeType) Description() string {
	return `
*sigstore config* Configuration to use for sigstore based sogning.
The following fields are used.
- *<code>fulcioURL</code>* *string*  default is https://v1.fulcio.sigstore.dev
- *<code>OIDCIssuer</code>* *string*  default is https://oauth2.sigstore.dev/auth
- *<code>OIDCClientID</code>* *string*  default is sigstore
`
}

func (a AttributeType) Encode(v interface{}, marshaller runtime.Marshaler) ([]byte, error) {
	if marshaller == nil {
		marshaller = runtime.DefaultJSONEncoding
	}
	if a, ok := v.(*Attribute); !ok {
		return nil, fmt.Errorf("sigstore attribute")
	} else {
		return marshaller.Marshal(a)
	}
}

func (a AttributeType) Decode(data []byte, unmarshaller runtime.Unmarshaler) (interface{}, error) {
	if unmarshaller == nil {
		unmarshaller = runtime.DefaultJSONEncoding
	}
	var attr Attribute
	err := unmarshaller.Unmarshal(data, &attr)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid attribute value for %s", ATTR_KEY)
	}
	return &attr, nil
}

////////////////////////////////////////////////////////////////////////////////

type Attribute struct {
	FulcioURL    string `json:"fulcioURL"`
	OIDCIssuer   string `json:"OIDCIssuer"`
	OIDCClientID string `json:"OIDCClientID"`
}

////////////////////////////////////////////////////////////////////////////////

func Get(ctx datacontext.Context) *Attribute {
	v := ctx.GetAttributes().GetAttribute(ATTR_KEY)

	if v != nil {
		a, _ := v.(*Attribute)
		return a
	}
	return &Attribute{
		FulcioURL:    "https://v1.fulcio.sigstore.dev",
		OIDCIssuer:   "https://oauth2.sigstore.dev/auth",
		OIDCClientID: "sigstore",
	}
}

func Set(ctx datacontext.Context, a *Attribute) {
	ctx.GetAttributes().SetAttribute(ATTR_KEY, a)
}
