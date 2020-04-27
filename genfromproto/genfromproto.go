package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

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
		for _, field := range st.Fields.List {
			var found bool
			for _, name := range field.Names {
				if name.Name != "fieldMask" {
					found = true
					break
				}
			}
			if !found {
				continue
			}
			typ, ok := field.Type.(*ast.Ident)
			if !ok {
				errFound = true
				fmt.Printf("expected ident type for field, got %t\n", field.Type)
				return
			}
			typ.Name = "*" + typ.Name
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
