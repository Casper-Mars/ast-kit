package parserz

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
)

type Pkg struct {
	//Vars   []*Field
	//Const  []*Field
	Imports ImportManager
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
	s := p.fileSet.Position(pos)
	e := p.fileSet.Position(end)
	file := p.fileSrc[s.Filename]
	return string(file[s.Offset:e.Offset])
}

type PkgBuilder struct {
	path   string
	filter func(info os.FileInfo) bool
}

func NewPkgBuilder(options ...Option) *PkgBuilder {
	p := &PkgBuilder{
		path: ".",
		filter: func(info os.FileInfo) bool {
			return strings.HasSuffix(info.Name(), ".go") && !(info.IsDir() || strings.HasSuffix(info.Name(), "_test.go") || strings.HasSuffix(info.Name(), "gen.go"))
		},
	}
	for _, option := range options {
		option(p)
	}
	return p
}

// Deprecated: use WithPath instead
func (p *PkgBuilder) WithPkgPath(path string) *PkgBuilder {
	p.path = path
	return p
}

func (p *PkgBuilder) Build() (*Pkg, error) {
	fileSet := token.NewFileSet()
	pkgSet, err := parser.ParseDir(fileSet, p.path, p.filter, 0)
	if err != nil {
		return nil, err
	}
	result := &Pkg{
		fileSet: fileSet,
	}
	for pkgName, astPkg := range pkgSet {
		result.fileSrc = make(map[string][]byte, len(astPkg.Files))
		result.name = pkgName
		pkgImportManager := newPkgImportManager()
		result.Imports = pkgImportManager
		for filename, file := range astPkg.Files {
			importManger := newFileImportManger(file)
			pkgImportManager.add(importManger.All())
			readFile, err := ioutil.ReadFile(filename)
			if err != nil {
				fmt.Printf("error reading file %s: %s\n", filename, err)
				return nil, err
			}
			result.fileSrc[filename] = readFile
			for _, decl := range file.Decls {
				switch astDecl := decl.(type) {
				case *ast.FuncDecl:
					// add func
					var rec *ast.Field
					if astDecl.Recv != nil {
						rec = astDecl.Recv.List[0]
					}
					result.Func = append(result.Func, NewFunc(result, importManger, astDecl.Name.Name, astDecl.Type, rec))
				case *ast.GenDecl:
					for _, spec := range astDecl.Specs {
						switch astSpec := spec.(type) {
						case *ast.TypeSpec:
							// add struct
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
