package gin

import (
	"errors"
	"github.com/ls6-events/gengo"
	"github.com/ls6-events/gengo/utils"
	"github.com/rs/zerolog"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"strconv"
	"strings"
)

func extractContext(node *ast.FuncDecl) (string, error) {
	var ctxName string
	for _, param := range node.Type.Params.List {
		if len(param.Names) == 0 {
			continue
		}

		starExpr, ok := param.Type.(*ast.StarExpr)
		if !ok {
			continue
		}

		selectorExpr, ok := starExpr.X.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		ident, ok := selectorExpr.X.(*ast.Ident)
		if !ok || ident.Name != "gin" {
			continue
		}

		if selectorExpr.Sel.Name == "Context" {
			ctxName = param.Names[0].Name
			break
		}
	}

	if ctxName == "" {
		return "", errors.New("context parameter not found")
	}

	return ctxName, nil
}

func extractStatusCode(status ast.Node) (int, error) {
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

func parseFunction(s *gengo.Service, log zerolog.Logger, currRoute *gengo.Route, node *ast.FuncDecl, imports []*ast.ImportSpec, pkgName, pkgPath string, level int) error {
	// Get the variable name of the context parameter
	ctxName, err := extractContext(node)
	if err != nil {
		return err
	}

	// Loop over every statement in the function
	ast.Inspect(node.Body, func(n ast.Node) bool {
		// If a function is called
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Then it is either running a function from the gin.Context struct or a function utilising the context (any others are ignored)
		switch fun := callExpr.Fun.(type) {
		case *ast.SelectorExpr: // A method is called with a package name
			ident, ok := fun.X.(*ast.Ident)
			if !ok {
				return true
			}
			if ident.Name == ctxName { // If any gin.Context method is called
				switch fun.Sel.Name {
				// Return Types Below
				case "JSON":
					fallthrough
				case "XML":
					fallthrough
				case "YAML":
					fallthrough
				case "ProtoBuf":
					fallthrough
				case "Data":

					switch fun.Sel.Name {
					case "JSON":
						currRoute.ContentType = "application/json"
					case "XML":
						currRoute.ContentType = "application/xml"
					case "YAML":
						currRoute.ContentType = "application/yaml"
					case "ProtoBuf":
						currRoute.ContentType = "application/protobuf"
					}

					// Get the status code
					var statusCode int
					statusCode, err = extractStatusCode(callExpr.Args[0])
					if err != nil {
						return true
					}

					onExtract := func(result utils.ParseResult) {
						returnType := gengo.ReturnType{
							StatusCode: statusCode,
							Field: gengo.Field{
								Type:      result.VarName,
								SliceType: result.SliceType,
								MapKeyPkg: result.MapKeyPkg,
								MapKey:    result.MapKey,
								MapValue:  result.MapVal,
							},
						}
						if !gengo.IsAcceptedType(result.VarName) {
							s.AddToBeProcessed(result.PkgName, result.VarName)
							returnType.Field.Package = result.PkgName
						}

						currRoute.ReturnTypes = utils.AddReturnType(currRoute.ReturnTypes, returnType)
					}
					argNo := 1
					if fun.Sel.Name == "Data" {
						argNo = 2
					}

					err, ok = parseFromCalledFunction(log, callExpr, argNo, pkgName, pkgPath, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "String": // c.String
					currRoute.ContentType = "text/plain"

					var statusCode int
					statusCode, err = extractStatusCode(callExpr.Args[0])
					if err != nil {
						return false
					}

					returnType := gengo.ReturnType{
						StatusCode: statusCode,
						Field: gengo.Field{
							Type: "string",
						},
					}
					currRoute.ReturnTypes = utils.AddReturnType(currRoute.ReturnTypes, returnType)
					return true
				case "Status": // c.Status
					if len(currRoute.ReturnTypes) == 0 {
						var statusCode int
						statusCode, err = extractStatusCode(callExpr.Args[0])
						if err != nil {
							return false
						}

						returnType := gengo.ReturnType{
							StatusCode: statusCode,
							Field: gengo.Field{
								Type: "nil",
							},
						}
						currRoute.ReturnTypes = utils.AddReturnType(currRoute.ReturnTypes, returnType)
						return true
					}
				// Query Param methods
				case "GetQuery":
					fallthrough
				case "Query":
					onExtract := func(result utils.ParseResult) {
						currRoute.QueryParams = append(currRoute.QueryParams, gengo.Param{
							Name: strings.ReplaceAll(result.Value, "\"", ""),
							Type: result.VarName,
						})
					}

					err, ok = parseFromCalledFunction(log, callExpr, 0, pkgName, pkgPath, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "GetQueryArray":
					fallthrough
				case "QueryArray":
					onExtract := func(result utils.ParseResult) {
						currRoute.QueryParams = append(currRoute.QueryParams, gengo.Param{
							Name:    strings.ReplaceAll(result.Value, "\"", ""),
							Type:    result.VarName,
							IsArray: true,
						})
					}

					err, ok = parseFromCalledFunction(log, callExpr, 0, pkgName, pkgPath, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "GetQueryMap":
					fallthrough
				case "QueryMap":
					onExtract := func(result utils.ParseResult) {
						currRoute.QueryParams = append(currRoute.QueryParams, gengo.Param{
							Name:  strings.ReplaceAll(result.Value, "\"", ""),
							Type:  result.VarName,
							IsMap: true,
						})
					}

					err, ok = parseFromCalledFunction(log, callExpr, 0, pkgName, pkgPath, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "ShouldBindQuery":
					fallthrough
				case "BindQuery":
					onExtract := func(result utils.ParseResult) {
						queryParam := gengo.Param{
							IsBound: true,
							Type:    result.VarName,
						}
						if !gengo.IsAcceptedType(queryParam.Type) {
							s.AddToBeProcessed(result.PkgName, queryParam.Type)
							queryParam.Package = result.PkgName
						}

						currRoute.QueryParams = append(currRoute.QueryParams, queryParam)
					}

					err, ok = parseFromCalledFunction(log, callExpr, 0, pkgName, pkgPath, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				// Body Param methods
				case "ShouldBind":
					fallthrough
				case "Bind":
					onExtract := func(result utils.ParseResult) {
						bodyParam := gengo.Param{
							IsBound: true,
							Type:    result.VarName,
						}
						if !gengo.IsAcceptedType(bodyParam.Type) {
							s.AddToBeProcessed(result.PkgName, bodyParam.Type)
							bodyParam.Package = result.PkgName
						}

						currRoute.Body = append(currRoute.Body, bodyParam)
					}

					currRoute.BodyType = "form"

					err, ok = parseFromCalledFunction(log, callExpr, 0, pkgName, pkgPath, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "ShouldBindJSON":
					fallthrough
				case "BindJSON":
					onExtract := func(result utils.ParseResult) {
						bodyParam := gengo.Param{
							IsBound: true,
							Type:    result.VarName,
						}
						if !gengo.IsAcceptedType(bodyParam.Type) {
							s.AddToBeProcessed(result.PkgName, bodyParam.Type)
							bodyParam.Package = result.PkgName
						}

						currRoute.Body = append(currRoute.Body, bodyParam)
					}

					currRoute.BodyType = "application/json"

					err, ok = parseFromCalledFunction(log, callExpr, 0, pkgName, pkgPath, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "ShouldBindXML":
					fallthrough
				case "BindXML":
					onExtract := func(result utils.ParseResult) {
						bodyParam := gengo.Param{
							IsBound: true,
							Type:    result.VarName,
						}
						if !gengo.IsAcceptedType(bodyParam.Type) {
							s.AddToBeProcessed(result.PkgName, bodyParam.Type)
							bodyParam.Package = result.PkgName
						}

						currRoute.Body = append(currRoute.Body, bodyParam)
					}

					currRoute.BodyType = "application/xml"

					err, ok = parseFromCalledFunction(log, callExpr, 0, pkgName, pkgPath, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "ShouldBindYAML":
					fallthrough
				case "BindYAML":
					onExtract := func(result utils.ParseResult) {
						bodyParam := gengo.Param{
							IsBound: true,
							Type:    result.VarName,
						}
						if !gengo.IsAcceptedType(bodyParam.Type) {
							s.AddToBeProcessed(result.PkgName, bodyParam.Type)
							bodyParam.Package = result.PkgName
						}

						currRoute.Body = append(currRoute.Body, bodyParam)
					}

					currRoute.BodyType = "application/yaml"

					err, ok = parseFromCalledFunction(log, callExpr, 0, pkgName, pkgPath, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "GetPostForm":
					fallthrough
				case "PostForm":
					onExtract := func(result utils.ParseResult) {
						currRoute.Body = append(currRoute.Body, gengo.Param{
							Name: strings.ReplaceAll(result.Value, "\"", ""),
							Type: result.VarName,
						})
					}

					currRoute.BodyType = "application/x-www-form-urlencoded"

					err, ok = parseFromCalledFunction(log, callExpr, 0, pkgName, pkgPath, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "GetPostFormArray":
					fallthrough
				case "PostFormArray":
					onExtract := func(result utils.ParseResult) {
						currRoute.Body = append(currRoute.Body, gengo.Param{
							Name:    strings.ReplaceAll(result.Value, "\"", ""),
							Type:    result.VarName,
							IsArray: true,
						})
					}

					currRoute.BodyType = "application/x-www-form-urlencoded"

					err, ok = parseFromCalledFunction(log, callExpr, 0, pkgName, pkgPath, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "GetPostFormMap":
					fallthrough
				case "PostFormMap":
					onExtract := func(result utils.ParseResult) {
						currRoute.Body = append(currRoute.Body, gengo.Param{
							Name:  strings.ReplaceAll(result.Value, "\"", ""),
							Type:  result.VarName,
							IsMap: true,
						})
					}

					currRoute.BodyType = "application/x-www-form-urlencoded"

					err, ok = parseFromCalledFunction(log, callExpr, 0, pkgName, pkgPath, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				default:
					return true
				}
			} else { // Check if parameters contain the context
				var hasContext bool
				for _, arg := range callExpr.Args {
					switch argType := arg.(type) {
					case *ast.Ident:
						if argType.Name == ctxName {
							hasContext = true
						}
					}
				}

				if hasContext {
					if pkgPath == "" {
						pkgPath = pkgName
					}
					nPkgPath := utils.ParseInputPath(imports, ident.Name, pkgPath)
					var nPkg *packages.Package
					nPkg, err = loadPackage(nPkgPath)
					if err != nil {
						return false
					}

					var funcDecl *ast.FuncDecl
					var nImports []*ast.ImportSpec
					for _, file := range nPkg.Syntax {
						for _, decl := range file.Decls {
							f, ok := decl.(*ast.FuncDecl)
							if !ok {
								continue
							}

							if f.Name.Name == fun.Sel.Name {
								nImports = file.Imports
								funcDecl = f
								break
							}
						}
					}

					if funcDecl == nil {
						return true
					}

					err = parseFunction(s, log, currRoute, funcDecl, nImports, nPkg.Name, nPkgPath, level+1)
					if err != nil {
						return false
					}
				}
			}
		case *ast.Ident: // A function is called without a package name
			var hasContext bool // Check if parameters contain the context
			for _, arg := range callExpr.Args {
				switch argType := arg.(type) {
				case *ast.Ident:
					if argType.Name == ctxName {
						hasContext = true
					}
				}
			}

			if hasContext {
				if pkgPath == "" {
					pkgPath = pkgName
				}
				nPkgPath := utils.ParseInputPath(imports, pkgName, pkgPath)

				if nPkgPath == "main" {
					nPkgPath, err = s.GetMainPackageName()
					if err != nil {
						return false
					}
				}

				var nPkg *packages.Package
				nPkg, err = loadPackage(nPkgPath)
				if err != nil {
					return false
				}

				var funcDecl *ast.FuncDecl
				var nImports []*ast.ImportSpec
				for _, file := range nPkg.Syntax {
					for _, decl := range file.Decls {
						f, ok := decl.(*ast.FuncDecl)
						if !ok {
							continue
						}

						if f.Name.Name == fun.Name {
							nImports = file.Imports
							funcDecl = f
							break
						}
					}
				}

				if funcDecl == nil {
					err = errors.New("function not found")
					return true
				}

				nSplitPkg := strings.Split(nPkgPath, "/")
				err = parseFunction(s, log, currRoute, funcDecl, nImports, nSplitPkg[len(nSplitPkg)-1], strings.Join(nSplitPkg[:len(nSplitPkg)-1], "/"), level+1)
				if err != nil {
					return false
				}
			}
		}

		return true
	})

	if err != nil {
		return err
	}

	if len(currRoute.ReturnTypes) == 0 && level == 0 {
		return errors.New("return type not found")
	}

	return nil
}

func parseFromCalledFunction(log zerolog.Logger, callExpr *ast.CallExpr, argNo int, pkgName, pkgPath string, imports []*ast.ImportSpec, onExtract func(result utils.ParseResult)) (error, bool) {
	var err error
	arg := callExpr.Args[argNo]
	switch argType := arg.(type) {
	case *ast.CompositeLit: // A constant defined in the arguments
		switch compositLit := argType.Type.(type) {
		case *ast.Ident: // A constant defined in this package
			onExtract(utils.ParseResult{
				PkgName: pkgName,
				VarName: compositLit.Name,
			})

			return nil, true
		case *ast.SelectorExpr: // A constant defined in another package
			ident, ok := compositLit.X.(*ast.Ident)
			if !ok {
				return nil, false
			}

			onExtract(utils.ParseResult{
				PkgName: utils.ParseInputPath(imports, ident.Name, pkgPath),
				VarName: compositLit.Sel.Name,
			})

			return nil, true
		}
	case *ast.Ident: // A variable used in the arguments
		assignStmt, ok := argType.Obj.Decl.(*ast.AssignStmt)
		if !ok {
			return nil, false
		}

		var assignedIndex int
		for i, expr := range assignStmt.Lhs {
			if expr.(*ast.Ident).Name == argType.Name {
				assignedIndex = i
				break
			}
		}

		var assignedExpr ast.Expr
		if len(assignStmt.Lhs) == len(assignStmt.Rhs) { // If the number of variables and values are the same
			assignedExpr = assignStmt.Rhs[assignedIndex]
		} else { // If the number of variables and values are different (i.e. a function call)
			assignedExpr = assignStmt.Rhs[0]
		}

		onExternalPkg := func(funcName, pkgName, pkgPath string) error {
			// We need all this logic here because we need to check the return type of the function against that package's imports
			nPkgPath := utils.ParseInputPath(imports, pkgName, pkgPath)
			var pkg *packages.Package
			pkg, err = loadPackage(nPkgPath)
			if err != nil {
				return err
			}

			var pkgImports []*ast.ImportSpec
			var funcDecl *ast.FuncDecl
			for _, file := range pkg.Syntax {
				for _, decl := range file.Decls {
					if f, ok := decl.(*ast.FuncDecl); ok {
						if f.Name.Name == funcName {
							pkgImports = file.Imports
							funcDecl = f
							break
						}
					}
				}
			}

			var funcReturnIndex int
			for i, field := range assignStmt.Lhs {
				if f, ok := field.(*ast.Ident); ok {
					if f.Name == argType.Name {
						funcReturnIndex = i
					}
				}
			}

			field := funcDecl.Type.Results.List[funcReturnIndex]

			res, ok := parseFunctionReturnTypes(log, field.Type, argType)
			if !ok {
				return nil
			}

			onExtract(utils.ParseResult{
				PkgName:   utils.ParseInputPath(pkgImports, res.PkgName, nPkgPath),
				VarName:   res.VarName,
				Value:     res.Value,
				MapKeyPkg: res.MapKeyPkg,
				MapKey:    res.MapKey,
				MapVal:    res.MapVal,
				SliceType: res.SliceType,
			})

			return nil
		}

		var res utils.ParseResult
		res, err, ok = parseAssignStatement(log, assignedExpr, assignStmt, pkgPath, pkgName, imports, argType, onExternalPkg)
		if !ok {
			return err, false
		}

		onExtract(utils.ParseResult{
			VarName:   res.VarName,
			PkgName:   utils.ParseInputPath(imports, res.PkgName, pkgPath),
			Value:     res.Value,
			MapKeyPkg: res.MapKeyPkg,
			MapKey:    res.MapKey,
			MapVal:    res.MapVal,
			SliceType: res.SliceType,
		})

		return nil, true
	case *ast.BasicLit: // A literal used in the arguments
		onExtract(utils.ParseResult{
			PkgName: pkgName,
			VarName: strings.ToLower(argType.Kind.String()),
			Value:   argType.Value,
		})

		return nil, true
	default:
		return nil, false
	}

	return nil, false
}

func parseAssignStatement(log zerolog.Logger, expr ast.Expr, assignStmt *ast.AssignStmt, pkgPath string, pkgName string, imports []*ast.ImportSpec, argType *ast.Ident, onExternalPkg func(funcName, pkgName, pkgPath string) error) (utils.ParseResult, error, bool) {
	var err error
	var res utils.ParseResult
	switch rhs := expr.(type) {
	case *ast.UnaryExpr:
		return parseAssignStatement(log, rhs.X, assignStmt, pkgPath, pkgName, imports, argType, onExternalPkg)
	case *ast.CompositeLit:
		switch compositLit := rhs.Type.(type) {
		case *ast.Ident:
			res = utils.SplitIdentSelectorExpr(compositLit, pkgName)
		case *ast.SelectorExpr:
			res = utils.SplitIdentSelectorExpr(compositLit, pkgName)
		case *ast.ArrayType:
			embeddedType := utils.SplitIdentSelectorExpr(compositLit.Elt, pkgName)
			res = utils.ParseResult{
				VarName:   "slice",
				PkgName:   embeddedType.PkgName,
				SliceType: embeddedType.VarName,
			}
		case *ast.MapType:
			keyType := utils.SplitIdentSelectorExpr(compositLit.Key, pkgName)
			valueType := utils.SplitIdentSelectorExpr(compositLit.Value, pkgName)
			res = utils.ParseResult{
				VarName:   "map",
				MapKey:    keyType.VarName,
				MapKeyPkg: keyType.PkgName,
				MapVal:    valueType.VarName,
				PkgName:   valueType.PkgName,
			}
		}

	case *ast.BasicLit:
		res = utils.ParseResult{
			VarName: strings.ToLower(rhs.Kind.String()),
			PkgName: pkgName,
			Value:   rhs.Value,
		}
	case *ast.Ident:
		assignStmt, ok := rhs.Obj.Decl.(*ast.AssignStmt)
		if !ok {
			return utils.ParseResult{}, nil, false
		}

		var assignedIndex int
		for i, expr := range assignStmt.Lhs {
			if expr.(*ast.Ident).Name == rhs.Name {
				assignedIndex = i
				break
			}
		}

		var assignedExpr ast.Expr
		if len(assignStmt.Lhs) == len(assignStmt.Rhs) { // If the number of variables and values are the same
			assignedExpr = assignStmt.Rhs[assignedIndex]
		} else { // If the number of variables and values are different (i.e. a function call)
			assignedExpr = assignStmt.Rhs[0]
		}

		return parseAssignStatement(log, assignedExpr, assignStmt, pkgPath, pkgName, imports, argType, onExternalPkg)
	case *ast.CallExpr:
		switch fun := rhs.Fun.(type) {
		case *ast.SelectorExpr: // foo.Bar()
			ident, ok := fun.X.(*ast.Ident)
			if !ok {
				return utils.ParseResult{}, nil, false
			}

			err = onExternalPkg(fun.Sel.Name, ident.Name, pkgPath)
			if err != nil {
				return utils.ParseResult{}, err, false
			} else {
				return utils.ParseResult{}, nil, true
			}
		case *ast.Ident: // Bar()
			funcDecl, ok := fun.Obj.Decl.(*ast.FuncDecl)
			if !ok {
				return utils.ParseResult{}, nil, false
			}

			var funcReturnIndex int
			for i, field := range assignStmt.Lhs {
				if f, ok := field.(*ast.Ident); ok {
					if f.Name == argType.Name {
						funcReturnIndex = i
					}
				}
			}

			field := funcDecl.Type.Results.List[funcReturnIndex]

			res, ok = parseFunctionReturnTypes(log, field.Type, argType)
			if !ok {
				return utils.ParseResult{}, nil, false
			}
		default:
			return utils.ParseResult{}, nil, false
		}
	default:
		return utils.ParseResult{}, nil, false
	}
	return res, nil, true
}

func parseFunctionReturnTypes(log zerolog.Logger, node ast.Node, argType *ast.Ident) (utils.ParseResult, bool) {
	switch fieldType := node.(type) {
	case *ast.StarExpr:
		return parseFunctionReturnTypes(log, fieldType.X, argType)
	case *ast.SelectorExpr:
		return utils.SplitIdentSelectorExpr(fieldType, argType.Name), true
	case *ast.Ident:
		return utils.SplitIdentSelectorExpr(fieldType, argType.Name), true
	default:
		return utils.ParseResult{}, false
	}
}
