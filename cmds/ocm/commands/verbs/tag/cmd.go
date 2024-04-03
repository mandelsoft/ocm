// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package tag

import (
	"github.com/spf13/cobra"

	artifacts "github.com/open-component-model/ocm/cmds/ocm/commands/ocicmds/artifacts/tag"
	"github.com/open-component-model/ocm/cmds/ocm/commands/verbs"
	"github.com/open-component-model/ocm/cmds/ocm/pkg/utils"
	"github.com/open-component-model/ocm/pkg/contexts/clictx"
)

// NewCommand creates a new command.
func NewCommand(ctx clictx.Context) *cobra.Command {
	cmd := utils.MassageCommand(&cobra.Command{
		Short: "tag OCI artifacts",
	}, verbs.Tag)
	cmd.AddCommand(artifacts.NewCommand(ctx))
	return cmd
}
