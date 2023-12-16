package astra

// Config is the configuration for the generator
// It matches very closely to the OpenAPI specification
type Config struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Version     string  `json:"version"`
	Contact     Contact `json:"contact"`
	License     License `json:"license"`

	Secure   bool   `json:"secure"`
	Host     string `json:"host"`
	BasePath string `json:"basePath"`
	Port     int    `json:"port"`
}

type Contact struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Email string `json:"email"`
}

type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Validate validates the configuration
// It will also set default values for some fields
func (c *Config) Validate() error {
	// For now, only the port is required, the rest are set to default values
	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.Port == 0 {
		return ErrConfigPortRequired
	}
	if c.BasePath == "" {
		c.BasePath = "/"
	}

	return nil
}

// SetConfig sets the configuration for the generator
func (s *Service) SetConfig(config *Config) *Service {
	s.Config = config
	return s
}

type ConfigOption struct{}

// WithConfig sets the configuration for the generator in option pattern
func (o ConfigOption) With(config *Config) FunctionalOption {
	return func(s *Service) {
		s.SetConfig(config)
	}
}

func (o ConfigOption) LoadFromPlugin(s *Service, p *ConfigurationPlugin) error {
	configSymbol, found := p.Lookup("Config")
	if found {
		config, ok := configSymbol.(*Config)
		if !ok {
			return ErrInvalidTypeFromConfigurationFile
		} else {
			o.With(config)(s)
		}
	} else {
		config := &Config{}

		title, found := p.Lookup("Title")
		if found {
			if title, ok := title.(string); ok {
				config.Title = title
			}
		}

		description, found := p.Lookup("Description")
		if found {
			if description, ok := description.(string); ok {
				config.Description = description
			}
		}

		version, found := p.Lookup("Version")
		if found {
			if version, ok := version.(string); ok {
				config.Version = version
			}
		}

		contact, found := p.Lookup("Contact")
		if found {
			if contact, ok := contact.(Contact); ok {
				config.Contact = contact
			}
		}

		license, found := p.Lookup("License")
		if found {
			if license, ok := license.(License); ok {
				config.License = license
			}
		}

		secure, found := p.Lookup("Secure")
		if found {
			if secure, ok := secure.(bool); ok {
				config.Secure = secure
			}
		}

		host, found := p.Lookup("Host")
		if found {
			if host, ok := host.(string); ok {
				config.Host = host
			}
		}

		basePath, found := p.Lookup("BasePath")
		if found {
			if basePath, ok := basePath.(string); ok {
				config.BasePath = basePath
			}
		}

		port, found := p.Lookup("Port")
		if found {
			if port, ok := port.(int); ok {
				config.Port = port
			}
		}

		o.With(config)(s)
	}

	return nil
}

func init() {
	RegisterOption(ConfigOption{})
}
