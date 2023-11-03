// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package genericblob

import (
	"github.com/open-component-model/ocm/pkg/blobaccess"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/cpi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements/artifactaccess/epi"
	"github.com/open-component-model/ocm/pkg/generics"
	"github.com/open-component-model/ocm/pkg/optionutils"
)

func Access[M any, P compdesc.ArtifactMetaPointer[M]](ctx cpi.Context, meta P, blob blobaccess.BlobAccessProvider, opts ...Option) cpi.ArtifactAccess[M] {
	eff := optionutils.EvalOptions(opts...)
	accprov := cpi.NewAccessProviderForBlobAccessProvider(ctx, blob, eff.Hint, eff.Global)
	// strange type cast is required by Go compiler, meta has the correct type.
	return cpi.NewArtifactAccessForProvider(generics.As[*M](meta), accprov)
}

func ResourceAccess(ctx cpi.Context, name, typ string, opts ...elements.ResourceMetaOption) func(blob blobaccess.BlobAccessProvider, opts ...Option) (cpi.ResourceAccess, error) {
	return epi.ResourceAccessA[blobaccess.BlobAccessProvider, Option](ctx, name, typ, Access[compdesc.ResourceMeta], opts...)
}

func SourceAccess(ctx cpi.Context, name, typ string, opts ...elements.SourceMetaOption) func(blob blobaccess.BlobAccessProvider, opts ...Option) (cpi.SourceAccess, error) {
	return epi.SourceAccessA[blobaccess.BlobAccessProvider, Option](ctx, name, typ, Access[compdesc.SourceMeta], opts...)
}
