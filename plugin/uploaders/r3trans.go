// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package uploaders

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-component-model/ocm/local/r3trans/plugin/accessmethods"
	"github.com/open-component-model/ocm/local/r3trans/plugin/common"
	"github.com/open-component-model/ocm/local/r3trans/plugin/config"
	"github.com/open-component-model/ocm/pkg/contexts/credentials"
	"github.com/open-component-model/ocm/pkg/contexts/oci/identity"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/plugin/ppi"
	"github.com/open-component-model/ocm/pkg/runtime"
)

const (
	NAME    = "r3trans"
	VERSION = "v1"
)

type TargetSpec struct {
	runtime.ObjectVersionedType `json:",inline"`

	TransportSystem string `json:"transportSystem"`
	Path            string `json:"path"`
}

var types map[string]runtime.TypedObjectDecoder

func init() {
	decoder, err := runtime.NewDirectDecoder(&TargetSpec{})
	if err != nil {
		panic(err)
	}
	types = map[string]runtime.TypedObjectDecoder{NAME + runtime.VersionSeparator + VERSION: decoder}
}

type Uploader struct {
	ppi.UploaderBase
}

var _ ppi.Uploader = (*Uploader)(nil)

func New() ppi.Uploader {
	return &Uploader{
		UploaderBase: ppi.MustNewUploaderBase(NAME, `
Upload R/3 transport requests to the transport system.

It uses the following target specification fields:

- **<code>type</code** *string* constant <code>r3trans/v1</code>

- **<code>transportSystem</code>** *string*

  The address of the R/3 transport system.

- **<code>path</code>** *string*

  The sub path used in the transport system.
`,
		),
	}
}

func (a *Uploader) Decoders() map[string]runtime.TypedObjectDecoder {
	return types
}

func (a *Uploader) ValidateSpecification(p ppi.Plugin, spec ppi.UploadTargetSpec) (*ppi.UploadTargetSpecInfo, error) {
	var info ppi.UploadTargetSpecInfo
	my := spec.(*TargetSpec)

	if strings.Contains(my.TransportSystem, "/") {
		return nil, fmt.Errorf("invalid transport system (%s)", my.TransportSystem)
	}

	if strings.HasPrefix(my.Path, "/") {
		return nil, fmt.Errorf("path must be relative (%s)", my.Path)
	}

	info.ConsumerId = credentials.ConsumerIdentity{
		identity.ID_TYPE:       common.CONSUMER_TYPE,
		identity.ID_HOSTNAME:   my.TransportSystem,
		identity.ID_PATHPREFIX: my.Path,
	}
	return &info, nil
}

func (a *Uploader) Writer(p ppi.Plugin, arttype, mediatype, hint string, repo ppi.UploadTargetSpec, creds credentials.Credentials) (io.WriteCloser, ppi.AccessSpecProvider, error) {
	var file *os.File
	var err error

	root, def, err := config.BasePath(p, repo.(*TargetSpec).TransportSystem)
	if err != nil {
		return nil, nil, err
	}

	my := repo.(*TargetSpec)
	if my.Path != "" {
		def = my.Path
	}
	if def != "" {
		root = filepath.Join(root, def)
	}

	err = os.MkdirAll(root, 0o700)
	if err != nil {
		return nil, nil, err
	}

	if hint == "" {
		file, err = os.CreateTemp(root, "r3trans.*.blob")
	} else {
		file, err = os.OpenFile(filepath.Join(root, hint), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	}
	if err != nil {
		return nil, nil, err
	}
	writer := NewWriter(file, root, hint, def, my.TransportSystem, accessmethods.NAME, accessmethods.VERSION)
	return writer, writer.Specification, nil
}
