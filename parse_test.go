package astra

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_SetupParse(t *testing.T) {
	t.Run("returns error if no inputs are provided", func(t *testing.T) {
		service := &Service{}

		err := service.SetupParse()
		require.Error(t, err)
		require.ErrorContains(t, err, "input not set")
	})

	t.Run("returns error if no outputs are provided", func(t *testing.T) {
		service := &Service{
			Inputs: []Input{
				{
					Mode:         "",
					CreateRoutes: nil,
					ParseRoutes:  nil,
				},
			},
		}

		err := service.SetupParse()
		require.Error(t, err)
		require.ErrorContains(t, err, "output not set")
	})

	t.Run("create routes is called", func(t *testing.T) {
		var createRoutesCalled bool
		var parseRoutesCalled bool
		service := &Service{
			Inputs: []Input{
				{
					Mode: "",
					CreateRoutes: func(service *Service) error {
						createRoutesCalled = true
						return nil
					},
					ParseRoutes: func(service *Service) error {
						parseRoutesCalled = true
						return nil
					},
				},
			},
			Outputs: []Output{
				{
					Mode:     "",
					Generate: nil,
				},
			},
		}

		err := service.SetupParse()
		require.NoError(t, err)

		require.True(t, createRoutesCalled)
		require.False(t, parseRoutesCalled)
	})
}

func TestService_CompleteParse(t *testing.T) {
	t.Run("returns error if no inputs are provided", func(t *testing.T) {
		service := &Service{}

		err := service.Parse()
		require.Error(t, err)
		require.ErrorContains(t, err, "input not set")
	})

	t.Run("returns error if no outputs are provided", func(t *testing.T) {
		service := &Service{
			Inputs: []Input{
				{
					Mode:         "",
					CreateRoutes: nil,
					ParseRoutes:  nil,
				},
			},
		}

		err := service.Parse()
		require.Error(t, err)
		require.ErrorContains(t, err, "output not set")
	})

	t.Run("create routes, parse routes and generate are called", func(t *testing.T) {
		var createRoutesCalled bool
		var parseRoutesCalled bool

		var generateCalled bool

		service := &Service{
			Inputs: []Input{
				{
					Mode: "",
					CreateRoutes: func(service *Service) error {
						createRoutesCalled = true
						return nil
					},
					ParseRoutes: func(service *Service) error {
						parseRoutesCalled = true
						return nil
					},
				},
			},
			Outputs: []Output{
				{
					Mode: "",
					Generate: func(service *Service) error {
						generateCalled = true
						return nil
					},
				},
			},
		}

		err := service.Setup()
		require.NoError(t, err)

		err = service.CompleteParse()
		require.NoError(t, err)

		require.False(t, createRoutesCalled)
		require.True(t, parseRoutesCalled)

		require.True(t, generateCalled)
	})
}

func TestService_Parse(t *testing.T) {
	t.Run("returns error if no inputs are provided", func(t *testing.T) {
		service := &Service{}

		err := service.CompleteParse()
		require.Error(t, err)
		require.ErrorContains(t, err, "input not set")
	})

	t.Run("returns error if no outputs are provided", func(t *testing.T) {
		service := &Service{
			Inputs: []Input{
				{
					Mode:         "",
					CreateRoutes: nil,
					ParseRoutes:  nil,
				},
			},
		}

		err := service.CompleteParse()
		require.Error(t, err)
		require.ErrorContains(t, err, "output not set")
	})

	t.Run("parse routes and generate are called", func(t *testing.T) {
		var createRoutesCalled bool
		var parseRoutesCalled bool

		var generateCalled bool

		service := &Service{
			Inputs: []Input{
				{
					Mode: "",
					CreateRoutes: func(service *Service) error {
						createRoutesCalled = true
						return nil
					},
					ParseRoutes: func(service *Service) error {
						parseRoutesCalled = true
						return nil
					},
				},
			},
			Outputs: []Output{
				{
					Mode: "",
					Generate: func(service *Service) error {
						generateCalled = true
						return nil
					},
				},
			},
		}

		err := service.Setup()
		require.NoError(t, err)

		err = service.Parse()
		require.NoError(t, err)

		require.True(t, createRoutesCalled)

		require.True(t, parseRoutesCalled)

		require.True(t, generateCalled)
	})
}
