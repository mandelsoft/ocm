// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package dockermultiblob

import (
	"github.com/open-component-model/ocm/pkg/blobaccess/dockermulti"
	"github.com/open-component-model/ocm/pkg/contexts/ocm"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/cpi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements/artifactaccess/epi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/resourcetypes"
	"github.com/open-component-model/ocm/pkg/generics"
	"github.com/open-component-model/ocm/pkg/optionutils"
)

const TYPE = resourcetypes.OCI_IMAGE

func Access[M any, P compdesc.ArtifactMetaPointer[M]](ctx ocm.Context, meta P, opts ...Option) cpi.ArtifactAccess[M] {
	eff := optionutils.EvalOptions(opts...)
	if meta.GetType() == "" {
		meta.SetType(TYPE)
	}
	eff.Blob.Context = ctx.OCIContext()

	blobprov := dockermulti.BlobAccessProviderForMultiImageFromDockerDaemon(&eff.Blob)
	accprov := cpi.NewAccessProviderForBlobAccessProvider(ctx, blobprov, eff.Hint, eff.Global)
	// strange type cast is required by Go compiler, meta has the correct type.
	return cpi.NewArtifactAccessForProvider(generics.As[*M](meta), accprov)
}

func ResourceAccess(ctx cpi.Context, name string, opts ...elements.ResourceMetaOption) func(pts ...Option) (cpi.ResourceAccess, error) {
	return epi.ResourceAccess[Option](ctx, name, TYPE, Access[compdesc.ResourceMeta], opts...)
}

func SourceAccess(ctx cpi.Context, name string, opts ...elements.SourceMetaOption) func(opts ...Option) (cpi.SourceAccess, error) {
	return epi.SourceAccess[Option](ctx, name, TYPE, Access[compdesc.SourceMeta], opts...)
}
