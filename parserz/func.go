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

func NewFunc(pkg *Pkg, importManager ImportManager, name string, astFunc *ast.FuncType, rec *ast.Field) *Func {
	f := &Func{
		pkg:     pkg,
		astFunc: astFunc,
		name:    name,
	}
	if rec != nil {
		f.Rec = NewField(pkg, importManager, rec)
	}
	if len(astFunc.Params.List) != 0 {
		f.Params = make([]*Field, 0, len(astFunc.Params.List))
		for _, param := range astFunc.Params.List {
			f.Params = append(f.Params, NewField(pkg, importManager, param))
		}
	}
	if len(astFunc.Results.List) != 0 {
		f.Results = make([]*Field, 0, len(astFunc.Results.List))
		for _, result := range astFunc.Results.List {
			f.Results = append(f.Results, NewField(pkg, importManager, result))
		}
	}
	return f
}

func (i *Func) Name() string {
	return i.name
}

//InterfaceFormat 转成接口方法格式
func (i *Func) InterfaceFormat() string {
	return ""
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
