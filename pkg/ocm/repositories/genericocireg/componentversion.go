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

package genericocireg

import (
	"path"
	"strings"

	"github.com/gardener/ocm/pkg/common"
	"github.com/gardener/ocm/pkg/common/accessobj"
	"github.com/gardener/ocm/pkg/errors"
	"github.com/gardener/ocm/pkg/oci"
	"github.com/gardener/ocm/pkg/oci/artdesc"
	"github.com/gardener/ocm/pkg/oci/repositories/ctf/artefactset"
	"github.com/gardener/ocm/pkg/ocm/accessmethods"
	ocihdlr "github.com/gardener/ocm/pkg/ocm/blobhandler/oci"
	"github.com/gardener/ocm/pkg/ocm/compdesc"
	"github.com/gardener/ocm/pkg/ocm/cpi"
	"github.com/gardener/ocm/pkg/ocm/repositories/ctf/comparch"
	"github.com/opencontainers/go-digest"
)

type ComponentVersion struct {
	container *ComponentVersionContainer
	*comparch.ComponentVersionAccess
}

var _ cpi.ComponentVersionAccess = (*ComponentVersion)(nil)

func NewComponentVersionAccess(mode accessobj.AccessMode, comp *ComponentAccess, version string, access oci.ManifestAccess) (*ComponentVersion, error) {
	c, err := NewComponentVersionContainer(mode, comp, version, access)
	if err != nil {
		return nil, err
	}
	return &ComponentVersion{
		container:              c,
		ComponentVersionAccess: comparch.NewComponentVersionAccess(c),
	}, nil
}

////////////////////////////////////////////////////////////////////////////////

type ComponentVersionContainer struct {
	comp     *ComponentAccess
	version  string
	manifest oci.ManifestAccess
	state    accessobj.State
}

var _ comparch.ComponentVersionContainer = (*ComponentVersionContainer)(nil)

func NewComponentVersionContainer(mode accessobj.AccessMode, comp *ComponentAccess, version string, manifest oci.ManifestAccess) (*ComponentVersionContainer, error) {
	state, err := NewState(mode, comp.name, version, manifest)
	if err != nil {
		return nil, err
	}
	return &ComponentVersionContainer{
		comp:     comp,
		version:  version,
		manifest: manifest,
		state:    state,
	}, nil
}

func (c *ComponentVersionContainer) Check() error {
	if c.version != c.GetDescriptor().Version {
		return errors.ErrInvalid("component version", c.GetDescriptor().Version)
	}
	if c.comp.name != c.GetDescriptor().Name {
		return errors.ErrInvalid("component name", c.GetDescriptor().Name)
	}
	return nil
}

func (c *ComponentVersionContainer) GetContext() cpi.Context {
	return c.comp.GetContext()
}

func (c *ComponentVersionContainer) IsReadOnly() bool {
	return c.state.IsReadOnly()
}

func (c *ComponentVersionContainer) IsClosed() bool {
	return c.manifest == nil
}

func (c *ComponentVersionContainer) Update() error {
	err := c.Check()
	desc := c.GetDescriptor()
	for i, r := range desc.Resources {
		s, err := c.evalLayer(r.Access)
		if err != nil {
			return err
		}
		if s != r.Access {
			desc.Resources[i].Access = s
		}
	}
	for i, r := range desc.Sources {
		s, err := c.evalLayer(r.Access)
		if err != nil {
			return err
		}
		if s != r.Access {
			desc.Sources[i].Access = s
		}
	}
	_, err = c.state.Update()
	if err != nil {
		return err
	}
	_, err = c.comp.namespace.AddArtefact(c.manifest, c.version)
	if err != nil {
		return err
	}
	return nil
}

func (c *ComponentVersionContainer) evalLayer(s compdesc.AccessSpec) (compdesc.AccessSpec, error) {
	spec, err := c.GetContext().AccessSpecForSpec(s)
	if err != nil {
		return s, err
	}
	if a, ok := spec.(*accessmethods.LocalBlobAccessSpec); ok {
		if ok, _ := artdesc.IsDigest(a.LocalReference); !ok {
			return s, errors.ErrInvalid("digest", a.LocalReference)
		}
	}
	return s, nil
}

func (c *ComponentVersionContainer) GetDescriptor() *compdesc.ComponentDescriptor {
	return c.state.GetState().(*compdesc.ComponentDescriptor)
}

func (c *ComponentVersionContainer) GetBlobData(name string) (cpi.DataAccess, error) {
	return c.manifest.GetBlob(digest.Digest((name)))
}

func (c *ComponentVersionContainer) AddBlob(blob cpi.BlobAccess, refName string, global cpi.AccessSpec) (cpi.AccessSpec, error) {
	if blob == nil {
		return nil, errors.New("a resource has to be defined")
	}

	storagectx := ocihdlr.New(c.comp.repo.ocirepo, c.comp.namespace, c.manifest)
	h := c.GetContext().BlobHandlers().GetHandler(oci.CONTEXT_TYPE, c.comp.repo.ocirepo.GetSpecification().GetKind(), blob.MimeType())
	if h != nil {
		acc, err := h.StoreBlob(c.comp.repo, blob, refName, storagectx)
		if err != nil {
			return nil, err
		}
		if acc != nil {
			return acc, nil
		}
	}

	err := c.manifest.AddBlob(blob)
	if err != nil {
		return nil, err
	}
	err = storagectx.AssureLayer(blob)
	if err != nil {
		return nil, err
	}
	return accessmethods.NewLocalBlobAccessSpec(common.DigestToFileName(blob.Digest()), refName, blob.MimeType(), global), nil
}

// assureGlobalRef provides a global manifest for a local OCI Artefact
func (c *ComponentVersionContainer) assureGlobalRef(d digest.Digest, url, name string) (cpi.AccessSpec, error) {

	blob, err := c.manifest.GetBlob(d)
	if err != nil {
		return nil, err
	}
	var namespace oci.NamespaceAccess
	var version string
	var tag string
	if name == "" {
		namespace = c.comp.namespace
	} else {
		i := strings.LastIndex(name, ":")
		if i > 0 {
			version = name[i+1:]
			name = name[:i]
			tag = version
		}
		namespace, err = c.comp.repo.ocirepo.LookupNamespace(name)
		if err != nil {
			return nil, err
		}
	}
	set, err := artefactset.OpenFromBlob(accessobj.ACC_READONLY, blob)
	if err != nil {
		return nil, err
	}
	defer set.Close()
	digest := set.GetMain()
	if version == "" {
		version = digest.String()
	}
	art, err := set.GetArtefact(digest.String())
	if err != nil {
		return nil, err
	}
	err = artefactset.TransferArtefact(art, namespace, oci.AsTags(tag)...)
	if err != nil {
		return nil, err
	}

	ref := path.Join(url+namespace.GetNamespace()) + ":" + version

	global := accessmethods.NewOCIRegistryAccessSpec(ref)
	return global, nil
}