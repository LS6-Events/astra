package gin

import (
	"errors"
	"github.com/ls6-events/gengo"
	"github.com/ls6-events/gengo/utils"
	"github.com/ls6-events/gengo/utils/astUtils"
	"github.com/rs/zerolog"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"strings"
)

// parseFunction parses a function and adds it to the service
// It is designed to be called recursively should it be required
// The level parameter is used to determine the depth of recursion
// And the package name and path are used to determine the package of the currently analysed function
// The currRoute reference is used to manipulate the current route being analysed
// The imports are used to determine the package of the context variable
func parseFunction(s *gengo.Service, log zerolog.Logger, currRoute *gengo.Route, node *ast.FuncLit, imports []*ast.ImportSpec, pkgName, pkgPath string, level int) error {
	// Get the variable name of the context parameter
	ctxName, err := astUtils.ExtractContext("github.com/gin-gonic/gin", "*Context", node, imports)
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
					statusCode, err = astUtils.ExtractStatusCode(callExpr.Args[0])
					if err != nil {
						return true
					}

					onExtract := func(result astUtils.ParseResult) {
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
						if result.SliceType != "" && !gengo.IsAcceptedType(result.SliceType) {
							s.AddToBeProcessed(result.PkgName, result.SliceType)
							returnType.Field.Package = result.PkgName
						}
						if result.MapKey != "" && !gengo.IsAcceptedType(result.MapKey) {
							s.AddToBeProcessed(result.MapKeyPkg, result.MapKey)
							returnType.Field.Package = result.MapKeyPkg
						}
						if result.MapVal != "" && !gengo.IsAcceptedType(result.MapVal) {
							s.AddToBeProcessed(result.MapKeyPkg, result.MapVal)
							returnType.Field.Package = result.MapKeyPkg
						}

						currRoute.ReturnTypes = utils.AddReturnType(currRoute.ReturnTypes, returnType)
					}
					argNo := 1
					if fun.Sel.Name == "Data" {
						argNo = 2
					}

					err, ok = parseFromCalledFunction(s, log, callExpr, argNo, pkgName, pkgPath, s.WorkDir, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "String": // c.String
					currRoute.ContentType = "text/plain"

					var statusCode int
					statusCode, err = astUtils.ExtractStatusCode(callExpr.Args[0])
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
					var statusCode int
					statusCode, err = astUtils.ExtractStatusCode(callExpr.Args[0])
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
				// Query Param methods
				case "GetQuery":
					fallthrough
				case "Query":
					onExtract := func(result astUtils.ParseResult) {
						currRoute.QueryParams = append(currRoute.QueryParams, gengo.Param{
							Name: strings.ReplaceAll(result.Value, "\"", ""),
							Type: result.VarName,
						})
					}

					err, ok = parseFromCalledFunction(s, log, callExpr, 0, pkgName, pkgPath, s.WorkDir, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "GetQueryArray":
					fallthrough
				case "QueryArray":
					onExtract := func(result astUtils.ParseResult) {
						currRoute.QueryParams = append(currRoute.QueryParams, gengo.Param{
							Name:    strings.ReplaceAll(result.Value, "\"", ""),
							Type:    result.VarName,
							IsArray: true,
						})
					}

					err, ok = parseFromCalledFunction(s, log, callExpr, 0, pkgName, pkgPath, s.WorkDir, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "GetQueryMap":
					fallthrough
				case "QueryMap":
					onExtract := func(result astUtils.ParseResult) {
						currRoute.QueryParams = append(currRoute.QueryParams, gengo.Param{
							Name:  strings.ReplaceAll(result.Value, "\"", ""),
							Type:  result.VarName,
							IsMap: true,
						})
					}

					err, ok = parseFromCalledFunction(s, log, callExpr, 0, pkgName, pkgPath, s.WorkDir, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "ShouldBindQuery":
					fallthrough
				case "BindQuery":
					onExtract := func(result astUtils.ParseResult) {
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

					err, ok = parseFromCalledFunction(s, log, callExpr, 0, pkgName, pkgPath, s.WorkDir, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				// Body Param methods
				case "ShouldBind":
					fallthrough
				case "Bind":
					onExtract := func(result astUtils.ParseResult) {
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

					err, ok = parseFromCalledFunction(s, log, callExpr, 0, pkgName, pkgPath, s.WorkDir, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "ShouldBindJSON":
					fallthrough
				case "BindJSON":
					onExtract := func(result astUtils.ParseResult) {
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

					err, ok = parseFromCalledFunction(s, log, callExpr, 0, pkgName, pkgPath, s.WorkDir, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "ShouldBindXML":
					fallthrough
				case "BindXML":
					onExtract := func(result astUtils.ParseResult) {
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

					err, ok = parseFromCalledFunction(s, log, callExpr, 0, pkgName, pkgPath, s.WorkDir, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "ShouldBindYAML":
					fallthrough
				case "BindYAML":
					onExtract := func(result astUtils.ParseResult) {
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

					err, ok = parseFromCalledFunction(s, log, callExpr, 0, pkgName, pkgPath, s.WorkDir, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "GetPostForm":
					fallthrough
				case "PostForm":
					onExtract := func(result astUtils.ParseResult) {
						currRoute.Body = append(currRoute.Body, gengo.Param{
							Name: strings.ReplaceAll(result.Value, "\"", ""),
							Type: result.VarName,
						})
					}

					currRoute.BodyType = "application/x-www-form-urlencoded"

					err, ok = parseFromCalledFunction(s, log, callExpr, 0, pkgName, pkgPath, s.WorkDir, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "GetPostFormArray":
					fallthrough
				case "PostFormArray":
					onExtract := func(result astUtils.ParseResult) {
						currRoute.Body = append(currRoute.Body, gengo.Param{
							Name:    strings.ReplaceAll(result.Value, "\"", ""),
							Type:    result.VarName,
							IsArray: true,
						})
					}

					currRoute.BodyType = "application/x-www-form-urlencoded"

					err, ok = parseFromCalledFunction(s, log, callExpr, 0, pkgName, pkgPath, s.WorkDir, imports, onExtract)
					if err != nil {
						return false
					}
					return !ok
				case "GetPostFormMap":
					fallthrough
				case "PostFormMap":
					onExtract := func(result astUtils.ParseResult) {
						currRoute.Body = append(currRoute.Body, gengo.Param{
							Name:  strings.ReplaceAll(result.Value, "\"", ""),
							Type:  result.VarName,
							IsMap: true,
						})
					}

					currRoute.BodyType = "application/x-www-form-urlencoded"

					err, ok = parseFromCalledFunction(s, log, callExpr, 0, pkgName, pkgPath, s.WorkDir, imports, onExtract)
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
					nPkgPath := astUtils.ParseInputPath(imports, ident.Name, pkgPath)
					var nPkg *packages.Package
					nPkg, err = astUtils.LoadPackage(nPkgPath, s.WorkDir)
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

					err = parseFunction(s, log, currRoute, astUtils.FuncDeclToFuncLit(funcDecl), nImports, nPkg.Name, nPkgPath, level+1)
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
				nPkgPath := astUtils.ParseInputPath(imports, pkgName, pkgPath)

				if nPkgPath == "main" {
					nPkgPath, err = s.GetMainPackageName()
					if err != nil {
						return false
					}
				}

				var nPkg *packages.Package
				nPkg, err = astUtils.LoadPackage(nPkgPath, s.WorkDir)
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
				err = parseFunction(s, log, currRoute, astUtils.FuncDeclToFuncLit(funcDecl), nImports, nSplitPkg[len(nSplitPkg)-1], strings.Join(nSplitPkg[:len(nSplitPkg)-1], "/"), level+1)
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

// parseFromCalledFunction parses a function call that an argument is a reference to something
// It returns true if a new pass is needed to parse the function (i.e. it has had to extract another different type)
func parseFromCalledFunction(s *gengo.Service, log zerolog.Logger, callExpr *ast.CallExpr, argNo int, pkgName, pkgPath, workDir string, imports []*ast.ImportSpec, onExtract func(result astUtils.ParseResult)) (error, bool) {
	arg := callExpr.Args[argNo]
	switch argType := arg.(type) {
	case *ast.UnaryExpr: // A reference to a constant defined in the arguments
		switch unaryExpr := argType.X.(type) {
		case *ast.Ident: // A constant defined in this package
			return parseIdentAndTrace(s, log, unaryExpr, pkgName, pkgPath, workDir, imports, onExtract)
		case *ast.SelectorExpr: // A constant defined in another package
			ident, ok := unaryExpr.X.(*ast.Ident)
			if !ok {
				return nil, false
			}

			onExtract(astUtils.ParseResult{
				PkgName: astUtils.ParseInputPath(imports, ident.Name, pkgPath),
				VarName: unaryExpr.Sel.Name,
			})

			return nil, true
		}
	case *ast.CompositeLit: // A constant defined in the arguments
		switch compositLit := argType.Type.(type) {
		case *ast.Ident: // A constant defined in this package
			onExtract(astUtils.ParseResult{
				PkgName: pkgName,
				VarName: compositLit.Name,
			})

			return nil, true
		case *ast.SelectorExpr: // A constant defined in another package
			ident, ok := compositLit.X.(*ast.Ident)
			if !ok {
				return nil, false
			}

			onExtract(astUtils.ParseResult{
				PkgName: astUtils.ParseInputPath(imports, ident.Name, pkgPath),
				VarName: compositLit.Sel.Name,
			})

			return nil, true
		}
	case *ast.Ident: // A variable used in the arguments
		return parseIdentAndTrace(s, log, argType, pkgPath, pkgName, workDir, imports, onExtract)
	case *ast.CallExpr: // A function call used in the arguments
		res, err := astUtils.HandleReservedFunctions(argType, pkgName)
		if err != nil {
			return err, false
		}
		if res.VarName != "" {
			onExtract(res)
			return nil, true
		}

		res, node, nImports, err := astUtils.FindDeclInPackage(argType.Fun, imports, pkgName, pkgPath, workDir, func(res astUtils.ParseResult) (astUtils.ParseResult, error) {
			if res.PkgName == "main" {
				var err error
				res.PkgName, err = s.GetMainPackageName()
				if err != nil {
					return astUtils.ParseResult{}, err
				}
			}

			return res, nil
		})
		if err != nil {
			return err, false
		}

		if funcDecl, ok := node.(*ast.FuncDecl); ok {
			returnType := funcDecl.Type.Results.List[0] // It is assumed that the function only returns one value, as it is inline in the arguments

			res, ok = astUtils.ParseFunctionReturnTypes(log, returnType.Type, res.PkgName)
			res.PkgName = astUtils.ParseInputPath(nImports, res.PkgName, pkgPath)
			if !ok {
				return nil, false
			} else {
				onExtract(res)
				return nil, true
			}

		} else {
			return nil, false
		}

	case *ast.BasicLit: // A literal used in the arguments
		onExtract(astUtils.ParseResult{
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

// parseIdentAndTrace parses an identifier and traces it back to its definition
// It returns true if a new pass is needed to parse the function (i.e. it has had to extract another different type)
// It is designed to match any number of arguments on either side
func parseIdentAndTrace(s *gengo.Service, log zerolog.Logger, argType *ast.Ident, pkgPath, pkgName, workDir string, imports []*ast.ImportSpec, onExtract func(result astUtils.ParseResult)) (error, bool) {
	var assignedExpr ast.Expr
	var assignStmt *ast.AssignStmt
	var ok bool
	if argType.Obj == nil {
		res, node, nImports, err := astUtils.FindDeclInPackage(argType, imports, pkgName, pkgPath, workDir, func(res astUtils.ParseResult) (astUtils.ParseResult, error) {
			if res.PkgName == "main" {
				var err error
				res.PkgName, err = s.GetMainPackageName()
				if err != nil {
					return astUtils.ParseResult{}, err
				}
			}

			return res, nil
		})
		if err != nil {
			return err, false
		}

		if valueSpec, ok := node.(*ast.ValueSpec); ok { // It is a variable that is assigned a value, so we need to trace it back to its definition
			var assignedIndex int
			for i, expr := range valueSpec.Names {
				if expr.Name == argType.Name {
					assignedIndex = i
					break
				}
			}

			assignedExpr = valueSpec.Values[assignedIndex]
			pkgPath = astUtils.ParseInputPath(nImports, res.PkgName, pkgPath)
			pkgName = res.PkgName
			imports = nImports
		}
	} else {
		assignStmt, ok = argType.Obj.Decl.(*ast.AssignStmt)
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

		if len(assignStmt.Lhs) == len(assignStmt.Rhs) { // If the number of variables and values are the same
			assignedExpr = assignStmt.Rhs[assignedIndex]
		} else { // If the number of variables and values are different (i.e. a function call)
			assignedExpr = assignStmt.Rhs[0]
		}
	}

	onExternalPkg := func(funcName, pkgName, pkgPath string) error {
		// We need all this logic here because we need to check the return type of the function against that package's imports

		if pkgName == "main" {
			var err error
			pkgName, err = s.GetMainPackageName()
			if err != nil {
				return err
			}
		}

		nPkgPath := astUtils.ParseInputPath(imports, pkgName, pkgPath)
		pkg, err := astUtils.LoadPackage(nPkgPath, workDir)
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

		res, ok := astUtils.ParseFunctionReturnTypes(log, field.Type, argType.Name)
		if !ok {
			return nil
		}

		onExtract(astUtils.ParseResult{
			PkgName:   astUtils.ParseInputPath(pkgImports, res.PkgName, nPkgPath),
			VarName:   res.VarName,
			Value:     res.Value,
			MapKeyPkg: res.MapKeyPkg,
			MapKey:    res.MapKey,
			MapVal:    res.MapVal,
			SliceType: res.SliceType,
		})

		return nil
	}

	var res astUtils.ParseResult
	res, err, isExtractRequired := astUtils.ParseAssignStatement(log, assignedExpr, assignStmt, pkgPath, pkgName, imports, argType, onExternalPkg)
	if err != nil {
		return err, false
	} else if !isExtractRequired {
		return nil, true
	}

	onExtract(astUtils.ParseResult{
		VarName:   res.VarName,
		PkgName:   astUtils.ParseInputPath(imports, res.PkgName, pkgPath),
		Value:     res.Value,
		MapKeyPkg: res.MapKeyPkg,
		MapKey:    res.MapKey,
		MapVal:    res.MapVal,
		SliceType: res.SliceType,
	})
	return nil, true
}
