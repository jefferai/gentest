package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"regexp"
)

var regex = regexp.MustCompile(`(json:".*")`)

type visitFn func(node ast.Node)

func (fn visitFn) Visit(node ast.Node) ast.Visitor {
	fn(node)
	return fn
}

func main() {
	in := os.Args[1]
	inName := os.Args[2]
	out := os.Args[3]
	outName := os.Args[4]

	_ = out

	fset := token.NewFileSet()
	inAst, err := parser.ParseFile(fset, in, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("error parsing %s: %v\n", in, err)
		os.Exit(1)
	}

	var errFound bool

	ast.Walk(visitFn(func(n ast.Node) {
		spec, ok := n.(*ast.TypeSpec)
		if !ok {
			return
		}
		if spec.Name == nil {
			return
		}
		if spec.Name.Name != inName {
			return
		}
		spec.Name.Name = outName
		st, ok := spec.Type.(*ast.StructType)
		if !ok {
			errFound = true
			fmt.Printf("expected struct type for identifier, got %t\n", spec.Type)
			return
		}

		if st.Fields.List == nil {
			errFound = true
			fmt.Printf("no fields found in %q\n", inName)
			return
		}
		var elideList []int
		defer func() {
			var cutCount int
			for _, val := range elideList {
				loc := val - cutCount
				st.Fields.List = append(st.Fields.List[:loc], st.Fields.List[loc+1:]...)
				cutCount++
			}
		}()

		for i, field := range st.Fields.List {
			if !field.Names[0].IsExported() {
				elideList = append(elideList, i)
				continue
			}

			// Anything that isn't a basic type we expect to be a star
			// expression with a selector; that is, a wrapper value, timestamp
			// value, etc.
			//
			// TODO: this isn't necessarily a good assumption, which means we
			// might get failures with other types. This is an internal tools
			// only; we can revisit as needed!
			var selectorExpr *ast.SelectorExpr
			switch typ := field.Type.(type) {
			case *ast.Ident:
				typ.Name = "*" + typ.Name
				goto TAGMODIFY
			case *ast.StarExpr:
				switch nextTyp := typ.X.(type) {
				case *ast.Ident:
					// Already a pointer, don't do anything
					goto TAGMODIFY
				case *ast.SelectorExpr:
					selectorExpr = nextTyp
				}
			case *ast.SelectorExpr:
				selectorExpr = typ
			}

			switch {
			case selectorExpr != nil:
				xident, ok := selectorExpr.X.(*ast.Ident)
				if !ok {
					fmt.Printf("unexpected non-ident type in selector\n")
					os.Exit(1)
				}

				switch xident.Name {
				case "wrappers":
					switch selectorExpr.Sel.Name {
					case "StringValue":
						st.Fields.List[i] = &ast.Field{
							Names: field.Names,
							Type: &ast.Ident{
								Name: "*string",
							},
							Tag: field.Tag,
						}
					case "BoolValue":
						st.Fields.List[i] = &ast.Field{
							Names: field.Names,
							Type: &ast.Ident{
								Name: "*bool",
							},
							Tag: field.Tag,
						}
					default:
						fmt.Printf("unhandled wrappers selector sel name %q\n", selectorExpr.Sel.Name)
						os.Exit(1)
					}

				case "timestamp":
					switch selectorExpr.Sel.Name {
					case "Timestamp":
						st.Fields.List[i] = &ast.Field{
							Names: field.Names,
							Type: &ast.Ident{
								Name: "time.Time",
							},
							Tag: field.Tag,
						}

					default:
						fmt.Printf("unhandled timestamp selector sel name %q\n", selectorExpr.Sel.Name)
						os.Exit(1)
					}

				default:
					fmt.Printf("unhandled xident name %q\n", xident.Name)
					os.Exit(1)
				}

			default:
				fmt.Println("unhandled non-ident, non-selector case")
				os.Exit(1)
			}

		TAGMODIFY:
			st.Fields.List[i].Tag.Value = "`" + regex.FindString(st.Fields.List[i].Tag.Value) + "`"
		}
	}), inAst)

	if errFound {
		os.Exit(1)
	}

	buf := new(bytes.Buffer)
	if err := format.Node(buf, fset, inAst); err != nil {
		fmt.Printf("error formatting new code: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(buf.String())
}
