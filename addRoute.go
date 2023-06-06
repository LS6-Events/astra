package gengo

func (s *Service) AddRoute(route Route) {
	s.Routes = append(s.Routes, route)
}
