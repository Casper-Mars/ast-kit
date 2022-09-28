package parserz

import (
	"go/ast"
	"strings"
)

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
	if astFunc.Params != nil && len(astFunc.Params.List) != 0 {
		f.Params = make([]*Field, 0, len(astFunc.Params.List))
		for _, param := range astFunc.Params.List {
			f.Params = append(f.Params, NewField(pkg, importManager, param))
		}
	}
	if astFunc.Results != nil && len(astFunc.Results.List) != 0 {
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
	builder := strings.Builder{}
	builder.WriteString(i.name)
	builder.WriteByte('(')
	// construct param
	if len(i.Params) > 0 {
		for num, param := range i.Params {
			builder.WriteString(param.String())
			if num+1 < len(i.Params) {
				builder.WriteString(", ")
			}
		}
	}
	builder.WriteString(") ")
	// construct result
	if len(i.Results) > 0 {
		if len(i.Results) > 1 || len(i.Results[0].Names) > 0 {
			builder.WriteByte('(')
		}
		for num, result := range i.Results {
			builder.WriteString(result.String())
			if num+1 < len(i.Results) {
				builder.WriteByte(',')
			}
		}
		if len(i.Results) > 1 || len(i.Results[0].Names) > 0 {
			builder.WriteByte(')')
		}
	}
	return builder.String()
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
