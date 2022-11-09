// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package install

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"sigs.k8s.io/yaml"

	ocmcommon "github.com/open-component-model/ocm/cmds/ocm/commands/ocmcmds/common"
	"github.com/open-component-model/ocm/cmds/ocm/commands/ocmcmds/common/options/repooption"
	"github.com/open-component-model/ocm/cmds/ocm/commands/ocmcmds/names"
	"github.com/open-component-model/ocm/cmds/ocm/commands/verbs"
	"github.com/open-component-model/ocm/cmds/ocm/pkg/utils"
	"github.com/open-component-model/ocm/pkg/common"
	"github.com/open-component-model/ocm/pkg/contexts/clictx"
	"github.com/open-component-model/ocm/pkg/contexts/ocm"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/attrs/plugindirattr"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/download"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/plugin/cache"
	"github.com/open-component-model/ocm/pkg/errors"
	"github.com/open-component-model/ocm/pkg/out"
)

var (
	Names = names.Plugins
	Verb  = verbs.Install
)

type Command struct {
	utils.BaseCommand

	Describe bool
	Force    bool

	Ref  string
	Name string
}

// NewCommand creates a new ctf command.
func NewCommand(ctx clictx.Context, names ...string) *cobra.Command {
	return utils.SetupCommand(
		&Command{
			BaseCommand: utils.NewBaseCommand(ctx, repooption.New()),
		},
		utils.Names(Names, names...)...,
	)
}

func (o *Command) ForName(name string) *cobra.Command {
	return &cobra.Command{
		Use:   "[<options>] <component version> <plugin name>",
		Short: "get plugins",
		Long: `
Download and install a plugin provided by an OCM component version.
`,
		Args: cobra.RangeArgs(1, 2),
		Example: `
$ ocm install plugin ghcr.io/github.com/mandelsoft/cnudie//github.com/mandelsoft/ocmplugin:0.1.0-dev
`,
	}
}
func (o *Command) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVarP(&o.Describe, "describe", "d", false, "describe plugin, only")
	fs.BoolVarP(&o.Force, "force", "f", false, "overwrite existing plugin")
}

func (o *Command) Complete(args []string) error {
	o.Ref = args[0]
	if len(args) > 1 {
		o.Name = args[2]
	}
	return nil
}

func (o *Command) Run() error {
	session := ocm.NewSession(nil)
	defer session.Close()

	err := o.ProcessOnOptions(ocmcommon.CompleteOptionsWithSession(o.Context, session))
	if err != nil {
		return err
	}

	repo := repooption.From(o)
	if repo.Repository != nil {
		return o.downloadFromRepo(session, repo.Repository, o.Ref, o.Name)
	}
	return o.downloadRef(session, o.Ref, o.Name)
}

/////////////////////////////////////////////////////////////////////////////

func (o *Command) downloadRef(session ocm.Session, ref string, name string) error {
	result, err := session.EvaluateVersionRef(o.Context.OCMContext(), ref)
	if err != nil {
		return err
	}
	if result.Version == nil {
		return fmt.Errorf("component version required")
	}
	return o.download(session, result.Version, name)
}

func (o *Command) downloadFromRepo(session ocm.Session, repo ocm.Repository, ref, name string) error {
	cr, err := ocm.ParseComp(ref)
	if err != nil {
		return err
	}

	var cv ocm.ComponentVersionAccess
	comp, err := session.LookupComponent(repo, cr.Component)
	if err != nil {
		return err
	}
	if cr.IsVersion() {
		cv, err = session.GetComponentVersion(comp, *cr.Version)
	} else {
		var vers []string

		vers, err = comp.ListVersions()
		if err != nil {
			return errors.Wrapf(err, "cannot list versions for component %s", cr.Component)
		}
		if len(vers) > 1 {
			return errors.Wrapf(err, "nultiple versions found for component %s", cr.Component)
		}
		if len(vers) == 0 {
			return errors.Wrapf(err, "no versions found for component %s", cr.Component)
		}
		cv, err = session.GetComponentVersion(comp, vers[0])
	}
	if err != nil {
		return err
	}
	return o.download(session, cv, name)
}

func (o *Command) download(session ocm.Session, cv ocm.ComponentVersionAccess, name string) (err error) {
	defer errors.PropagateErrorf(&err, nil, "%s", common.VersionedElementKey(cv))

	var found ocm.ResourceAccess
	var wrong ocm.ResourceAccess

	for _, r := range cv.GetResources() {
		if name != "" && r.Meta().Name != name {
			continue
		}
		if r.Meta().Type == "ocmPlugin" {
			if r.Meta().ExtraIdentity.Get("os") == runtime.GOOS &&
				r.Meta().ExtraIdentity.Get("architecture") == runtime.GOARCH {
				found = r
				break
			}
			wrong = r
		} else {
			if name != "" {
				wrong = r
			}
		}
	}
	if found == nil {
		if wrong != nil {
			if wrong.Meta().Type != "ocmPlugin" {
				return fmt.Errorf("resource %q has wrong type: %s", wrong.Meta().Name, wrong.Meta().Type)
			}
			return fmt.Errorf("os %s architecture %s not found for resource %q", runtime.GOOS, runtime.GOARCH, wrong.Meta().Name)
		}
		if name != "" {
			return fmt.Errorf("resource %q not found", name)
		}
		return fmt.Errorf("no ocmPlugin found")
	}
	out.Outf(o.Context, "found resource %s[%s]\n", found.Meta().Name, found.Meta().ExtraIdentity.String())

	file, err := os.CreateTemp(os.TempDir(), "plugin-*")
	if err != nil {
		return errors.Wrapf(err, "cannot create temp file")
	}
	file.Close()
	fs := osfs.New()
	_, _, err = download.For(o.Context).Download(o, found, file.Name(), fs)
	if err != nil {
		return errors.Wrapf(err, "cannot download resource %s", found.Meta().Name)
	}

	desc, err := cache.GetPluginInfo(file.Name())
	if err != nil {
		return err
	}
	if o.Describe {
		data, err := yaml.Marshal(desc)
		if err != nil {
			return errors.Wrapf(err, "cannot marshal plugin descriptor")
		}
		out.Outln(o.Context, string(data))
	} else {
		dir := plugindirattr.Get(o.Context)
		if dir != "" {
			target := filepath.Join(dir, desc.PluginName)

			if ok, _ := vfs.FileExists(fs, target); ok {
				if !o.Force {
					return fmt.Errorf("plugin %s already found is %s", desc.PluginName, dir)
				}
				fs.Remove(target)
			}
			out.Outf(o, "installing plugin %s[%s] in %s...\n", desc.PluginName, desc.PluginVersion, dir)
			dst, err := fs.OpenFile(target, vfs.O_CREATE|vfs.O_TRUNC|vfs.O_WRONLY, 0o755)
			if err != nil {
				return errors.Wrapf(err, "cannot create plugin file %s", target)
			}
			defer dst.Close()
			src, err := fs.OpenFile(file.Name(), vfs.O_RDONLY, 0)
			if err != nil {
				return errors.Wrapf(err, "cannot open plugin executable %s", file.Name())
			}
			defer src.Close()
			_, err = io.Copy(dst, src)
			if err != nil {
				return errors.Wrapf(err, "cannot copy plugin file %s", target)
			}
		}
	}

	return nil
}
