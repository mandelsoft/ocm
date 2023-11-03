// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package helmaccess

import (
	access "github.com/open-component-model/ocm/pkg/contexts/ocm/accessmethods/helm"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/cpi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements/artifactaccess/genericaccess"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/resourcetypes"
)

const TYPE = resourcetypes.HELM_CHART

func Access[M any, P compdesc.ArtifactMetaPointer[M]](ctx cpi.Context, meta P, chart, repourl string) cpi.ArtifactAccess[M] {
	if meta.GetType() == "" {
		meta.SetType(TYPE)
	}

	spec := access.New(chart, repourl)
	// is global access, must work, otherwise there is an error in the lib.
	return genericaccess.MustAccess(ctx, meta, spec)
}

func ResourceAccess(ctx cpi.Context, name string, chart, repourl string, opts ...elements.ResourceMetaOption) (cpi.ResourceAccess, error) {
	meta, err := elements.ResourceMeta(name, TYPE, opts...)
	if err != nil {
		return nil, err
	}

	return Access(ctx, meta, chart, repourl), nil
}

func SourceAccess(ctx cpi.Context, name string, chart, repourl string, opts ...elements.SourceMetaOption) (cpi.SourceAccess, error) {
	meta, err := elements.SourceMeta(name, TYPE, opts...)
	if err != nil {
		return nil, err
	}

	return Access(ctx, meta, chart, repourl), nil
}
