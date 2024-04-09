package smrcptr

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"regexp"
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
	constructor       string
	skip              string
	skipGenerated     bool
	constructorRegExp *regexp.Regexp
)

func init() {
	Analyzer.Flags.StringVar(&constructor, "constructor", "^New(?P<Type>.*)", `regexp to detect constructor and type that it belongs to, if empty then skipping constructor`)
	Analyzer.Flags.StringVar(&skip, "skip", "^UnmarshalJSON$|^UnmarshalText$|^UnmarshalBinary$|^UnmarshalXML$|^UnmarshalXMLAttr$|^Scanner$|^Scan$|^Read$", `regexp to skip methods, if empty then none skipped`)
	Analyzer.Flags.BoolVar(&skipGenerated, "skip-generated", true, `skip generated files`)

	if constructor != "" {
		constructorRegExp = regexp.MustCompile(constructor)
	}
}

type functionSelector interface {
	SelectFunction(fn *ast.FuncDecl) bool
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	var fnSelector functionSelector = noopFunctionSelector{v: true}
	if skip != "" {
		fnSelector = regExpFunctionSelector{regexp: regexp.MustCompile(skip)}
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
			if constructor != "" {
				if tname, ok := isConstructor(fn); ok {
					hasPtr, hasVal := checkConstructorReturns(tname, fn)
					if hasPtr {
						typePtrFns[tname] = append(typePtrFns[tname], fn)
					}
					if hasVal {
						typeValFns[tname] = append(typeValFns[tname], fn)
					}
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
	keys = make(map[string]bool)
	for _, v := range vs {
		for k := range v {
			keys[k] = true
		}
	}
	return keys
}

type regExpFunctionSelector struct {
	regexp *regexp.Regexp
}

func (s regExpFunctionSelector) SelectFunction(fn *ast.FuncDecl) bool {
	if fn.Name == nil {
		return false
	}
	return !s.regexp.MatchString(fn.Name.Name)
}

type noopFunctionSelector struct{ v bool }

func (s noopFunctionSelector) SelectFunction(fn *ast.FuncDecl) bool { return s.v }

func isConstructor(v *ast.FuncDecl) (tname string, ok bool) {
	if v.Name == nil {
		return "", false
	}
	matches := constructorRegExp.FindStringSubmatch(v.Name.Name)
	if idx := constructorRegExp.SubexpIndex("Type"); len(matches) > 0 && idx >= 0 && matches[idx] != "" {
		return matches[idx], true
	}
	return "", false
}

func checkConstructorReturns(name string, fn *ast.FuncDecl) (hasPtr, hasVal bool) {
	if fn.Type == nil {
		return false, false
	}
	if fn.Type.Results == nil {
		return false, false
	}
	if len(fn.Type.Results.List) == 0 {
		return false, false
	}
	for _, q := range fn.Type.Results.List {
		if q == nil {
			continue
		}
		tname, isPointer := isPointer(*q)
		if tname != name {
			continue
		}
		if isPointer {
			hasPtr = true
		} else {
			hasVal = true
		}
	}
	return hasPtr, hasVal
}
