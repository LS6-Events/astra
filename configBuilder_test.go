package astra

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewConfigBuilder(t *testing.T) {
	configBuilder := NewConfigBuilder()

	require.NotNil(t, configBuilder)

	require.NotNil(t, configBuilder.config)
}

func TestConfigBuilder_SetTitle(t *testing.T) {
	configBuilder := NewConfigBuilder()

	require.Empty(t, configBuilder.config.Title)

	configBuilder.SetTitle("test")

	require.Equal(t, "test", configBuilder.config.Title)
}

func TestConfigBuilder_SetDescription(t *testing.T) {
	configBuilder := NewConfigBuilder()

	require.Empty(t, configBuilder.config.Description)

	configBuilder.SetDescription("test")

	require.Equal(t, "test", configBuilder.config.Description)
}

func TestConfigBuilder_SetVersion(t *testing.T) {
	configBuilder := NewConfigBuilder()

	require.Empty(t, configBuilder.config.Version)

	configBuilder.SetVersion("test")

	require.Equal(t, "test", configBuilder.config.Version)
}

func TestConfigBuilder_SetContact(t *testing.T) {
	configBuilder := NewConfigBuilder()

	require.Empty(t, configBuilder.config.Contact)

	configBuilder.SetContact(Contact{
		Name:  "John Doe",
		URL:   "https://www.google.com",
		Email: "john@doe.com",
	})

	require.Equal(t, Contact{
		Name:  "John Doe",
		URL:   "https://www.google.com",
		Email: "john@doe.com",
	}, configBuilder.config.Contact)
}

func TestConfigBuilder_SetLicense(t *testing.T) {
	configBuilder := NewConfigBuilder()

	require.Empty(t, configBuilder.config.License)

	configBuilder.SetLicense(License{
		Name: "MIT",
		URL:  "https://opensource.org/licenses/MIT",
	})

	require.Equal(t, License{
		Name: "MIT",
		URL:  "https://opensource.org/licenses/MIT",
	}, configBuilder.config.License)
}

func TestConfigBuilder_SetSecure(t *testing.T) {
	configBuilder := NewConfigBuilder()

	require.Empty(t, configBuilder.config.Secure)

	configBuilder.SetSecure(true)

	require.Equal(t, true, configBuilder.config.Secure)
}

func TestConfigBuilder_SetHost(t *testing.T) {
	configBuilder := NewConfigBuilder()

	require.Empty(t, configBuilder.config.Host)

	configBuilder.SetHost("localhost")

	require.Equal(t, "localhost", configBuilder.config.Host)
}

func TestConfigBuilder_SetBasePath(t *testing.T) {
	configBuilder := NewConfigBuilder()

	require.Empty(t, configBuilder.config.BasePath)

	configBuilder.SetBasePath("/base-path")

	require.Equal(t, "/base-path", configBuilder.config.BasePath)
}

func TestConfigBuilder_SetPort(t *testing.T) {
	configBuilder := NewConfigBuilder()

	require.Empty(t, configBuilder.config.Port)

	configBuilder.SetPort(8000)

	require.Equal(t, 8000, configBuilder.config.Port)
}

func TestConfigBuilder_Build(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		configBuilder := NewConfigBuilder()

		configBuilder.SetPort(8000)

		config, err := configBuilder.Build()
		require.NoError(t, err)

		require.NotNil(t, config)
	})

	t.Run("invalid config", func(t *testing.T) {
		configBuilder := NewConfigBuilder()

		config, err := configBuilder.Build()
		require.Error(t, err)
		require.Equal(t, ErrConfigPortRequired, err)

		require.Nil(t, config)
	})
}
