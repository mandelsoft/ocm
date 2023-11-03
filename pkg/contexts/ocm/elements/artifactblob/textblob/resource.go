// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package textblob

import (
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/cpi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements/artifactaccess/epi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements/artifactblob/datablob"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/resourcetypes"
	"github.com/open-component-model/ocm/pkg/mime"
	"github.com/open-component-model/ocm/pkg/optionutils"
)

const TYPE = resourcetypes.PLAIN_TEXT

func Access[M any, P compdesc.ArtifactMetaPointer[M]](ctx cpi.Context, meta P, blob string, opts ...Option) cpi.ArtifactAccess[M] {
	eff := optionutils.EvalOptions(opts...)
	if eff.MimeType == "" {
		eff.MimeType = mime.MIME_TEXT
	}
	if meta.GetType() == "" {
		meta.SetType(TYPE)
	}
	return datablob.Access(ctx, meta, []byte(blob), eff)
}

func ResourceAccess(ctx cpi.Context, name string, opts ...elements.ResourceMetaOption) func(blob string, opts ...Option) (cpi.ResourceAccess, error) {
	return epi.ResourceAccessA[string, Option](ctx, name, TYPE, Access[compdesc.ResourceMeta], opts...)
}

func SourceAccess(ctx cpi.Context, name string, opts ...elements.SourceMetaOption) func(blob string, opts ...Option) (cpi.SourceAccess, error) {
	return epi.SourceAccessA[string, Option](ctx, name, TYPE, Access[compdesc.SourceMeta], opts...)
}
