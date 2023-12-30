package astra

import "github.com/ls6-events/astra/astTraversal"

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
		//nolint:nilnil // This is intentional
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

func WithCustomFunc(customFunc CustomFunc) Option {
	return func(service *Service) {
		service.CustomFuncs = append(service.CustomFuncs, customFunc)
	}
}
