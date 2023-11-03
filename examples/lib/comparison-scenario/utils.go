// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/open-component-model/ocm/pkg/contexts/credentials"
	"github.com/open-component-model/ocm/pkg/contexts/ocm"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/attrs/signingattr"
	metav1 "github.com/open-component-model/ocm/pkg/contexts/ocm/compdesc/meta/v1"
	"github.com/open-component-model/ocm/pkg/signing/handlers/rsa"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli"
)

func DescribeVersion(cv ocm.ComponentVersionAccess) error {
	// many elements of the API keep trak of their context
	ctx := cv.GetContext()

	// Have a look at the component descriptor
	cd := cv.GetDescriptor()
	fmt.Printf("resources of version of %s:%s:\n", cv.GetName(), cv.GetVersion())
	fmt.Printf("  provider: %s\n", cd.Provider.Name)

	// and list all the included resources.
	for i, r := range cv.GetResources() {
		fmt.Printf("  %2d: name:           %s\n", i+1, r.Meta().GetName())
		fmt.Printf("      extra identity: %s\n", r.Meta().GetExtraIdentity())
		fmt.Printf("      resource type:  %s\n", r.Meta().GetType())
		acc, err := r.Access()
		if err != nil {
			fmt.Printf("      access:         error: %s\n", err)
		} else {
			fmt.Printf("      access:         %s\n", acc.Describe(ctx))
		}
		PrintDigest(r.Meta().Digest)
	}
	return nil
}

func PrintPublicKey(ctx ocm.Context, name string) {
	info := signingattr.Get(ctx)
	key := info.GetPublicKey(name)
	if key == nil {
		fmt.Printf("public key for %s not found\n", name)
	} else {
		buf := bytes.NewBuffer(nil)
		err := rsa.WriteKeyData(key, buf)
		if err != nil {
			fmt.Printf("key error: %s\n", err)
		} else {
			fmt.Printf("public key for %s:\n%s\n", name, buf.String())
		}
	}
}

func PrintDigest(dig *metav1.DigestSpec) {
	fmt.Printf("      digest:\n")
	fmt.Printf("        algorithm:     %s\n", dig.HashAlgorithm)
	fmt.Printf("        normalization: %s\n", dig.NormalisationAlgorithm)
	fmt.Printf("        value:         %s\n", dig.Value)

}
func PrintSignatures(cv ocm.ComponentVersionAccess) {
	fmt.Printf("signatures:\n")
	for i, s := range cv.GetDescriptor().Signatures {
		fmt.Printf("%2d    name: %s\n", i, s.Name)
		PrintDigest(&s.Digest)
		fmt.Printf("      signature:\n")
		fmt.Printf("        algorithm: %s\n", s.Signature.Algorithm)
		fmt.Printf("        mediaType: %s\n", s.Signature.MediaType)
		fmt.Printf("        value:     %s\n", s.Signature.Value)
	}
}

func PrintConsumerId(o interface{}, msg string) {
	// register credentials for given OCI registry in context.
	id := credentials.GetProvidedConsumerId(o)
	if id == nil {
		fmt.Printf("no consumer id for %s\n", msg)
	} else {
		fmt.Printf("consumer id for %s: %s\n", msg, id)
	}
}

func InstallChart(chart *chart.Chart, release, namespace string) error {
	settings := cli.New()
	settings.SetNamespace(namespace)
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(
		settings.RESTClientGetter(),
		namespace,
		os.Getenv("HELM_DRIVER"),
		func(msg string, args ...interface{}) { fmt.Printf(msg, args...) },
	); err != nil {
		return err
	}

	client := action.NewInstall(actionConfig)
	client.ReleaseName = release
	client.Namespace = namespace
	if _, err := client.Run(chart, nil); err != nil {
		return err
	}

	return nil
}

func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		panic(err)
	}
}
