// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"github.com/open-component-model/ocm/examples/lib/helper"
	"github.com/open-component-model/ocm/pkg/contexts/ocm"
	metav1 "github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc/meta/v1"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements/artifactaccess/helmaccess"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements/artifactaccess/ociartifactaccess"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/elements/artifactblob/fileblob"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/repositories/composition"
	"github.com/open-component-model/ocm/pkg/errors"
	"github.com/open-component-model/ocm/pkg/mime"
)

const (
	PODINFO_IMAGE  = "ghcr.io/stefanprodan/podinfo:6.5.2"
	HELMCHART_REPO = "oci://ghcr.io/stefanprodan/charts"
	HELMCHART_NAME = "podinfo:6.5.2"
)

const (
	RSC_IMAGE     = "podinfo-image"
	RSC_HELMCHART = "helmchart"
	RSC_DEPLOY    = "deplyscript"
)
const DEPLOY_SCRIPT_TYPE = "helmDeployScript"

func Create(cfg *helper.Config) error {
	ctx := ocm.DefaultContext()

	cv, err := CreateComponentVersion(ctx)
	if err != nil {
		return err
	}
	defer cv.Close()

	DescribeVersion(cv)
	return nil
}

// CreateComponentVersion creates the scenario component version with
// three resources: the podinfo image, the helm chart and a locally
// found deploy script.
func CreateComponentVersion(ctx ocm.Context) (ocm.ComponentVersionAccess, error) {
	fmt.Printf("*** composing component version %s:%s\n", COMPONENT_NAME, COMPONENT_VERSION)

	cv := composition.NewComponentVersion(ctx, COMPONENT_NAME, COMPONENT_VERSION)

	cv.SetProvider(&metav1.Provider{
		Name: "acme.org",
	})

	// podinfo image as external resource reference
	image_res, err := ociartifactaccess.ResourceAccess(ctx, RSC_IMAGE, PODINFO_IMAGE)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create resource meta for podinfo-image")
	}
	err = cv.SetResourceAccess(image_res)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot add resource podinfo-image")
	}

	// helm chart as external resource reference
	helm_res, err := helmaccess.ResourceAccess(ctx, RSC_HELMCHART, HELMCHART_NAME, HELMCHART_REPO)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create resource meta for helmchart")
	}
	err = cv.SetResourceAccess(helm_res)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot add resource helmchart")
	}

	// deploy script found in filesystem
	script_res, err := fileblob.ResourceAccess(ctx, RSC_DEPLOY, elements.WithType(DEPLOY_SCRIPT_TYPE))(mime.MIME_YAML, "resources/deployscript")
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create resource meta for podinfo-image")
	}
	err = cv.SetResourceAccess(script_res)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot add resource helmchart")
	}

	return cv, nil
}
