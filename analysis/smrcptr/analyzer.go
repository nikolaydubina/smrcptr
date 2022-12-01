package smrcptr

import (
	"flag"
	"go/ast"
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
)

func init() {
	Analyzer.Flags.BoolVar(&enableConstructorCheck, "constructor", false, `enable constructor return type check`)
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

        // non-nil pointers for memory efficiency
	typePtrFns := map[string][]*ast.FuncDecl{}
	typeValFns := map[string][]*ast.FuncDecl{}

	inspect.Preorder([]ast.Node{&ast.FuncDecl{}}, func(n ast.Node) {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn == nil {
			return
		}

		if !strings.HasSuffix(pass.Fset.Position(fn.Pos()).Filename, ".go") {
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
