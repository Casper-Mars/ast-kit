package parserz

import (
	"go/ast"
	"strings"
)

type ImportManager struct {
	importMap map[string]string
}

func NewImportManger(astFile *ast.File) *ImportManager {
	i := &ImportManager{
		importMap: make(map[string]string),
	}
	for _, imp := range astFile.Imports {
		if imp.Name == nil {
			split := strings.Split(strings.Trim(imp.Path.Value, "\""), "/")
			i.importMap[split[len(split)-1]] = imp.Path.Value
		} else {
			i.importMap[imp.Name.Name] = imp.Path.Value
		}
	}
	return i
}

func (i *ImportManager) GetImportPath(alias string) string {
	return i.importMap[alias]
}
