// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package accessmethods

import (
	out "fmt"
	"io"
	"os"
	"strings"

	"github.com/mandelsoft/filepath/pkg/filepath"

	"github.com/open-component-model/ocm/local/r3trans/plugin/common"
	"github.com/open-component-model/ocm/local/r3trans/plugin/config"
	"github.com/open-component-model/ocm/pkg/cobrautils/flagsets"
	"github.com/open-component-model/ocm/pkg/contexts/credentials"
	"github.com/open-component-model/ocm/pkg/contexts/oci/identity"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/accessmethods/options"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/plugin/ppi"
	"github.com/open-component-model/ocm/pkg/errors"
	"github.com/open-component-model/ocm/pkg/runtime"
)

const (
	NAME    = "r3trans.sap.com"
	VERSION = "v1"
)

type AccessSpec struct {
	runtime.ObjectVersionedType `json:",inline"`

	Transport       string `json:"transport"`
	TransportSystem string `json:"transportSystem"`
	Path            string `json:"path,omitempty"`
}

const OPT_TRANSPORT = "accessTransport"
const OPT_TRANSPORT_PATH = "accessTransportPath"
const OPT_TRANSPORT_SYSTEM = "accessTransportSystem"

type AccessMethod struct {
	ppi.AccessMethodBase
}

var TransportOption = options.NewStringOptionType(OPT_TRANSPORT, "name of R/3 transport request")
var SystemOption = options.NewStringOptionType(OPT_TRANSPORT_SYSTEM, "R/3 transport system")
var PathOption = options.NewStringOptionType(OPT_TRANSPORT_PATH, "path in transport system")

var _ ppi.AccessMethod = (*AccessMethod)(nil)

func New() ppi.AccessMethod {
	return &AccessMethod{
		AccessMethodBase: ppi.MustNewAccessMethodBase(NAME, "", &AccessSpec{}, "demo access R/3 transport files", `
The type specific specification fields are:

- **<code>transport</code>** *string*

  name of transport request.

- **<code>transportSystem</code>** *string*

  address of transport system

- **<code>path</code>** *string*

  sub path in transport system.
`),
	}
}

func (a *AccessMethod) Options() []options.OptionType {
	return []options.OptionType{
		SystemOption,
		TransportOption,
		PathOption,
	}
}

func (a *AccessMethod) Decode(data []byte, unmarshaler runtime.Unmarshaler) (runtime.TypedObject, error) {
	if unmarshaler == nil {
		unmarshaler = runtime.DefaultYAMLEncoding
	}
	var spec AccessSpec
	err := unmarshaler.Unmarshal(data, &spec)
	if err != nil {
		return nil, err
	}
	return &spec, nil
}

func (a *AccessMethod) ValidateSpecification(p ppi.Plugin, spec ppi.AccessSpec) (*ppi.AccessSpecInfo, error) {
	var info ppi.AccessSpecInfo

	my := spec.(*AccessSpec)

	if my.Transport == "" {
		return nil, out.Errorf("path not specified")
	}
	if strings.Contains(my.Transport, "/") {
		return nil, out.Errorf("invalid transport request name (%s)", my.Transport)
	}
	if my.TransportSystem == "" {
		return nil, out.Errorf("address of transport system not specified")
	}
	if strings.HasPrefix(my.Transport, "/") {
		return nil, out.Errorf("invalid transport system address (%s)", my.TransportSystem)
	}
	if strings.HasPrefix(my.Path, "/") {
		return nil, out.Errorf("path must be relative (%s)", my.Path)
	}
	info.MediaType = "application/vnd.sap.com.r3trans.v1"
	info.ConsumerId = credentials.ConsumerIdentity{
		identity.ID_TYPE:       common.CONSUMER_TYPE,
		identity.ID_HOSTNAME:   my.TransportSystem,
		identity.ID_PATHPREFIX: my.Path,
	}

	info.Short = out.Sprintf("R/3 transport request %s[%s:]", my.Transport, my.TransportSystem, my.Path)
	info.Hint = my.Transport
	return &info, nil
}

func (a *AccessMethod) ComposeAccessSpecification(p ppi.Plugin, opts ppi.Config, config ppi.Config) error {
	list := errors.ErrListf("configuring options")
	list.Add(flagsets.AddFieldByOptionP(opts, TransportOption, config, "transport"))
	list.Add(flagsets.AddFieldByOptionP(opts, PathOption, config, "path"))
	list.Add(flagsets.AddFieldByOptionP(opts, SystemOption, config, "transportSystem"))
	return list.Result()
}

func (a *AccessMethod) Reader(p ppi.Plugin, spec ppi.AccessSpec, creds credentials.Credentials) (io.ReadCloser, error) {
	my := spec.(*AccessSpec)

	root, _, err := config.BasePath(p, my.TransportSystem)
	if err != nil {
		return nil, err
	}

	if my.Path != "" {
		root = filepath.Join(root, my.Path)
	}

	return os.Open(filepath.Join(root, my.Transport))
}
