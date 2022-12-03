package smrcptr

import (
	"go/ast"
)

type mapNameFunctionSelector struct {
	fns map[string]bool
	def bool
}

func (s mapNameFunctionSelector) SelectFunction(fn *ast.FuncDecl) bool {
	if fn.Name == nil {
		return false
	}
	v, ok := s.fns[fn.Name.Name]
	if ok {
		return v
	}
	return s.def
}
