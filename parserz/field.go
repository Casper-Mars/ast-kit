package parserz

import "go/ast"

type FieldType struct {
	astType    ast.Expr
	key        *FieldType // only for map
	x          *FieldType // for star/array/chan and the value for map
	funcType   *Func      // only for func
	importPath string
}

func NewFieldType(pkg *Pkg, importManager *ImportManager, astType ast.Expr) *FieldType {
	f := &FieldType{
		astType:    astType,
		key:        nil,
		x:          nil,
		funcType:   nil,
		importPath: "",
	}
	switch t := astType.(type) {
	case *ast.StarExpr:
		f.x = NewFieldType(pkg, importManager, t.X)
	case *ast.SelectorExpr:
		f.x = NewFieldType(pkg, importManager, t.X)
		f.importPath = importManager.GetImportPath(t.Sel.Name)
	case *ast.ArrayType:
		f.x = NewFieldType(pkg, importManager, t.Elt)
	case *ast.Ellipsis:
		f.x = NewFieldType(pkg, importManager, t.Elt)
	case *ast.MapType:
		f.key = NewFieldType(pkg, importManager, t.Key)
		f.x = NewFieldType(pkg, importManager, t.Value)
	case *ast.ChanType:
		f.x = NewFieldType(pkg, importManager, t.Value)
	case *ast.FuncType:
		f.funcType = NewFunc(pkg, importManager, "", t, nil)
	}
	return f
}

func (f *FieldType) ImportPaths() []string {
	switch f.astType.(type) {
	case *ast.Ident:
		return []string{}
	case *ast.StarExpr:
		return f.x.ImportPaths()
	case *ast.SelectorExpr:
		return []string{f.importPath}
	case *ast.ArrayType:
		return f.x.ImportPaths()
	case *ast.InterfaceType:
		return []string{}
	case *ast.MapType:
		return append(f.key.ImportPaths(), f.x.ImportPaths()...)
	case *ast.ChanType:
		return f.x.ImportPaths()
	case *ast.FuncType:
		return f.funcType.ImportPaths()
	case *ast.Ellipsis:
		return f.x.ImportPaths()
	default:
		return []string{}
	}
}

type Field struct {
	pkg      *Pkg
	astField *ast.Field
	Names    []string   // 名称，这里可能有多个，例如：a, b, c int
	Type     *FieldType // 类型
}

func NewField(pkg *Pkg, importManager *ImportManager, astField *ast.Field) *Field {
	f := &Field{
		pkg:      pkg,
		astField: astField,
		Names:    make([]string, 0, len(astField.Names)),
		Type:     NewFieldType(pkg, importManager, astField.Type),
	}
	for _, name := range astField.Names {
		f.Names = append(f.Names, name.Name)
	}
	return f
}

func (f *Field) ImportPaths() []string {
	return f.Type.ImportPaths()
}

func (f *Field) String() string {
	return f.pkg.Read(f.astField.Pos(), f.astField.End())
}
