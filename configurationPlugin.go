package astra

import (
	"errors"
	"plugin"
)

const CONFIGURATION_FUNCTION_PREFIX = "Astra_"

var (
	ErrInvalidTypeFromConfigurationFile = errors.New("invalid type from configuration file")
)

type ConfigurationPlugin struct {
	plugin *plugin.Plugin
}

func NewConfigurationPlugin(plugin *plugin.Plugin) *ConfigurationPlugin {
	return &ConfigurationPlugin{
		plugin: plugin,
	}
}

func (p *ConfigurationPlugin) Lookup(name string) (plugin.Symbol, bool) {
	symbol, err := p.plugin.Lookup(CONFIGURATION_FUNCTION_PREFIX + name)
	if err != nil {
		return nil, false
	}

	return symbol, true
}
