// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/open-component-model/ocm/local/r3trans/plugin/accessmethods"
	"github.com/open-component-model/ocm/local/r3trans/plugin/config"
	"github.com/open-component-model/ocm/local/r3trans/plugin/uploaders"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/plugin/ppi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/plugin/ppi/cmds"
)

func main() {
	p := ppi.NewPlugin("r3trans.sap.com", "0.0.1")

	p.SetShort("SAP R/3 Transport System Adapter Demo")
	p.SetLong(`This plugin provided support to access R/3 transport requests and to upload
and to upload them to a transport environment again.

The plugin uses the following configuration fields:

- **<code>systems</code>** *map[string]&lt;config>*

  The configuration used  fr a set of transport systems:

  - **<code>path</code>** *string* (default <code>/tmp/r3trans</code>)

    The base address to be used for the transport system.
`)
	p.SetConfigParser(config.GetConfig)

	p.RegisterAccessMethod(accessmethods.New())
	u := uploaders.New()
	p.RegisterUploader("r3trans.sap.com/transportRequest", "", u)
	err := cmds.NewPluginCommand(p).Execute(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
}
