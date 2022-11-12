// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"encoding/json"
	"os"

	"github.com/mandelsoft/filepath/pkg/filepath"
	"github.com/open-component-model/ocm/pkg/contexts/ocm/plugin/ppi"
	"github.com/open-component-model/ocm/pkg/errors"
)

type Config struct {
	Systems map[string]Values `json:"systems"`
}

type Values struct {
	Address string `json:"address"`
	Path    string `json:"path,omitempty"`
}

func GetConfig(raw json.RawMessage) (interface{}, error) {
	var cfg Config

	err := json.Unmarshal(raw, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func BasePath(p ppi.Plugin, name string) (string, string, error) {
	root := filepath.Join(os.TempDir(), "r3trans", name)
	path := ""
	cfg, _ := p.GetConfig()
	if cfg != nil {
		cfg := cfg.(*Config)
		if cfg.Systems != nil {
			base := cfg.Systems[name]
			if base.Address != "" {
				root = base.Address
			}
			path = base.Path
		}
	}
	err := os.MkdirAll(root, 0o700)
	if err != nil {
		return "", "", errors.Wrapf(err, "cannot create root dir")
	}
	return root, path, nil
}
