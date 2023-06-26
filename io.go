package gengo

type ServiceFunction func(service *Service) error

type IOConfigurationKey string

const (
	IOConfigurationKeyFilePath IOConfigurationKey = "filePath"
)

type IOConfiguration map[IOConfigurationKey]any
