package astra

// Config is the configuration for the generator.
// It matches very closely to the OpenAPI specification.
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

// Validate validates the configuration.
// It will also set default values for some fields.
func (c *Config) Validate() error {
	// For now, only the port is required, the rest are set to default values.
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

// SetConfig sets the configuration for the generator.
func (s *Service) SetConfig(config *Config) *Service {
	s.Config = config
	return s
}

// WithConfig sets the configuration for the generator in option pattern.
func WithConfig(config *Config) Option {
	return func(s *Service) {
		s.SetConfig(config)
	}
}
