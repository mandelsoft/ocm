// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package github

import (
	access "github.com/open-component-model/ocm/pkg/contexts/ocm/accessmethods/s3"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/cpi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements/artifactaccess/epi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements/artifactaccess/genericaccess"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/resourcetypes"
	"github.com/open-component-model/ocm/pkg/mime"
	"github.com/open-component-model/ocm/pkg/optionutils"
)

const TYPE = resourcetypes.BLOB

func Access[M any, P compdesc.ArtifactMetaPointer[M]](ctx cpi.Context, meta P, bucket, key string, opts ...Option) cpi.ArtifactAccess[M] {
	eff := optionutils.EvalOptions(opts...)
	if meta.GetType() == "" {
		meta.SetType(TYPE)
	}

	media := eff.MediaType
	if media == "" {
		media = mime.MIME_OCTET
	}
	spec := access.New(eff.Region, bucket, key, eff.Version, media)
	// is global access, must work, otherwise there is an error in the lib.
	return genericaccess.MustAccess(ctx, meta, spec)
}

func ResourceAccess(ctx cpi.Context, name string, opts ...elements.ResourceMetaOption) func(bucket, key string, opts ...Option) (cpi.ResourceAccess, error) {
	return epi.ResourceAccessAB[string, string, Option](ctx, name, TYPE, Access[compdesc.ResourceMeta], opts...)
}

func SourceAccess(ctx cpi.Context, name string, opts ...elements.SourceMetaOption) func(bucket, key string, opts ...Option) (cpi.SourceAccess, error) {
	return epi.SourceAccessAB[string, string, Option](ctx, name, TYPE, Access[compdesc.SourceMeta], opts...)
}
