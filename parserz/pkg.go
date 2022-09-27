package parserz

import "go/token"

type Pkg struct {
	//Vars   []*Field
	//Const  []*Field
	Struct []*Struct
	Func   []*Func
}

func (p *Pkg) Read(pos, end token.Pos) string {
	return ""
}

type PkgBuilder struct {
	path string
}

func NewPkgBuilder() *PkgBuilder {
	return &PkgBuilder{}
}

func (p *PkgBuilder) WithPkgPath(path string) *PkgBuilder {
	p.path = path
	return p
}

func (p *PkgBuilder) Build() *Pkg {

	return &Pkg{}
}
