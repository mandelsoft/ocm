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
	"sync"

	"github.com/gardener/ocm/pkg/credentials"
	"github.com/gardener/ocm/pkg/oci/cpi"
)

const ATTR_REPOS = "github.com/gardener/ocm/pkg/oci/repositories/ocireg"

type Repositories struct {
	lock  sync.Mutex
	repos map[string]*Repository
}

func newRepositories(context.Context) interface{} {
	return &Repositories{
		repos: map[string]*Repository{},
	}
}

func (r *Repositories) GetRepository(ctx cpi.Context, spec *RepositorySpec, creds credentials.Credentials) (*Repository, error) {
	name := spec.BaseURL // TODO: switch to ID
	r.lock.Lock()
	defer r.lock.Unlock()
	repo := r.repos[name]
	if repo == nil {
		new, err := spec.Repository(ctx, creds)
		if err != nil {
			return nil, err
		}
		repo = new.(*Repository)
		r.repos[name] = repo
	}
	return repo, nil
}
