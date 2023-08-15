package smrcptr

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "smrcptr",
	Doc:      "detect mixing pointer and value method receivers for the same type",
	Run:      run,
	Flags:    flag.FlagSet{},
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var (
	enableConstructorCheck bool
	skipSTD                bool
	skipGenerated          bool
)

func init() {
	Analyzer.Flags.BoolVar(&enableConstructorCheck, "constructor", false, `enable constructor return type check`)
	Analyzer.Flags.BoolVar(&skipSTD, "skip-std", true, `skip methods that satisfy typical interfaces from standard pacakges`)
	Analyzer.Flags.BoolVar(&skipGenerated, "skip-generated", true, `skip generated files`)
}

type functionSelector interface {
	SelectFunction(fn *ast.FuncDecl) bool
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	var fnSelector functionSelector = mapNameFunctionSelector{def: true}
	if skipSTD {
		// TODO: use type of function too or find native way to test if interface is satisfied
		fnSelector = mapNameFunctionSelector{
			def: true,
			fns: map[string]bool{
				"UnmarshalJSON":    false, // encoding
				"UnmarshalText":    false, // encoding
				"UnmarshalBinary":  false, // encoding
				"UnmarshalXML":     false, // encoding
				"UnmarshalXMLAttr": false, // encoding
				"Scanner":          false, // database/sql
				"Scan":             false, // fmt
				"Read":             false, // io
			},
		}
	}

	// non-nil pointers for memory efficiency
	typePtrFns := map[string][]*ast.FuncDecl{}
	typeValFns := map[string][]*ast.FuncDecl{}

	fset := token.NewFileSet()

	inspect.Preorder([]ast.Node{&ast.FuncDecl{}}, func(n ast.Node) {
		if skipGenerated {
			// TODO: find way to reuse ast.File from analysis inspector
			fname := pass.Fset.Position(n.Pos()).Filename
			f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments|parser.PackageClauseOnly)
			if err != nil {
				log.Fatalf("failed to parse file %s: %v", fname, err)
			}
			if ast.IsGenerated(f) {
				return
			}
		}

		if !strings.HasSuffix(pass.Fset.Position(n.Pos()).Filename, ".go") {
			return
		}

		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn == nil {
			return
		}

		// constructor
		if fn.Recv == nil {
			if tname, ok := isConstructor(fn); ok && enableConstructorCheck {
				hasPtr, hasVal := checkConstructorReturns(tname, fn)
				if hasPtr {
					typePtrFns[tname] = append(typePtrFns[tname], fn)
				}
				if hasVal {
					typeValFns[tname] = append(typeValFns[tname], fn)
				}
			}
			return
		}

		// method
		for _, v := range fn.Recv.List {
			if v == nil {
				continue
			}
			if v.Type == nil {
				continue
			}

			if !fnSelector.SelectFunction(fn) {
				continue
			}

			tname, isPointer := isPointer(*v)
			if isPointer {
				typePtrFns[tname] = append(typePtrFns[tname], fn)
			} else {
				typeValFns[tname] = append(typeValFns[tname], fn)
			}
		}
	})

	for tname := range mergekeys(typePtrFns, typeValFns) {
		if !((len(typePtrFns[tname]) > 0) && (len(typeValFns[tname]) > 0)) {
			continue
		}
		for _, fn := range typePtrFns[tname] {
			pass.Reportf(fn.Pos(), `%s.%s uses pointer`, tname, fn.Name.Name)
		}
		for _, fn := range typeValFns[tname] {
			pass.Reportf(fn.Pos(), `%s.%s uses value`, tname, fn.Name.Name)
		}
	}

	return nil, nil
}

func isPointer(v ast.Field) (tname string, ok bool) {
	// if it is star, then it is nested and first child is indent that has type name
	if v, ok := v.Type.(*ast.StarExpr); ok {
		if v, ok := v.X.(*ast.Ident); ok {
			return v.Name, true
		}
	}
	// if it is indent right away, it is not star, it also contains name
	if v, ok := v.Type.(*ast.Ident); ok {
		return v.Name, false
	}
	return "", false
}

func mergekeys(vs ...map[string][]*ast.FuncDecl) (keys map[string]bool) {
	keys = map[string]bool{}
	for _, v := range vs {
		for k := range v {
			keys[k] = true
		}
	}
	return keys
}
