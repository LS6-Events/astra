package astra

// ConfigBuilder is a builder for the Config struct
// It has methods to set the fields of the Config struct
// It also has a Build method that returns the Config struct and an error, allowing for easy chaining and basic validation

type ConfigBuilder struct {
	config *Config
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: &Config{},
	}
}

func (c *ConfigBuilder) SetTitle(title string) *ConfigBuilder {
	c.config.Title = title
	return c
}

func (c *ConfigBuilder) SetDescription(description string) *ConfigBuilder {
	c.config.Description = description
	return c
}

func (c *ConfigBuilder) SetVersion(version string) *ConfigBuilder {
	c.config.Version = version
	return c
}

func (c *ConfigBuilder) SetContact(contact Contact) *ConfigBuilder {
	c.config.Contact = contact
	return c
}

func (c *ConfigBuilder) SetLicense(license License) *ConfigBuilder {
	c.config.License = license
	return c
}

func (c *ConfigBuilder) SetSecure(secure bool) *ConfigBuilder {
	c.config.Secure = secure
	return c
}

func (c *ConfigBuilder) SetHost(host string) *ConfigBuilder {
	c.config.Host = host
	return c
}

func (c *ConfigBuilder) SetBasePath(basePath string) *ConfigBuilder {
	c.config.BasePath = basePath
	return c
}

func (c *ConfigBuilder) SetPort(port int) *ConfigBuilder {
	c.config.Port = port
	return c
}

func (c *ConfigBuilder) Build() (*Config, error) {
	if err := c.config.Validate(); err != nil {
		return nil, err
	}
	return c.config, nil
}
