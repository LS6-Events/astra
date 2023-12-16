package astra

import (
	"github.com/ls6-events/astra/astTraversal"
)

type ContextFuncBuilder struct {
	Route           *Route
	Traverser       *astTraversal.CallExpressionTraverser
	paramOperations []func() (any, error)
}

func NewContextFuncBuilder(route *Route, traverser *astTraversal.CallExpressionTraverser) *ContextFuncBuilder {
	paramOperations := make([]func() (any, error), 0)
	return &ContextFuncBuilder{
		Route:           route,
		Traverser:       traverser,
		paramOperations: paramOperations,
	}
}

func (c *ContextFuncBuilder) getCurrentParamIndex() int {
	return len(c.paramOperations)
}

func (c *ContextFuncBuilder) Ignored() *ContextFuncBuilder {
	c.paramOperations = append(c.paramOperations, func() (any, error) {
		return nil, nil
	})

	return c
}

func (c *ContextFuncBuilder) StatusCode() *ContextFuncBuilder {
	currIndex := c.getCurrentParamIndex()
	c.paramOperations = append(c.paramOperations, func() (any, error) {
		return c.Traverser.Traverser.ExtractStatusCode(c.Traverser.Node.Args[currIndex])
	})

	return c
}

func (c *ContextFuncBuilder) ExpressionResult() *ContextFuncBuilder {
	currIndex := c.getCurrentParamIndex()
	c.paramOperations = append(c.paramOperations, func() (any, error) {
		exprType, err := c.Traverser.Traverser.Expression(c.Traverser.Node.Args[currIndex]).Type()
		if err != nil {
			return nil, err
		}

		return c.Traverser.Traverser.Type(exprType, c.Traverser.File.Package).Result()
	})

	return c
}

func (c *ContextFuncBuilder) Value() *ContextFuncBuilder {
	currIndex := c.getCurrentParamIndex()
	c.paramOperations = append(c.paramOperations, func() (any, error) {
		expr := c.Traverser.Traverser.Expression(c.Traverser.Node.Args[currIndex])

		return expr.Value()
	})

	return c
}

func (c *ContextFuncBuilder) Build(mapper func(*Route, []any) (*Route, error)) (*Route, error) {
	results := make([]any, len(c.paramOperations))
	for i, operation := range c.paramOperations {
		result, err := operation()
		if err != nil {
			return nil, err
		}

		results[i] = result
	}

	return mapper(c.Route, results)
}

type CustomFunc func(contextVarName string, contextFuncBuilder *ContextFuncBuilder) (*Route, error)

type CustomFuncOption struct{}

func (o CustomFuncOption) With(customFunc CustomFunc) FunctionalOption {
	return func(s *Service) {
		s.CustomFuncs = append(s.CustomFuncs, customFunc)
	}
}

func (o CustomFuncOption) LoadFromPlugin(s *Service, p *ConfigurationPlugin) error {
	customFuncsSymbol, found := p.Lookup("CustomFuncs")
	if found {
		customFuncs, ok := customFuncsSymbol.([]CustomFunc)
		if !ok {
			return ErrInvalidTypeFromConfigurationFile
		}

		for _, customFunc := range customFuncs {
			o.With(customFunc)(s)
		}
	}

	return nil
}
