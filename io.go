package astra

// ServiceFunction is the function used by the inputs/outputs to interact with the service.
type ServiceFunction func(service *Service) error

// IOConfigurationKey is the key for the IOConfiguration.
// IOConfiguration is used to configure the inputs/outputs with expected data that needs to be stored in JSON for caching.
// At the moment, this only includes a file path that is used by both the outputs.

type IOConfigurationKey string

const (
	IOConfigurationKeyFilePath      IOConfigurationKey = "filePath"
	IOConfigurationKeyDirectoryPath IOConfigurationKey = "directoryPath"
)

type IOConfiguration map[IOConfigurationKey]any
