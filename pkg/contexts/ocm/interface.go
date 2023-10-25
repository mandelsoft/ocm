// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package ocm

import (
	"context"
	"io"

	"github.com/open-component-model/ocm/pkg/blobaccess"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc"
	metav1 "github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc/meta/v1"
	ocm "github.com/open-component-model/ocm/pkg/contexts/ocm/context"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/cpi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/internal"
	"github.com/open-component-model/ocm/pkg/runtime"
)

const (
	KIND_COMPONENTVERSION   = internal.KIND_COMPONENTVERSION
	KIND_COMPONENTREFERENCE = "component reference"
	KIND_RESOURCE           = internal.KIND_RESOURCE
	KIND_SOURCE             = internal.KIND_SOURCE
	KIND_REFERENCE          = internal.KIND_REFERENCE
	KIND_REPOSITORYSPEC     = internal.KIND_REPOSITORYSPEC
)

const CONTEXT_TYPE = internal.CONTEXT_TYPE

const CommonTransportFormat = internal.CommonTransportFormat

type (
	Context                          = internal.Context
	ContextProvider                  = internal.ContextProvider
	LocalContextProvider             = internal.LocalContextProvider
	ComponentVersionResolver         = internal.ComponentVersionResolver
	Repository                       = internal.Repository
	RepositorySpecHandlers           = internal.RepositorySpecHandlers
	RepositorySpecHandler            = internal.RepositorySpecHandler
	UniformRepositorySpec            = internal.UniformRepositorySpec
	ComponentLister                  = internal.ComponentLister
	ComponentAccess                  = internal.ComponentAccess
	ComponentVersionAccess           = internal.ComponentVersionAccess
	AccessSpec                       = internal.AccessSpec
	GenericAccessSpec                = internal.GenericAccessSpec
	HintProvider                     = internal.HintProvider
	AccessMethod                     = internal.AccessMethod
	AccessType                       = internal.AccessType
	DataAccess                       = internal.DataAccess
	BlobAccess                       = internal.BlobAccess
	AccessProvider                   = internal.AccessProvider
	SourceAccess                     = internal.SourceAccess
	SourceMeta                       = internal.SourceMeta
	ResourceAccess                   = internal.ResourceAccess
	ResourceMeta                     = internal.ResourceMeta
	RepositorySpec                   = internal.RepositorySpec
	GenericRepositorySpec            = internal.GenericRepositorySpec
	IntermediateRepositorySpecAspect = internal.IntermediateRepositorySpecAspect
	RepositoryType                   = internal.RepositoryType
	RepositoryTypeScheme             = internal.RepositoryTypeScheme
	RepositoryDelegationRegistry     = internal.RepositoryDelegationRegistry
	AccessTypeScheme                 = internal.AccessTypeScheme
	ComponentReference               = internal.ComponentReference
)

type (
	DigesterType         = internal.DigesterType
	BlobDigester         = internal.BlobDigester
	BlobDigesterRegistry = internal.BlobDigesterRegistry
	DigestDescriptor     = internal.DigestDescriptor
	HasherProvider       = internal.HasherProvider
	Hasher               = internal.Hasher
)

type (
	BlobHandlerRegistry = internal.BlobHandlerRegistry
	BlobHandler         = internal.BlobHandler
	BlobHandlerProvider = internal.BlobHandlerProvider
)

func NewDigestDescriptor(digest, hashAlgo, normAlgo string) *DigestDescriptor {
	return internal.NewDigestDescriptor(digest, hashAlgo, normAlgo)
}

// DefaultContext is the default context initialized by init functions.
func DefaultContext() internal.Context {
	return internal.DefaultContext
}

func DefaultBlobHandlers() BlobHandlerRegistry {
	return internal.DefaultBlobHandlerRegistry
}

func DefaultBlobHandlerProvider(ctx Context) BlobHandlerProvider {
	return internal.DefaultBlobHandlerProvider(ctx)
}

func DefaultRepositoryDelegationRegistry() RepositoryDelegationRegistry {
	return internal.DefaultRepositoryDelegationRegistry
}

func NewRepositoryDelegationRegistry(base ...RepositoryDelegationRegistry) RepositoryDelegationRegistry {
	return internal.NewDelegationRegistry[Context, RepositorySpec](base...)
}

// FromContext returns the Context to use for context.Context.
// This is either an explicit context or the default context.
func FromContext(ctx context.Context) Context {
	return internal.FromContext(ctx)
}

func FromProvider(p ContextProvider) Context {
	return internal.FromProvider(p)
}

func DefinedForContext(ctx context.Context) (Context, bool) {
	return internal.DefinedForContext(ctx)
}

func NewGenericAccessSpec(spec string) (AccessSpec, error) {
	return internal.NewGenericAccessSpec([]byte(spec))
}

func ToGenericAccessSpec(spec AccessSpec) (*GenericAccessSpec, error) {
	return internal.ToGenericAccessSpec(spec)
}

func ToGenericRepositorySpec(spec RepositorySpec) (*GenericRepositorySpec, error) {
	return internal.ToGenericRepositorySpec(spec)
}

func IsNoneAccess(a compdesc.AccessSpec) bool {
	return compdesc.IsNoneAccess(a)
}

func IsNoneAccessKind(k string) bool {
	return compdesc.IsNoneAccessKind(k)
}

type AccessSpecRef = internal.AccessSpecRef

func NewAccessSpecRef(spec cpi.AccessSpec) *AccessSpecRef {
	return internal.NewAccessSpecRef(spec)
}

func NewRawAccessSpecRef(data []byte, unmarshaler runtime.Unmarshaler) (*AccessSpecRef, error) {
	return internal.NewRawAccessSpecRef(data, unmarshaler)
}

func NewResourceMeta(name string, typ string, relation metav1.ResourceRelation) *ResourceMeta {
	return compdesc.NewResourceMeta(name, typ, relation)
}

///////////////////////////////////////////////////////

type AccessMethodView = cpi.AccessMethodView

func BlobAccessForAccessMethod(m AccessMethodView) (blobaccess.AnnotatedBlobAccess[AccessMethodView], error) {
	return cpi.BlobAccessForAccessMethod(m)
}

// AccessMethodAsView wrap an access method object into
// a multi-view version. The original method is closed when
// the last view is closed.
// After an access method is used as base object, it should not
// explicitly closed anymore, because the views will stop
// functioning.
func AccessMethodAsView(acc ocm.AccessMethod, closer ...io.Closer) AccessMethodView {
	return cpi.AccessMethodAsView(acc, closer...)
}

func AccessMethodViewForSpec(spec ocm.AccessSpec, cv ocm.ComponentVersionAccess) (AccessMethodView, error) {
	return cpi.AccessMethodViewForSpec(spec, cv)
}

func AccessMethodViewForAccessProvider(p AccessProvider) (AccessMethodView, error) {
	return cpi.AccessMethodViewForAccessProvider(p)
}
