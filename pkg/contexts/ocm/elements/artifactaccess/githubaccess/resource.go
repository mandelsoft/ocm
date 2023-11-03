// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package githubaccess

import (
	access "github.com/open-component-model/ocm/pkg/contexts/ocm/accessmethods/github"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/cpi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements/artifactaccess/epi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements/artifactaccess/genericaccess"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/resourcetypes"
	"github.com/open-component-model/ocm/pkg/optionutils"
)

const TYPE = resourcetypes.DIRECTORY_TREE

func Access[M any, P compdesc.ArtifactMetaPointer[M]](ctx cpi.Context, meta P, repo, commit string, opts ...Option) cpi.ArtifactAccess[M] {
	eff := optionutils.EvalOptions(opts...)
	if meta.GetType() == "" {
		meta.SetType(TYPE)
	}

	spec := access.New(repo, eff.APIHostName, commit)
	// is global access, must work, otherwise there is an error in the lib.
	return genericaccess.MustAccess(ctx, meta, spec)
}

func ResourceAccess(ctx cpi.Context, name string, opts ...elements.ResourceMetaOption) func(repo, commit string, opts ...Option) (cpi.ResourceAccess, error) {
	return epi.ResourceAccessAB[string, string, Option](ctx, name, TYPE, Access[compdesc.ResourceMeta], opts...)
}

func SourceAccess(ctx cpi.Context, name string, opts ...elements.SourceMetaOption) func(repo, commit string, opts ...Option) (cpi.SourceAccess, error) {
	return epi.SourceAccessAB[string, string, Option](ctx, name, TYPE, Access[compdesc.SourceMeta], opts...)
}
