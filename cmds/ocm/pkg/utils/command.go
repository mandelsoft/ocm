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

package utils

import (
	"os"
	"strings"

	"github.com/gardener/ocm/cmds/ocm/clictx"
	"github.com/gardener/ocm/cmds/ocm/pkg/output/out"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// OCMCommand is a command pattern, thta can be instantiated for a dediated
// sub command name.
type OCMCommand interface {
	clictx.Context
	ForName(name string) *cobra.Command
	AddFlags(fs *pflag.FlagSet)
	Complete(args []string) error
	Run() error
}

type BaseCommand struct {
	clictx.Context
}

func NewBaseCommand(ctx clictx.Context) BaseCommand {
	return BaseCommand{Context: ctx}
}

func (BaseCommand) AddFlags(fs *pflag.FlagSet)   {}
func (BaseCommand) Complete(args []string) error { return nil }

func SetupCommand(ocmcmd OCMCommand, names ...string) *cobra.Command {
	c := ocmcmd.ForName(names[0])
	if !strings.HasSuffix(c.Use, names[0]+" ") {
		c.Use = names[0] + " " + c.Use
	}
	c.Aliases = names[1:]
	c.Run = func(cmd *cobra.Command, args []string) {
		if err := ocmcmd.Complete(args); err != nil {
			out.Error(ocmcmd, err.Error())
			os.Exit(1)
		}

		if err := ocmcmd.Run(); err != nil {
			out.Error(ocmcmd, err.Error())
			os.Exit(1)
		}
	}
	c.TraverseChildren = true
	ocmcmd.AddFlags(c.Flags())
	return c
}