// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package epi

import (
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/cpi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements"
)

type CreatorFunc[M, O any] func(ctx cpi.Context, meta *M, opts ...O) cpi.ArtifactAccess[M]

func ResourceAccess[O any](ctx cpi.Context, name, typ string, access CreatorFunc[compdesc.ResourceMeta, O], opts ...elements.ResourceMetaOption) func(opts ...O) (cpi.ResourceAccess, error) {
	meta, err := elements.ResourceMeta(name, typ, opts...)
	return func(opts ...O) (cpi.ResourceAccess, error) {
		if err != nil {
			return nil, err
		}

		return access(ctx, meta, opts...), nil
	}
}

func SourceAccess[O any](ctx cpi.Context, name, typ string, access CreatorFunc[compdesc.SourceMeta, O], opts ...elements.SourceMetaOption) func(opts ...O) (cpi.SourceAccess, error) {
	meta, err := elements.SourceMeta(name, typ, opts...)
	return func(opts ...O) (cpi.SourceAccess, error) {
		if err != nil {
			return nil, err
		}
		return access(ctx, meta, opts...), nil
	}
}

type CreatorFuncA[M, A, O any] func(ctx cpi.Context, meta *M, arga A, opts ...O) cpi.ArtifactAccess[M]

func ResourceAccessA[A, O any](ctx cpi.Context, name, typ string, access CreatorFuncA[compdesc.ResourceMeta, A, O], opts ...elements.ResourceMetaOption) func(arga A, opts ...O) (cpi.ResourceAccess, error) {
	meta, err := elements.ResourceMeta(name, typ, opts...)
	return func(arga A, opts ...O) (cpi.ResourceAccess, error) {
		if err != nil {
			return nil, err
		}

		return access(ctx, meta, arga, opts...), nil
	}
}

func SourceAccessA[A, O any](ctx cpi.Context, name, typ string, access CreatorFuncA[compdesc.SourceMeta, A, O], opts ...elements.SourceMetaOption) func(arga A, opts ...O) (cpi.SourceAccess, error) {
	meta, err := elements.SourceMeta(name, typ, opts...)
	return func(arga A, opts ...O) (cpi.SourceAccess, error) {
		if err != nil {
			return nil, err
		}
		return access(ctx, meta, arga, opts...), nil
	}
}

type CreatorFuncAB[M, A, B, O any] func(ctx cpi.Context, meta *M, arga A, argb B, opts ...O) cpi.ArtifactAccess[M]

func ResourceAccessAB[A, B, O any](ctx cpi.Context, name, typ string, access CreatorFuncAB[compdesc.ResourceMeta, A, B, O], opts ...elements.ResourceMetaOption) func(arga A, argb B, opts ...O) (cpi.ResourceAccess, error) {
	meta, err := elements.ResourceMeta(name, typ, opts...)
	return func(arga A, argb B, opts ...O) (cpi.ResourceAccess, error) {
		if err != nil {
			return nil, err
		}
		return access(ctx, meta, arga, argb, opts...), nil
	}
}

func SourceAccessAB[A, B, O any](ctx cpi.Context, name, typ string, access CreatorFuncAB[compdesc.SourceMeta, A, B, O], opts ...elements.SourceMetaOption) func(arga A, argb B, opts ...O) (cpi.SourceAccess, error) {
	meta, err := elements.SourceMeta(name, typ, opts...)
	return func(arga A, argb B, opts ...O) (cpi.SourceAccess, error) {
		if err != nil {
			return nil, err
		}
		return access(ctx, meta, arga, argb, opts...), nil
	}
}

type CreatorFuncABC[M, A, B, C, O any] func(ctx cpi.Context, meta *M, arga A, argb B, argc C, opts ...O) cpi.ArtifactAccess[M]

func ResourceAccessABC[A, B, C, O any](ctx cpi.Context, name, typ string, access CreatorFuncABC[compdesc.ResourceMeta, A, B, C, O], opts ...elements.ResourceMetaOption) func(arga A, argb B, argc C, opts ...O) (cpi.ResourceAccess, error) {
	meta, err := elements.ResourceMeta(name, typ, opts...)
	return func(arga A, argb B, argc C, opts ...O) (cpi.ResourceAccess, error) {
		if err != nil {
			return nil, err
		}
		return access(ctx, meta, arga, argb, argc, opts...), nil
	}
}

func SourceAccessABC[A, B, C, O any](ctx cpi.Context, name, typ string, access CreatorFuncABC[compdesc.SourceMeta, A, B, C, O], opts ...elements.SourceMetaOption) func(arga A, argb B, argc C, opts ...O) (cpi.SourceAccess, error) {
	meta, err := elements.SourceMeta(name, typ, opts...)
	return func(arga A, argb B, argc C, opts ...O) (cpi.SourceAccess, error) {
		if err != nil {
			return nil, err
		}
		return access(ctx, meta, arga, argb, argc, opts...), nil
	}
}
