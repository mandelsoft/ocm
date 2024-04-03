// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package tag

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/open-component-model/ocm/cmds/ocm/commands/ocicmds/common"
	"github.com/open-component-model/ocm/cmds/ocm/commands/ocicmds/common/handlers/artifacthdlr"
	"github.com/open-component-model/ocm/cmds/ocm/commands/ocicmds/common/options/repooption"
	"github.com/open-component-model/ocm/cmds/ocm/commands/ocicmds/names"
	"github.com/open-component-model/ocm/cmds/ocm/commands/verbs"
	"github.com/open-component-model/ocm/cmds/ocm/pkg/output"
	"github.com/open-component-model/ocm/cmds/ocm/pkg/utils"
	"github.com/open-component-model/ocm/pkg/contexts/clictx"
	"github.com/open-component-model/ocm/pkg/contexts/oci"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/repositories/genericocireg/componentmapping"
	"github.com/open-component-model/ocm/pkg/out"
)

var (
	Names = names.Artifacts
	Verb  = verbs.Tag
)

type Command struct {
	utils.BaseCommand

	OCMPrefix string

	Tag  string
	Refs []string
}

// NewCommand creates a new artifact command.
func NewCommand(ctx clictx.Context, names ...string) *cobra.Command {
	return utils.SetupCommand(&Command{BaseCommand: utils.NewBaseCommand(ctx, repooption.New())}, utils.Names(Names, names...)...)
}

func (o *Command) ForName(name string) *cobra.Command {
	return &cobra.Command{
		Use:   "[<options>] {<artifact-reference>}",
		Short: "tag artifact version",
		Long: `
Sdd a tag to given artifact references.
	`,
		Example: `
$ ocm tag artifact latest ghcr.io/mandelsoft/kubelink:v1 
`,
	}
}

func (o *Command) AddFlags(flags *pflag.FlagSet) {
	o.BaseCommand.AddFlags(flags)
	// TODO after label-based discovery has been implemented for OCM
	// flags.StringVarP(&o.OCMPrefix, "ocm-repository-prefix", "O", "", "use as ocm repository with prefix (and component versions as refs)")
}

func (o *Command) Complete(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("a tag and at least one artifact reference required")
	}
	o.Tag = args[0]

	if o.OCMPrefix != "" {
		for _, r := range args[1:] {
			o.Refs = append(o.Refs, componentmapping.MapComponentName(o.OCMPrefix, r))
		}
	} else {
		o.Refs = args[1:]
	}
	return nil
}

func (o *Command) Run() error {
	session := oci.NewSession(nil)
	defer session.Close()
	err := o.ProcessOnOptions(common.CompleteOptionsWithContext(o.Context, session))
	if err != nil {
		return err
	}
	handler := artifacthdlr.NewTypeHandler(o.Context.OCI(), session, repooption.From(o).Repository)
	return utils.HandleOutput(&action{cmd: o}, handler, utils.StringElemSpecs(o.Refs...)...)
}

type action struct {
	cmd    *Command
	ok     int
	failed int
}

var _ output.Output = (*action)(nil)

func (a *action) Add(e interface{}) error {
	p := e.(*artifacthdlr.Object)

	err := p.Namespace.AddTags(p.Artifact.Digest(), a.cmd.Tag)
	if err == nil {
		a.ok++
		out.Outf(a.cmd, "tagged %s@%s\n", p.Namespace.GetNamespace(), p.Artifact.Digest())
	} else {
		a.failed++
		out.Outf(a.cmd, "failed for %s@%s: %s\n", p.Namespace.GetNamespace(), p.Artifact.Digest(), err.Error())
	}
	return nil
}

func (a *action) Close() error {
	return nil
}

func (a *action) Out() error {
	if a.failed > 0 {
		return fmt.Errorf("failed for %d artifacts", a.failed)
	}
	return nil
}
