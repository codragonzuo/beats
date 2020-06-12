// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package management

import (
	"github.com/gofrs/uuid"

	"github.com/codragonzuo/beats/libbeat/common"
	"github.com/codragonzuo/beats/libbeat/common/reload"
	"github.com/codragonzuo/beats/libbeat/feature"
)

// Namespace is the feature namespace for queue definition.
var Namespace = "libbeat.management"

// DebugK used as key for all things central management
var DebugK = "centralmgmt"

var centralMgmtKey = "x-pack-cm"

// ConfigManager interacts with the beat to update configurations
// from an external source
type ConfigManager interface {
	// Enabled returns true if config manager is enabled
	Enabled() bool

	// Start the config manager
	Start()

	// Stop the config manager
	Stop()

	// CheckRawConfig check settings are correct before launching the beat
	CheckRawConfig(cfg *common.Config) error
}

// PluginFunc for creating FactoryFunc if it matches a config
type PluginFunc func(*common.Config) FactoryFunc

// FactoryFunc for creating a config manager
type FactoryFunc func(*common.Config, *reload.Registry, uuid.UUID) (ConfigManager, error)

// Register a config manager
func Register(name string, fn PluginFunc, stability feature.Stability) {
	f := feature.New(Namespace, name, fn, feature.MakeDetails(name, "", stability))
	feature.MustRegister(f)
}

// Factory retrieves config manager constructor. If no one is registered
// it will create a nil manager
func Factory(cfg *common.Config) FactoryFunc {
	factories, err := feature.GlobalRegistry().LookupAll(Namespace)
	if err != nil {
		return nilFactory
	}

	for _, f := range factories {
		if plugin, ok := f.Factory().(PluginFunc); ok {
			if factory := plugin(cfg); factory != nil {
				return factory
			}
		}
	}

	return nilFactory
}

type modeConfig struct {
	Mode string `config:"mode" yaml:"mode"`
}

func defaultModeConfig() *modeConfig {
	return &modeConfig{
		Mode: centralMgmtKey,
	}
}

// nilManager, fallback when no manager is present
type nilManager struct{}

func nilFactory(*common.Config, *reload.Registry, uuid.UUID) (ConfigManager, error) {
	return nilManager{}, nil
}

func (nilManager) Enabled() bool                           { return false }
func (nilManager) Start()                                  {}
func (nilManager) Stop()                                   {}
func (nilManager) CheckRawConfig(cfg *common.Config) error { return nil }
