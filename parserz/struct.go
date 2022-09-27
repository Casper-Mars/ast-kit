package parserz

import "go/ast"

type Struct struct {
	pkg    *Pkg
	Name   string
	Fields []*Field
}

func NewStruct(pkg *Pkg, importManager *ImportManager, name string, astSpec *ast.StructType) *Struct {
	s := &Struct{
		pkg:    pkg,
		Name:   name,
		Fields: make([]*Field, 0, len(astSpec.Fields.List)),
	}
	for _, field := range astSpec.Fields.List {
		s.Fields = append(s.Fields, NewField(pkg, importManager, field))
	}
	return s
}
