package astra

// AddRoute adds a route to the service.
func (s *Service) AddRoute(route Route) {
	s.Routes = append(s.Routes, route)
}

// ReplaceRoute replaces a route in the service using the path and method as indexes.
func (s *Service) ReplaceRoute(route Route) {
	for i, r := range s.Routes {
		if r.Path == route.Path && r.Method == route.Method {
			s.Routes[i] = route
			return
		}
	}
}
