package gengo

type Config struct {
	Title       string
	Description string
	Version     string
	Contact     Contact
	License     License

	Secure   bool
	Host     string
	BasePath string
	Port     int
}

type Contact struct {
	Name  string
	URL   string
	Email string
}

type License struct {
	Name string
	URL  string
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
