package astTraversal

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/token"
	"testing"
)

func TestExpressionTraverser_SetReturnNum(t *testing.T) {
	bt := &BaseTraverser{}
	et := bt.Expression(&ast.BasicLit{})
	assert.Equal(t, 0, et.ReturnNum)

	et.SetReturnNum(5)
	assert.Equal(t, 5, et.ReturnNum)
}

func TestExpressionTraverser_DoesNeedTracing(t *testing.T) {
	bt := &BaseTraverser{}
	// *ast.StarExpr test
	etStar := bt.Expression(&ast.StarExpr{
		X: &ast.BasicLit{},
	})
	assert.False(t, etStar.DoesNeedTracing())

	// *ast.UnaryExpr test
	etUnary := bt.Expression(&ast.UnaryExpr{
		X: &ast.BasicLit{},
	})
	assert.False(t, etUnary.DoesNeedTracing())

	// *ast.BasicLit test
	etBasic := bt.Expression(&ast.BasicLit{})
	assert.False(t, etBasic.DoesNeedTracing())

	// *ast.CompositeLit test
	etComposite := bt.Expression(&ast.CompositeLit{})
	assert.False(t, etComposite.DoesNeedTracing())

	// Default test
	etDefault := bt.Expression(&ast.FuncLit{})
	assert.True(t, etDefault.DoesNeedTracing())
}

func TestExpressionTraverser_Result(t *testing.T) {
	bt, err := createTraverserFromTestFile("usefulTypes.go")
	assert.NoError(t, err)

	t.Run("StarExpr", func(t *testing.T) {
		etStar := bt.Expression(&ast.StarExpr{
			X: &ast.BasicLit{
				Kind:  9, // token.STRING
				Value: "\"Hello\"",
			},
		})
		_, err := etStar.Result()
		assert.NoError(t, err)
	})

	t.Run("UnaryExpr", func(t *testing.T) {
		etUnary := bt.Expression(&ast.UnaryExpr{
			X: &ast.BasicLit{
				Kind:  5, // token.INT
				Value: "1234",
			},
		})
		_, err := etUnary.Result()
		assert.NoError(t, err)
	})

	t.Run("Ident", func(t *testing.T) {
		etIdent := bt.Expression(&ast.Ident{
			Name: "VariableName",
		})
		res, err := etIdent.Result()
		assert.NoError(t, err)
		assert.Equal(t, "VariableName", res.Type)
	})

	t.Run("SelectorExpr", func(t *testing.T) {
		etSelector := bt.Expression(&ast.SelectorExpr{
			X:   &ast.Ident{Name: "package"},
			Sel: &ast.Ident{Name: "VariableName"},
		})
		res, err := etSelector.Result()
		assert.NoError(t, err)
		assert.Equal(t, "VariableName", res.Type)
	})

	t.Run("ArrayType", func(t *testing.T) {
		etArray := bt.Expression(&ast.ArrayType{
			Elt: &ast.Ident{Name: "int"},
		})
		res, err := etArray.Result()
		assert.Nil(t, err)
		assert.Equal(t, "slice", res.Type)
		assert.Equal(t, "int", res.SliceType)
	})

	t.Run("MapType", func(t *testing.T) {
		etMap := bt.Expression(&ast.MapType{
			Key:   &ast.Ident{Name: "string"},
			Value: &ast.Ident{Name: "int"},
		})
		res, err := etMap.Result()
		assert.NoError(t, err)
		assert.Equal(t, "map", res.Type)
		assert.Equal(t, "int", res.MapValType)
		assert.Equal(t, "string", res.MapKeyType)
	})

	t.Run("CompositeLit", func(t *testing.T) {
		etComposite := bt.Expression(&ast.CompositeLit{
			Type: &ast.ArrayType{
				Elt: &ast.Ident{Name: "int"},
			},
		})
		res, err := etComposite.Result()
		assert.NoError(t, err)
		assert.Equal(t, "slice", res.Type)
		assert.Equal(t, "int", res.SliceType)
	})

	t.Run("BasicLit_String", func(t *testing.T) {
		etBasicString := bt.Expression(&ast.BasicLit{
			Kind:  token.STRING,
			Value: "\"Hello World\"",
		})
		res, err := etBasicString.Result()
		assert.NoError(t, err)
		assert.Equal(t, "string", res.Type)
		assert.Equal(t, "Hello World", res.ConstantValue)
	})

	t.Run("BasicLit_Int", func(t *testing.T) {
		etBasicInt := bt.Expression(&ast.BasicLit{
			Kind:  token.INT,
			Value: "42",
		})
		res, err := etBasicInt.Result()
		assert.NoError(t, err)
		assert.Equal(t, "int", res.Type)
		assert.Equal(t, "42", res.ConstantValue)
	})

	t.Run("CallExpr_ReservedFunction", func(t *testing.T) {
		etCallReserved := bt.Expression(&ast.CallExpr{
			Fun: &ast.Ident{Name: "len"},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"test\"",
				},
			},
		})
		res, err := etCallReserved.Result()
		assert.NoError(t, err)
		assert.Equal(t, "int", res.Type)
	})

	t.Run("CallExpr_OtherFunction", func(t *testing.T) {
		etCallOther := bt.Expression(&ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "strings"},
				Sel: &ast.Ident{Name: "Contains"},
			},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"test\"",
				},
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"es\"",
				},
			},
		})
		res, err := etCallOther.Result()
		assert.NoError(t, err)
		assert.Equal(t, "bool", res.Type)
	})
}
