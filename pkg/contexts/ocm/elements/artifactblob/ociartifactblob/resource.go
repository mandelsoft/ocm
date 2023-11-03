// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package ociartifactblob

import (
	blob "github.com/open-component-model/ocm/pkg/blobaccess/ociartifact"
	"github.com/open-component-model/ocm/pkg/contexts/oci"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/cpi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements/artifactaccess/epi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/resourcetypes"
	"github.com/open-component-model/ocm/pkg/generics"
	"github.com/open-component-model/ocm/pkg/optionutils"
)

const TYPE = resourcetypes.OCI_IMAGE

func Access[M any, P compdesc.ArtifactMetaPointer[M]](ctx cpi.Context, meta P, refname string, opts ...Option) cpi.ArtifactAccess[M] {
	eff := optionutils.EvalOptions(append(opts, WithContext(ctx))...)
	if meta.GetType() == "" {
		meta.SetType(TYPE)
	}

	hint := eff.Hint
	if hint == "" {
		ref, err := oci.ParseRef(refname)
		if err == nil {
			hint = ref.String()
		}
	}

	blobprov := blob.BlobAccessProviderForOCIArtifact(refname, &eff.Blob)
	accprov := cpi.NewAccessProviderForBlobAccessProvider(ctx, blobprov, hint, eff.Global)
	// strange type cast is required by Go compiler, meta has the correct type.
	return cpi.NewArtifactAccessForProvider(generics.As[*M](meta), accprov)
}

func ResourceAccess(ctx cpi.Context, name string, opts ...elements.ResourceMetaOption) func(refname string, opts ...Option) (cpi.ResourceAccess, error) {
	return epi.ResourceAccessA[string, Option](ctx, name, TYPE, Access[compdesc.ResourceMeta], opts...)
}

func SourceAccess(ctx cpi.Context, name string, opts ...elements.SourceMetaOption) func(refname string, opts ...Option) (cpi.SourceAccess, error) {
	return epi.SourceAccessA[string, Option](ctx, name, TYPE, Access[compdesc.SourceMeta], opts...)
}
