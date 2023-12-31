package astra

import (
	"github.com/ls6-events/astra/astTraversal"
	"github.com/stretchr/testify/require"
	"go/ast"
	"testing"
)

func TestNewContextFuncBuilder(t *testing.T) {
	route := &Route{}
	traverser := &astTraversal.CallExpressionTraverser{}

	contextFuncBuilder := NewContextFuncBuilder(route, traverser)

	require.NotNil(t, contextFuncBuilder)
	require.Equal(t, route, contextFuncBuilder.Route)
	require.Equal(t, traverser, contextFuncBuilder.Traverser)
}

func TestContextFuncBuilder_getCurrentParamIndex(t *testing.T) {
	contextFuncBuilder := &ContextFuncBuilder{
		paramOperations: []func() (any, error){},
	}

	require.Equal(t, 0, contextFuncBuilder.getCurrentParamIndex())

	contextFuncBuilder.paramOperations = append(contextFuncBuilder.paramOperations, func() (any, error) {
		return nil, nil
	})

	require.Equal(t, 1, contextFuncBuilder.getCurrentParamIndex())
}

func TestContextFuncBuilder_Ignored(t *testing.T) {
	contextFuncBuilder := &ContextFuncBuilder{
		paramOperations: []func() (any, error){},
	}

	require.Empty(t, contextFuncBuilder.paramOperations)

	contextFuncBuilder.Ignored()

	require.NotEmpty(t, contextFuncBuilder.paramOperations)
}

func TestContextFuncBuilder_StatusCode(t *testing.T) {
	contextFuncBuilder := &ContextFuncBuilder{
		paramOperations: []func() (any, error){},
	}

	require.Empty(t, contextFuncBuilder.paramOperations)

	contextFuncBuilder.StatusCode()

	require.NotEmpty(t, contextFuncBuilder.paramOperations)
}

func TestContextFuncBuilder_ExpressionResult(t *testing.T) {
	contextFuncBuilder := &ContextFuncBuilder{
		paramOperations: []func() (any, error){},
	}

	require.Empty(t, contextFuncBuilder.paramOperations)

	contextFuncBuilder.ExpressionResult()

	require.NotEmpty(t, contextFuncBuilder.paramOperations)
}

func TestContextFuncBuilder_Value(t *testing.T) {
	contextFuncBuilder := &ContextFuncBuilder{
		paramOperations: []func() (any, error){},
	}

	require.Empty(t, contextFuncBuilder.paramOperations)

	contextFuncBuilder.Value()

	require.NotEmpty(t, contextFuncBuilder.paramOperations)
}

func setupTestCallExpressionTraverser(t *testing.T, funcName string) *astTraversal.CallExpressionTraverser {
	t.Helper()

	traverser, err := astTraversal.CreateTraverserFromTestFile("callExpression.go")
	require.NoError(t, err)

	var function *ast.FuncDecl
	ast.Inspect(traverser.ActiveFile().AST, func(node ast.Node) bool {
		if node == nil {
			return false
		}

		funcDecl, ok := node.(*ast.FuncDecl)
		if ok && funcDecl.Name.Name == "contextFuncBuilderTest" {
			function = funcDecl
			return false
		}

		return true
	})

	for i, funcExpr := range function.Body.List {
		assignStmt, ok := funcExpr.(*ast.AssignStmt)
		if !ok {
			continue
		}

		callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr)
		if !ok {
			continue
		}

		ident, ok := callExpr.Fun.(*ast.Ident)
		if !ok {
			continue
		}

		if ident.Name == funcName {
			callExpression, err := traverser.CallExpression(function.Body.List[i].(*ast.AssignStmt).Rhs[0])
			require.NoError(t, err)

			return callExpression
		}
	}

	t.Fatal("could not find function")
	return nil
}

func TestContextFuncBuilder_Build(t *testing.T) {
	t.Run("build can manipulate route", func(t *testing.T) {
		callExpression := setupTestCallExpressionTraverser(t, "contextFuncBuilderIgnored")

		contextFuncBuilder := NewContextFuncBuilder(&Route{
			Method: "GET",
			Path:   "/",
		}, callExpression)

		route, err := contextFuncBuilder.Build(func(route *Route, params []any) (*Route, error) {
			route.Method = "POST"
			route.Path = "/new"
			return route, nil
		})
		require.NoError(t, err)
		require.NotNil(t, route)
		require.Equal(t, "POST", route.Method)
		require.Equal(t, "/new", route.Path)
	})

	t.Run("valid ignored param", func(t *testing.T) {
		callExpression := setupTestCallExpressionTraverser(t, "contextFuncBuilderIgnored")

		contextFuncBuilder := NewContextFuncBuilder(&Route{
			Method: "GET",
			Path:   "/",
		}, callExpression)

		contextFuncBuilder.Ignored()

		route, err := contextFuncBuilder.Build(func(route *Route, params []any) (*Route, error) {
			return route, nil
		})
		require.NoError(t, err)
		require.NotNil(t, route)
	})

	t.Run("valid status code param", func(t *testing.T) {
		callExpression := setupTestCallExpressionTraverser(t, "contextFuncBuilderStatusCode")

		_, err := callExpression.Traverser.Packages.Get(callExpression.Traverser.Packages.Find("github.com/ls6-events/astra/astTraversal/testfiles"))
		require.NoError(t, err)

		contextFuncBuilder := NewContextFuncBuilder(&Route{
			Method: "GET",
			Path:   "/",
		}, callExpression)

		contextFuncBuilder.StatusCode()

		route, err := contextFuncBuilder.Build(func(route *Route, params []any) (*Route, error) {
			require.Equal(t, 200, params[0])
			return route, nil
		})
		require.NoError(t, err)
		require.NotNil(t, route)
	})

	t.Run("valid expression result param", func(t *testing.T) {
		callExpression := setupTestCallExpressionTraverser(t, "contextFuncBuilderExpressionResult")

		_, err := callExpression.Traverser.Packages.Get(callExpression.Traverser.Packages.Find("github.com/ls6-events/astra/astTraversal/testfiles"))
		require.NoError(t, err)

		contextFuncBuilder := NewContextFuncBuilder(&Route{
			Method: "GET",
			Path:   "/",
		}, callExpression)

		contextFuncBuilder.ExpressionResult()

		route, err := contextFuncBuilder.Build(func(route *Route, params []any) (*Route, error) {
			param := params[0].(astTraversal.Result)
			require.Equal(t, "MyStruct", param.Type)

			return route, nil
		})

		require.NoError(t, err)
		require.NotNil(t, route)
	})

	t.Run("valid value param", func(t *testing.T) {
		callExpression := setupTestCallExpressionTraverser(t, "contextFuncBuilderValue")

		pkgNode := callExpression.Traverser.Packages.AddPackage("github.com/ls6-events/astra/astTraversal/testfiles")

		_, err := callExpression.Traverser.Packages.Get(pkgNode)
		require.NoError(t, err)

		contextFuncBuilder := NewContextFuncBuilder(&Route{
			Method: "GET",
			Path:   "/",
		}, callExpression)

		contextFuncBuilder.Value()

		route, err := contextFuncBuilder.Build(func(route *Route, params []any) (*Route, error) {
			require.Equal(t, "bar", params[0])

			return route, nil
		})
		require.NoError(t, err)
		require.NotNil(t, route)
	})
}
