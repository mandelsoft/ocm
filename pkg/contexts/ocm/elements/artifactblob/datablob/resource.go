// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package datablob

import (
	"github.com/open-component-model/ocm/pkg/blobaccess"
	"github.com/open-component-model/ocm/pkg/contexts/ocm"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/cpi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements/artifactaccess/epi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/resourcetypes"
	"github.com/open-component-model/ocm/pkg/mime"
	"github.com/open-component-model/ocm/pkg/optionutils"
)

const TYPE = resourcetypes.BLOB

func Access[M any, P compdesc.ArtifactMetaPointer[M]](ctx ocm.Context, meta P, blob []byte, opts ...Option) cpi.ArtifactAccess[M] {
	eff := optionutils.EvalOptions(opts...)

	media := eff.MimeType
	if media == "" {
		media = mime.MIME_OCTET
	}
	if meta.GetType() == "" {
		meta.SetType(TYPE)
	}

	var blobprov blobaccess.BlobAccessProvider
	switch eff.Compression {
	case NONE:
		blobprov = blobaccess.ProviderForData(media, blob)
	case COMPRESSION:
		blob := blobaccess.ForData(media, blob)
		defer blob.Close()
		blob, _ = blobaccess.WithCompression(blob)
		blobprov = blobaccess.ProviderForBlobAccess(blob)
	case DECOMPRESSION:
		blob := blobaccess.ForData(media, blob)
		defer blob.Close()
		blob, _ = blobaccess.WithDecompression(blob)
		blobprov = blobaccess.ProviderForBlobAccess(blob)
	}

	accprov := cpi.NewAccessProviderForBlobAccessProvider(ctx, blobprov, eff.Hint, eff.Global)
	return cpi.NewArtifactAccessForProvider(meta, accprov)
}

func ResourceAccess(ctx cpi.Context, name, typ string, opts ...elements.ResourceMetaOption) func(blob []byte, opts ...Option) (cpi.ResourceAccess, error) {
	return epi.ResourceAccessA[[]byte, Option](ctx, name, typ, Access[compdesc.ResourceMeta], opts...)
}

func SourceAccess(ctx cpi.Context, name, typ string, opts ...elements.SourceMetaOption) func(blob []byte, opts ...Option) (cpi.SourceAccess, error) {
	return epi.SourceAccessA[[]byte, Option](ctx, name, typ, Access[compdesc.SourceMeta], opts...)
}
