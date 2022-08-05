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

package install

import (
	"encoding/json"

	metav1 "github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc/meta/v1"
)

const TypeTOIPackage = "toiPackage"
const PackageSpecificationMimeType = "application/vnd.toi.gardener.cloud.package.v1+yaml"

const TypeTOIExecutor = "toiExecutor"
const ExecutorSpecificationMimeType = "application/vnd.toi.gardener.cloud.executor.v1+yaml"

type PackageSpecification struct {
	CredentialsRequest `json:",inline"`
	Template           json.RawMessage            `json:"configTemplate"`
	Libraries          []metav1.ResourceReference `json:"templateLibraries"`
	Scheme             json.RawMessage            `json:"configScheme"`
	Executors          []Executor                 `json:"executors"`
}

type Executor struct {
	Actions          []string                 `json:"actions,omitempty"`
	ResourceRef      metav1.ResourceReference `json:"resourceRef"`
	Image            *Image                   `json:"image,omitempty"`
	ParameterMapping json.RawMessage          `json:"parameterMapping,omitempty"`
	Config           json.RawMessage          `json:"config,omitempty"`
	Outputs          map[string]string        `json:"outputs,omitempty"`
}

type Image struct {
	Ref    string `json:"image"`
	Digest string `json:"digest"`
}

////////////////////////////////////////////////////////////////////////////////

type ExecutorSpecification struct {
	CredentialsRequest `json:",inline"`
	Actions            []string                   `json:"actions,omitempty"`
	ResourceRef        metav1.ResourceReference   `json:"resourceRef"`
	ConfigTemplate     json.RawMessage            `json:"configTemplate"`
	Libraries          []metav1.ResourceReference `json:"templateLibraries"`
	Scheme             json.RawMessage            `json:"configScheme,omitempty"`
	Outputs            map[string]OutputSpec      `json:"outputs,omitempty"`
}

type OutputSpec struct {
	Description string `json:"description,omitempty"`
}
