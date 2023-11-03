// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package ociartifactaccess

import (
	access "github.com/open-component-model/ocm/pkg/contexts/ocm/accessmethods/ociartifact"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/cpi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements/artifactaccess/genericaccess"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/resourcetypes"
)

const TYPE = resourcetypes.OCI_IMAGE

func Access[M any, P compdesc.ArtifactMetaPointer[M]](ctx cpi.Context, meta P, refname string, nooption ...any) cpi.ArtifactAccess[M] {
	if meta.GetType() == "" {
		meta.SetType(TYPE)
	}

	spec := access.New(refname)
	// is global access, must work, otherwise there is an error in the lib.
	return genericaccess.MustAccess(ctx, meta, spec)
}

func ResourceAccess(ctx cpi.Context, name string, refname string, opts ...elements.ResourceMetaOption) (cpi.ResourceAccess, error) {
	meta, err := elements.ResourceMeta(name, TYPE, opts...)
	if err != nil {
		return nil, err
	}

	return Access(ctx, meta, refname), nil
}

func SourceAccess(ctx cpi.Context, name string, path string, opts ...elements.SourceMetaOption) (cpi.SourceAccess, error) {
	meta, err := elements.SourceMeta(name, TYPE, opts...)
	if err != nil {
		return nil, err
	}

	return Access(ctx, meta, path), nil
}
