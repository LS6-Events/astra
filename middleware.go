package astra

func UnstableWithMiddleware() Option {
	return func(service *Service) {
		service.UnstableEnableMiddleware = true
	}
}
