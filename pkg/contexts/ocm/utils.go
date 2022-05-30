// Copyright 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package ocm

import (
	"fmt"
	"io"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/open-component-model/ocm/pkg/common"
	"github.com/open-component-model/ocm/pkg/common/accessio"
	"github.com/open-component-model/ocm/pkg/common/accessobj"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/core"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/cpi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/repositories/ctf"
	"github.com/open-component-model/ocm/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////

func AssureTargetRepository(session Session, ctx Context, targetref string, opts ...interface{}) (Repository, error) {
	var format accessio.FileFormat
	var archive string
	var fs vfs.FileSystem

	for _, o := range opts {
		switch v := o.(type) {
		case vfs.FileSystem:
			if fs == nil && v != nil {
				fs = v
			}
		case accessio.FileFormat:
			format = v
		case string:
			archive = v
		default:
			panic(fmt.Sprintf("invalid option type %T", o))
		}
	}

	ref, err := ParseRepo(targetref)
	if err != nil {
		return nil, err
	}
	if archive != "" && ref.Type != "" {
		for _, f := range ctf.SupportedFormats() {
			if f.String() == ref.Type {
				ref.Type = archive + "+" + ref.Type
			}
		}
	}
	ref.TypeHint = archive
	ref.CreateIfMissing = true
	target, err := session.DetermineRepositoryBySpec(ctx, &ref)
	if err != nil {
		if !errors.IsErrUnknown(err) || vfs.IsErrNotExist(err) || ref.Info == "" {
			return nil, err
		}
		if ref.Type == "" {
			ref.Type = format.String()
		}
		if ref.Type == "" {
			return nil, fmt.Errorf("ctf format type required to create ctf")
		}
		target, err = ctf.Create(ctx, accessobj.ACC_CREATE, ref.Info, 0770, accessio.PathFileSystem(accessio.FileSystem(fs)))
		if err != nil {
			return nil, err
		}
		session.Closer(target)
	}
	return target, nil
}

type AccessMethodSource = cpi.AccessMethodSource

// ResourceReader gets a Reader for a given resource/source access.
// It provides a Reader handling the Close contract for the access method
// by connecting the access method's Close method to the Readers Close method .
func ResourceReader(s AccessMethodSource) (io.ReadCloser, error) {
	return cpi.ResourceReader(s)
}

// ResourceData extracts the data for a given resource/source access.
// It handles the Close contract for the access method for a singular use.
func ResourceData(s AccessMethodSource) ([]byte, error) {
	return cpi.ResourceData(s)
}

type CompoundResolver struct {
	resolvers []ComponentVersionResolver
}

var _ ComponentVersionResolver = (*CompoundResolver)(nil)

func NewCompoundResolver(res ...ComponentVersionResolver) *CompoundResolver {
	return &CompoundResolver{res}
}

func (c *CompoundResolver) LookupComponentVersion(name string, version string) (core.ComponentVersionAccess, error) {
	for _, r := range c.resolvers {
		cv, err := r.LookupComponentVersion(name, version)
		if err == nil {
			return cv, nil
		}
		if !errors.IsErrNotFoundKind(err, KIND_COMPONENTVERSION) {
			return nil, err
		}
	}
	return nil, errors.ErrNotFound(KIND_OCM_REFERENCE, common.NewNameVersion(name, version).String())
}
