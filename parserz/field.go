package parserz

import "go/ast"

type FieldType struct {
	astType    ast.Expr
	key        *FieldType // only for map
	x          *FieldType // for star/array/chan and the value for map
	funcType   *Func      // only for func
	importPath string
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

func (f *Field) ImportPaths() []string {
	return f.Type.ImportPaths()
}

func (f *Field) String() string {
	return f.pkg.Read(f.astField.Pos(), f.astField.End())
}
