package parserz

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type Pkg struct {
	//Vars   []*Field
	//Const  []*Field
	Struct  []*Struct
	Func    []*Func
	fileSet *token.FileSet
	fileSrc map[string][]byte
	name    string
}

func (p *Pkg) Name() string {
	return p.name
}

func (p *Pkg) Read(pos, end token.Pos) string {
	return ""
}

type PkgBuilder struct {
	path   string
	filter func(info os.FileInfo) bool
}

func NewPkgBuilder() *PkgBuilder {
	return &PkgBuilder{
		path: ".",
		filter: func(info os.FileInfo) bool {
			return strings.HasSuffix(info.Name(), ".go") && !(info.IsDir() || strings.HasSuffix(info.Name(), "_test.go") || strings.HasSuffix(info.Name(), "gen.go"))
		},
	}
}

func (p *PkgBuilder) WithPkgPath(path string) *PkgBuilder {
	p.path = path
	return p
}

func (p *PkgBuilder) Build() (*Pkg, error) {
	fileSet := token.NewFileSet()
	pkgSet, err := parser.ParseDir(fileSet, p.path, func(info os.FileInfo) bool {
		if !p.filter(info) {
			return false
		}
		return true
	}, 0)
	if err != nil {
		return nil, err
	}
	result := &Pkg{
		fileSet: fileSet,
	}
	for pkgName, astPkg := range pkgSet {
		result.fileSrc = make(map[string][]byte, len(astPkg.Files))
		result.name = pkgName
		for _, file := range astPkg.Files {
			readFile, err := os.ReadFile(file.Name.String())
			importManger := NewImportManger(file)
			if err != nil {
				fmt.Printf("error reading file %s: %s", file.Name.String(), err)
				return nil, err
			}
			result.fileSrc[file.Name.String()] = readFile
			for _, decl := range file.Decls {
				switch astDecl := decl.(type) {
				case *ast.FuncDecl:
					result.Func = append(result.Func, NewFunc(result, importManger, astDecl))
				case *ast.GenDecl:
					for _, spec := range astDecl.Specs {
						switch astSpec := spec.(type) {
						case *ast.TypeSpec:
							structType, ok := astSpec.Type.(*ast.StructType)
							if ok {
								result.Struct = append(result.Struct, NewStruct(result, importManger, astSpec.Name.String(), structType))
							}
						}
					}
				}
			}
		}
		break
	}
	return result, nil
}
