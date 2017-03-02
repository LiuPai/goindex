package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type item struct {
	Title   string
	Element string
	Pos     token.Pos
}

func index(file []byte) ([]item, error) {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", file, 0)
	if err != nil {
		return nil, err
	}
	var result []item
	for _, decl := range f.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			for iSpec, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.ValueSpec:
					for iIdent, ident := range spec.Names {
						item := item{
							Element: ident.Obj.Name,
							Pos:     ident.Pos(),
						}
						if iIdent == 0 && iSpec == 0 {
							item.Pos = decl.Pos()
						}
						kind := ident.Obj.Kind.String()
						switch kind {
						case "var":
							kind = "Variables"
						case "const":
							kind = "Constants"
						}
						item.Title = kind
						result = append(result, item)
					}
				case *ast.TypeSpec:
					item := item{
						Element: spec.Name.Name,
						Pos:     spec.Pos(),
					}
					if iSpec == 0 {
						item.Pos = decl.Pos()
					}
					switch spec.Type.(type) {
					case *ast.StructType:
						item.Title = "Struct"
					case *ast.InterfaceType:
						item.Title = "Interface"
					default:
						item.Title = "Alias"
						item.Element = string(file[spec.Name.Pos()-1 : spec.Type.End()-1])
					}
					result = append(result, item)
				}
			}
		case *ast.FuncDecl:
			element := strings.Replace(string(file[decl.Name.Pos()-1:decl.Type.End()-1]), "\n", " ", -1)
			item := item{
				Element: element,
				Pos:     decl.Pos(),
			}

			if decl.Recv == nil {
				item.Title = "Function"
			} else {
				if len(decl.Recv.List) != 1 {
					continue
				}
				switch t := decl.Recv.List[0].Type.(type) {
				case *ast.StarExpr:
					item.Title = fmt.Sprintf("Method %s", t.X.(*ast.Ident).Name)
					item.Element = fmt.Sprintf("(*%s) %s", t.X.(*ast.Ident).Name, item.Element)
				case *ast.Ident:
					item.Title = fmt.Sprintf("Method %s", t.Name)
					item.Element = fmt.Sprintf("(%s) %s", t.Name, item.Element)
				}

			}
			result = append(result, item)
		}
	}
	return result, nil
}
