// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package fileblob

import (
	"github.com/open-component-model/ocm/pkg/blobaccess"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/cpi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements/artifactaccess/epi"
	"github.com/open-component-model/ocm/pkg/generics"
	"github.com/open-component-model/ocm/pkg/mime"
	"github.com/open-component-model/ocm/pkg/optionutils"
)

const TYPE = "blob"

func Access[M any, P compdesc.ArtifactMetaPointer[M]](ctx cpi.Context, meta P, media, path string, opts ...Option) cpi.ArtifactAccess[M] {
	eff := optionutils.EvalOptions(opts...)

	if meta.GetType() == "" {
		meta.SetType(TYPE)
	}
	if media == "" {
		media = mime.MIME_OCTET
	}

	var blobprov blobaccess.BlobAccessProvider
	switch eff.Compression {
	case NONE:
		blobprov = blobaccess.ProviderForFile(media, path, eff.FileSystem)
	case COMPRESSION:
		blob := blobaccess.ForFile(media, path, eff.FileSystem)
		defer blob.Close()
		blob, _ = blobaccess.WithCompression(blob)
		blobprov = blobaccess.ProviderForBlobAccess(blob)
	case DECOMPRESSION:
		blob := blobaccess.ForFile(media, path, eff.FileSystem)
		defer blob.Close()
		blob, _ = blobaccess.WithDecompression(blob)
		blobprov = blobaccess.ProviderForBlobAccess(blob)
	}

	accprov := cpi.NewAccessProviderForBlobAccessProvider(ctx, blobprov, eff.Hint, eff.Global)
	// strange type cast is required by Go compiler, meta has the correct type.
	return cpi.NewArtifactAccessForProvider(generics.As[*M](meta), accprov)
}

func ResourceAccess(ctx cpi.Context, name string, opts ...elements.ResourceMetaOption) func(media, path string, opts ...Option) (cpi.ResourceAccess, error) {
	return epi.ResourceAccessAB[string, string, Option](ctx, name, TYPE, Access[compdesc.ResourceMeta], opts...)
}

func SourceAccess(ctx cpi.Context, name string, opts ...elements.SourceMetaOption) func(media, path string, opts ...Option) (cpi.SourceAccess, error) {
	return epi.SourceAccessAB[string, string, Option](ctx, name, TYPE, Access[compdesc.SourceMeta], opts...)
}
