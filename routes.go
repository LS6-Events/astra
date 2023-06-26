package gengo

func (s *Service) AddRoute(route Route) {
	s.Routes = append(s.Routes, route)
}

func (s *Service) ReplaceRoute(route Route) {
	for i, r := range s.Routes {
		if r.Path == route.Path && r.Method == route.Method {
			s.Routes[i] = route
			return
		}
	}
}
