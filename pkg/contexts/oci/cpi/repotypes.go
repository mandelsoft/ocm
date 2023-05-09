// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package cpi

// this file is identical for contexts oci and credentials and simillar for
// ocm.

import (
	"github.com/open-component-model/ocm/pkg/errors"
	"github.com/open-component-model/ocm/pkg/runtime"
)

type RepositoryTypeVersionScheme interface {
	Register(t RepositoryType) error
	AddToScheme(scheme RepositoryTypeScheme)
	runtime.TypedObjectEncoder
	runtime.TypedObjectDecoder
}

type repositoryTypeVersionScheme struct {
	kind   string
	scheme RepositoryTypeScheme
}

func NewRepositoryTypeVersionScheme(kind string) RepositoryTypeVersionScheme {
	return &repositoryTypeVersionScheme{kind, newStrictRepositoryTypeScheme()}
}

func (s *repositoryTypeVersionScheme) Register(t RepositoryType) error {
	if t.GetKind() != s.kind {
		return errors.ErrInvalid("access spec type", t.GetType(), "kind", s.kind)
	}
	s.scheme.Register(t.GetType(), t)
	return nil
}

func (s *repositoryTypeVersionScheme) AddToScheme(scheme RepositoryTypeScheme) {
	scheme.AddKnownTypes(s.scheme)
}

func (s *repositoryTypeVersionScheme) Encode(obj runtime.TypedObject, m runtime.Marshaler) ([]byte, error) {
	return s.scheme.Encode(obj, m)
}

func (s *repositoryTypeVersionScheme) Decode(data []byte, unmarshaler runtime.Unmarshaler) (runtime.TypedObject, error) {
	return s.scheme.Decode(data, unmarshaler)
}

func RegisterRepositoryType(rtype RepositoryType) {
	defaultRepositoryTypeScheme.Register(rtype.GetType(), rtype)
}

func RegisterRepositoryTypeVersions(s RepositoryTypeVersionScheme) {
	s.AddToScheme(defaultRepositoryTypeScheme)
}

////////////////////////////////////////////////////////////////////////////////

func NewRepositoryType(name string, proto RepositorySpec) RepositoryType {
	return runtime.NewVersionedTypeByProto(name, proto)
}

////////////////////////////////////////////////////////////////////////////////

type (
	RepositorySpecConverter = runtime.Converter[RepositorySpec]
	RepositorySpecVersion   = runtime.FormatVersion[RepositorySpec]
)

func NewRepositorySpecVersion(proto runtime.VersionedTypedObject, converter RepositorySpecConverter) RepositorySpecVersion {
	return runtime.NewProtoBasedVersion[RepositorySpec](proto, converter)
}

////////////////////////////////////////////////////////////////////////////////

type ConvertedRepositoryType = runtime.ConvertedType[RepositorySpec]

func NewConvertedRepositorySpecType(name string, v RepositorySpecVersion) *ConvertedRepositoryType {
	return runtime.NewConvertedType[RepositorySpec](name, v)
}
