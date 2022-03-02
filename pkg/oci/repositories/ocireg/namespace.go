// Copyright 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package ocireg

import (
	"context"

	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/remotes"
	"github.com/gardener/ocm/pkg/common/accessio"
	"github.com/gardener/ocm/pkg/errors"
	"github.com/gardener/ocm/pkg/oci/artdesc"
	"github.com/gardener/ocm/pkg/oci/cpi"
	"github.com/opencontainers/go-digest"
)

type Namespace struct {
	access *NamespaceContainer
}

func (n *Namespace) Close() error {
	return nil
}

type NamespaceContainer struct {
	repo      *Repository
	namespace string
	resolver  remotes.Resolver
	fetcher   remotes.Fetcher
	pusher    remotes.Pusher
}

var _ cpi.ArtefactSetContainer = (*NamespaceContainer)(nil)
var _ cpi.NamespaceAccess = (*Namespace)(nil)

func NewNamespace(repo *Repository, name string) (*Namespace, error) {
	ref := repo.getRef(name, "")
	resolver, err := repo.getResolver(ref)
	if err != nil {
		return nil, err
	}
	fetcher, err := resolver.Fetcher(context.Background(), ref)
	if err != nil {
		return nil, err
	}
	pusher, err := resolver.Pusher(context.Background(), ref)
	if err != nil {
		return nil, err
	}
	n := &Namespace{
		access: &NamespaceContainer{
			repo:      repo,
			namespace: name,
			resolver:  resolver,
			fetcher:   fetcher,
			pusher:    pusher,
		},
	}
	return n, nil
}

func (n *NamespaceContainer) getPusher(vers string) (remotes.Pusher, error) {
	ref := n.repo.getRef(n.namespace, vers)
	return n.resolver.Pusher(dummyContext, ref)
}

func (n *NamespaceContainer) push(vers string, blob cpi.BlobAccess) error {
	ref := n.repo.getRef(n.namespace, vers)
	p, err := n.resolver.Pusher(dummyContext, ref)
	if err != nil {
		return err
	}
	return push(dummyContext, p, blob)
}

func (n *NamespaceContainer) GetNamepace() string {
	return n.namespace
}

func (n *NamespaceContainer) IsReadOnly() bool {
	return n.repo.IsReadOnly()
}

func (n *NamespaceContainer) IsClosed() bool {
	return n.repo.IsClosed()
}

func (n *NamespaceContainer) GetBlobDescriptor(digest digest.Digest) *cpi.Descriptor {
	return nil
}

func (n *NamespaceContainer) ListTags() ([]string, error) {
	panic("implement me")
}

func (n *NamespaceContainer) GetBlobData(digest digest.Digest) (cpi.DataAccess, error) {
	return NewDataAcess(n.fetcher, digest, "", false)
}

func (n *NamespaceContainer) AddBlob(blob cpi.BlobAccess) error {
	return push(dummyContext, n.pusher, blob)
}

func (n *NamespaceContainer) GetArtefact(vers string) (cpi.ArtefactAccess, error) {
	ref := n.repo.getRef(n.namespace, vers)
	_, desc, err := n.resolver.Resolve(context.Background(), ref)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return nil, errors.ErrNotFound(cpi.KIND_OCIARTEFACT, ref, n.namespace)
		}
		return nil, err
	}
	acc, err := NewDataAcess(n.fetcher, desc.Digest, desc.MediaType, false)
	if err != nil {
		return nil, err
	}
	return cpi.NewArtefactForBlob(n, accessio.BlobAccessForDataAccess(desc.Digest, desc.Size, desc.MediaType, acc))
}

func (n *NamespaceContainer) AddArtefact(artefact cpi.Artefact, platform *artdesc.Platform) (access accessio.BlobAccess, err error) {
	blob, err := artefact.Blob()
	if err != nil {
		return nil, err
	}
	return blob, n.push(blob.Digest().String(), blob)
}

func (n *NamespaceContainer) AddTags(digest digest.Digest, tags ...string) error {
	_, desc, err := n.resolver.Resolve(context.Background(), n.repo.getRef(n.namespace, digest.String()))
	if err != nil {
		return err
	}
	acc, err := NewDataAcess(n.fetcher, desc.Digest, desc.MediaType, false)
	if err != nil {
		return err
	}
	blob := accessio.BlobAccessForDataAccess(desc.Digest, desc.Size, desc.MediaType, acc)
	for _, tag := range tags {
		err := n.push(tag, blob)
		if err != nil {
			return err
		}
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func (n *Namespace) GetRepository() cpi.Repository {
	return n.access.repo
}

func (n *Namespace) GetNamespace() string {
	return n.access.GetNamepace()
}

func (n *Namespace) ListTags() ([]string, error) {
	return n.access.ListTags()
}

func (n *Namespace) NewArtefact(art ...*artdesc.Artefact) (cpi.ArtefactAccess, error) {
	if n.access.IsReadOnly() {
		return nil, accessio.ErrReadOnly
	}
	return cpi.NewArtefact(n.access, art...), nil
}

func (n *Namespace) GetBlobData(digest digest.Digest) (cpi.DataAccess, error) {
	return n.access.GetBlobData(digest)
}

func (n *Namespace) GetArtefact(vers string) (cpi.ArtefactAccess, error) {
	return n.access.GetArtefact(vers)
}

func (n *Namespace) AddTaggedArtefact(artefact cpi.Artefact, tags ...string) (accessio.BlobAccess, error) {
	blob, err := n.access.AddArtefact(artefact, nil)
	if err != nil {
		return nil, err
	}
	return blob, n.AddTags(blob.Digest(), tags...)
}

func (n *Namespace) AddTags(digest digest.Digest, tags ...string) error {
	return n.access.AddTags(digest, tags...)
}

func (n *Namespace) AddBlob(blob cpi.BlobAccess) error {
	return n.access.AddBlob(blob)
}
