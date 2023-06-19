package astUtils

import (
	"errors"
	"github.com/ls6-events/gengo/utils"
	"go/ast"
	"strconv"
)

func ExtractStatusCode(status ast.Node) (int, error) {
	var statusCode int
	var err error

	switch statusType := status.(type) {
	case *ast.BasicLit: // A constant status code is used (e.g. 200)
		statusCode, err = strconv.Atoi(statusType.Value)
		if err != nil {
			return 0, err
		}
	case *ast.Ident: // A constant defined in this package
		assignStmt, ok := statusType.Obj.Decl.(*ast.AssignStmt)
		if !ok {
			return 0, errors.New("status code is not a constant")
		}

		// Get the index of the status code constant
		var statementIndex int
		for i, expr := range assignStmt.Lhs {
			if expr.(*ast.Ident).Name == statusType.Name {
				statementIndex = i
				break
			}
		}

		switch rhs := assignStmt.Rhs[statementIndex].(type) {
		case *ast.BasicLit: // A constant status code is used (e.g. 200)
			// Get the value of the constant
			statusCode, err = strconv.Atoi(rhs.Value)
			if err != nil {
				return 0, err
			}
		case *ast.SelectorExpr: // A constant defined in another package
			// TODO Account for other constants in other packages (atm we just net/http (cheating I know))
			statusCode, err = utils.ConvertStatusCodeTypeToInt(rhs.Sel.Name)
			if err != nil {
				return 0, err
			}
		}
	case *ast.SelectorExpr: // A constant defined in another package
		// TODO Account for other constants in other packages (atm we just net/http (cheating I know))
		statusCode, err = utils.ConvertStatusCodeTypeToInt(statusType.Sel.Name)
		if err != nil {
			return 0, err
		}
	}

	// TODO DRY Cleanup

	return statusCode, nil
}
