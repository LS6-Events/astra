package gengo

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

func (s *Service) SetConfig(config *Config) *Service {
	s.Config = config
	return s
}
