package parserz

import (
	"go/ast"
	"strings"
)

type FieldType struct {
	pkg        *Pkg
	astType    ast.Expr
	key        *FieldType // only for map
	x          *FieldType // for star/array/chan and the value for map
	funcType   *Func      // only for func
	importPath string
}

func NewFieldType(pkg *Pkg, importManager ImportManager, astType ast.Expr) *FieldType {
	f := &FieldType{
		pkg:        pkg,
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
		f.x = NewFieldType(pkg, importManager, t.Sel)
		f.importPath = importManager.GetImportByAlias(t.X.(*ast.Ident).Name).Path
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

func (f *FieldType) Name() string {
	switch t := f.astType.(type) {
	default:
		return ""
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return f.x.Name()
	case *ast.SelectorExpr:
		return f.x.Name()
	case *ast.ArrayType:
		return f.x.Name()
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.MapType:
		return "map"
	case *ast.ChanType:
		return f.x.Name()
	case *ast.FuncType:
		return "func"
	case *ast.Ellipsis:
		return f.x.Name()
	}
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

func (f *FieldType) String() string {
	switch t := f.astType.(type) {
	default:
		return ""
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + f.x.String()
	case *ast.SelectorExpr:
		return f.pkg.Imports.GetImportByPath(f.importPath).Alias + "." + f.x.String()
	case *ast.ArrayType:
		return "[]" + f.x.String()
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.MapType:
		return "map[" + f.key.String() + "]" + f.x.String()
	case *ast.ChanType:
		if t.Dir == ast.SEND {
			return "chan<- " + f.x.String()
		}
		if t.Dir == ast.RECV {
			return "<-chan " + f.x.String()
		}
		return "chan " + f.x.String()
	case *ast.FuncType:
		str := func(t *Func) (ret string) {
			ret = "("
			// 拼接参数
			for index, p := range t.Params {
				ret += p.String()
				if index < len(t.Params)-1 {
					ret += ", "
				}
			}
			ret += ")"

			// 拼接返回值
			if len(t.Results) > 0 {
				ret += " "
				// 多返回值时用括号包裹
				if len(t.Results) > 1 {
					ret += "("
					defer func() {
						ret += ")"
					}()
				}

				for index, p := range t.Results {
					ret += p.String()
					if index < len(t.Results)-1 {
						ret += ", "
					}
				}
			}

			return ret
		}(f.funcType)
		return "func" + str
	case *ast.Ellipsis:
		return "..." + f.x.String()
	}
}

type Field struct {
	pkg      *Pkg
	astField *ast.Field
	Names    []string   // 名称，这里可能有多个，例如：a, b, c int
	Type     *FieldType // 类型
}

func NewField(pkg *Pkg, importManager ImportManager, astField *ast.Field) *Field {
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
	builder := strings.Builder{}
	if len(f.Names) != 0 {
		builder.WriteString(strings.Join(f.Names, ", "))
		builder.WriteString(" ")
	}
	builder.WriteString(f.Type.String())
	return builder.String()
}
