package astra

type FunctionalOption func(*Service)

// PluginOption We separate this option from the others to not have to deal with the generics being the same for a list
type PluginOption interface {
	LoadFromPlugin(*Service, *ConfigurationPlugin) error
}

var pluginOptionsList []PluginOption

func RegisterOption(plugin PluginOption) {
	pluginOptionsList = append(pluginOptionsList, plugin)
}
