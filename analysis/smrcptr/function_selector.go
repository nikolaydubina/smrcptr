package smrcptr

import "go/ast"

type allFunctions struct{}

func (allFunctions) SelectFunction(fn *ast.FuncDecl) bool { return true }

type andFunctions struct {
	vs []functionSelector
}

func (v andFunctions) SelectFunction(fn *ast.FuncDecl) bool {
	for _, q := range v.vs {
		if !q.SelectFunction(fn) {
			return false
		}
	}
	return true
}

type nameFunction struct {
	name string
}

func (v nameFunction) SelectFunction(fn *ast.FuncDecl) bool {
	if fn.Name == nil {
		return false
	}
	return fn.Name.Name == v.name
}

type notFunction struct {
	v functionSelector
}

func (v notFunction) SelectFunction(fn *ast.FuncDecl) bool { return !v.v.SelectFunction(fn) }
