// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/open-component-model/ocm/pkg/contexts/config"
	cfgcpi "github.com/open-component-model/ocm/pkg/contexts/config/cpi"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/cpi"
	"github.com/open-component-model/ocm/pkg/runtime"
)

const (
	ConfigType   = "ocm" + cfgcpi.OCM_CONFIG_TYPE_SUFFIX
	ConfigTypeV1 = ConfigType + runtime.VersionSeparator + "v1"
)

func init() {
	cfgcpi.RegisterConfigType(cfgcpi.NewConfigType(ConfigType, &Config{}))
	cfgcpi.RegisterConfigType(cfgcpi.NewConfigType(ConfigTypeV1, &Config{}))
}

// Config describes a memory based config interface.
type Config struct {
	runtime.ObjectVersionedType `json:",inline"`
	Aliases                     map[string]*cpi.GenericRepositorySpec `json:"aliases,omitempty"`
}

// New creates a new memory ConfigSpec.
func New() *Config {
	return &Config{
		ObjectVersionedType: runtime.NewVersionedTypedObject(ConfigType),
	}
}

func (a *Config) GetType() string {
	return ConfigType
}

func (a *Config) SetAlias(name string, spec cpi.RepositorySpec) error {
	g, err := cpi.ToGenericRepositorySpec(spec)
	if err != nil {
		return err
	}
	if a.Aliases == nil {
		a.Aliases = map[string]*cpi.GenericRepositorySpec{}
	}
	a.Aliases[name] = g
	return nil
}

func (a *Config) ApplyTo(ctx config.Context, target interface{}) error {
	t, ok := target.(cpi.Context)
	if !ok {
		return config.ErrNoContext(ConfigType)
	}
	for n, s := range a.Aliases {
		t.SetAlias(n, s)
	}
	return nil
}
