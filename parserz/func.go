package parserz

import "go/ast"

type Func struct {
	pkg     *Pkg
	astFunc *ast.FuncType
	name    string
	Rec     *Field
	Params  []*Field
	Results []*Field
}

func (i *Func) Name() string {
	return i.name
}

func (i *Func) ImportPaths() []string {
	var paths []string
	if i.Rec != nil {
		paths = append(paths, i.Rec.ImportPaths()...)
	}
	if len(i.Params) > 0 {
		for _, param := range i.Params {
			paths = append(paths, param.ImportPaths()...)
		}
	}
	if len(i.Results) > 0 {
		for _, result := range i.Results {
			paths = append(paths, result.ImportPaths()...)
		}
	}
	return paths
}
