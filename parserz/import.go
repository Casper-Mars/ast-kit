package parserz

import (
	"go/ast"
	"strconv"
	"strings"
)

type Import struct {
	Path  string
	Alias string
}

type ImportManager interface {
	All() []*Import
	GetImportByPath(path string) *Import
	GetImportByAlias(alias string) *Import
}

type fileImportManager struct {
	importMap map[string]*Import
	all       []*Import
}

func (i *fileImportManager) All() []*Import {
	return i.all
}

func (i *fileImportManager) GetImportByPath(path string) *Import {
	return &Import{
		Path: path,
	}
}

func (i *fileImportManager) GetImportByAlias(alias string) *Import {
	result, ok := i.importMap[alias]
	if !ok {
		return &Import{}
	}
	return result
}

func newFileImportManger(astFile *ast.File) ImportManager {
	i := &fileImportManager{
		importMap: make(map[string]*Import),
	}
	for _, imp := range astFile.Imports {
		item := &Import{
			Path: imp.Path.Value,
		}
		if imp.Name == nil {
			split := strings.Split(strings.Trim(imp.Path.Value, "\""), "/")
			item.Alias = split[len(split)-1]
		} else {
			item.Alias = imp.Name.Name
		}
		i.all = append(i.all, item)
		i.importMap[item.Alias] = item
	}
	return i
}

type pkgImportManager struct {
	all      []*Import
	pathMap  map[string]*Import
	aliasMap map[string]*Import
}

func newPkgImportManager() *pkgImportManager {
	return &pkgImportManager{
		pathMap:  make(map[string]*Import),
		aliasMap: make(map[string]*Import),
	}
}

func (p *pkgImportManager) add(list []*Import) {
	for _, item := range list {
		if _, ok := p.pathMap[item.Path]; ok {
			continue
		}
		p.all = append(p.all, item)
		p.pathMap[item.Path] = item
		if _, ok := p.aliasMap[item.Alias]; !ok {
			p.aliasMap[item.Alias] = item
		} else {
			for i := 1; ; i++ {
				s := item.Alias + strconv.Itoa(i)
				if _, ok := p.aliasMap[s]; !ok {
					p.aliasMap[s] = item
					item.Alias = s
					break
				}
			}
		}
	}
}

func (p *pkgImportManager) All() []*Import {
	return p.all
}

func (p *pkgImportManager) GetImportByPath(path string) *Import {
	return p.pathMap[path]
}

func (p *pkgImportManager) GetImportByAlias(alias string) *Import {
	return p.aliasMap[alias]
}
