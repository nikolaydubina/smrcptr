package smrcptr

import (
	"go/ast"
	"strings"
)

const (
	constructorPrefix = "New"
)

func isConstructor(v *ast.FuncDecl) (tname string, ok bool) {
	if v.Name == nil {
		return "", false
	}
	if !strings.HasPrefix(v.Name.Name, constructorPrefix) {
		return "", false
	}
	return v.Name.Name[len(constructorPrefix):], true
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
