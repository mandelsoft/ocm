// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package uploaders

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/open-component-model/ocm/local/r3trans/plugin/accessmethods"
	"github.com/open-component-model/ocm/pkg/common/accessio"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/plugin/ppi"
	"github.com/open-component-model/ocm/pkg/errors"
	"github.com/open-component-model/ocm/pkg/runtime"
)

type writer = accessio.DigestWriter

type Writer struct {
	*writer
	file    *os.File
	acc     string
	version string
	root    string
	name    string
	tpath   string
	sys     string
	spec    *accessmethods.AccessSpec
}

func NewWriter(file *os.File, root, name, tpath, sys string, acc, version string) *Writer {
	return &Writer{
		writer:  accessio.NewDefaultDigestWriter(file),
		file:    file,
		acc:     acc,
		version: version,
		root:    root,
		name:    name,
		tpath:   tpath,
		sys:     sys,
	}
}

func (w *Writer) Close() error {
	err := w.writer.Close()
	if err == nil {
		name := w.name
		if name == "" {
			dig := string(w.writer.Digest())
			if i := strings.Index(dig, ":"); i >= 0 {
				dig = dig[i+1:]
			}
			name = strings.ToUpper(dig[:8])
			if len(w.tpath) >= 3 {
				name = strings.ToUpper(w.tpath[:3]) + "K" + name
			} else {
				name = "TMPK" + name
			}
			n := filepath.Join(w.root, name)
			err := os.Rename(w.file.Name(), n)
			if err != nil {
				return errors.Wrapf(err, "cannot rename %q to %q", w.file.Name(), n)
			}
		}
		w.spec = &accessmethods.AccessSpec{
			ObjectVersionedType: runtime.NewVersionedObjectType(w.acc, w.version),
			Transport:           name,
			Path:                w.tpath,
			TransportSystem:     w.sys,
		}
	}
	return err
}

func (w *Writer) Specification() ppi.AccessSpec {
	return w.spec
}
