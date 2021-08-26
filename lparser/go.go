package lparser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
)

//AbsRoot -- figure abs tree root
//figure root and doc package
//return root and doc.Packages
func AbsRoot(fileName string) (*ast.File, []byte) {
	var content, err = ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	var fSet = token.NewFileSet()
	var root *ast.File
	root, err = parser.ParseFile(fSet, fileName, content, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	return root, content
}

//ImportField -- import field
//Name -- the import alias, could be empty
//Path -- the import path
//For instance
/*
import (
	aliasA "github.com/xx/a"
	"github.com/xx/b"
)
===>>>
{Name:`aliasA`, Path:`"github.com/xx/a"`}, {Name:``, Path:`"github.com/xx/b"`}.
First item "Name" field is `aliasA` not empty.
Second item "Name" field is  empty.
Every item "Path" field has wrapped by `""`, like `"github.com/xx/a"` not `github.com/xx/a`.
*/
type ImportField struct {
	Name string
	Path string
}

//AbsImports -- figure out imports
func AbsImports(f *ast.File) []ImportField {
	var c = len(f.Imports)
	if c == 0 {
		return nil
	}
	var fs = make([]ImportField, c)
	for i := 0; i < c; i++ {
		if f.Imports[i].Name != nil {
			fs[i].Name = f.Imports[i].Name.Name
		}
		fs[i].Path = f.Imports[i].Path.Value
	}
	return fs
}

//Method : method in interface, actually it's function.
type Method struct {
	//Params input params
	Params []string
	//Results return params
	Results []string
	//Name method name
	Name string
	//Doc method comment
	Doc string
}

//Interface -- interface field
type Interface struct {
	Name    string
	Methods []Method
}

//AbsInterfaces -- figure out interfaces
func AbsInterfaces(f *ast.File, src []byte) []Interface {
	var typeList = filterInterfaces(f)
	var c = len(typeList)
	if c == 0 {
		return nil
	}
	var s = make([]Interface, 0, c)
	for _, e := range typeList {
		var i = convertInterfaceType(e, src)
		s = append(s, i)
	}
	return s
}

//convert interface type
func convertInterfaceType(t *ast.TypeSpec, src []byte) Interface {
	var i Interface
	i.Name = t.Name.String()
	var e = t.Type.(*ast.InterfaceType)
	if e.Methods == nil {
		return i
	}
	var c = len(e.Methods.List)
	if c == 0 {
		return i
	}
	i.Methods = make([]Method, 0, c)
	for _, m := range e.Methods.List {
		var method = convertMethod(m, src)
		i.Methods = append(i.Methods, method)
	}
	return i
}

//filterInterfaces -- figure out interfaces
func filterInterfaces(f *ast.File) []*ast.TypeSpec {
	var c = len(f.Decls)
	if c == 0 {
		return nil
	}
	var typeSpecs = make([]*ast.TypeSpec, 0, c)
	for _, i := range f.Decls {
		var g, ok = i.(*ast.GenDecl)
		if !ok {
			continue
		}
		if len(g.Specs) != 1 {
			continue
		}
		var t *ast.TypeSpec
		t, ok = g.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}
		_, ok = t.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}
		typeSpecs = append(typeSpecs, t)
	}
	return typeSpecs
}

//convert method
func convertMethod(e *ast.Field, src []byte) Method {
	var method Method
	method.Name = e.Names[0].String()
	method.Doc = e.Doc.Text()
	var f = e.Type.(*ast.FuncType)
	if f == nil {
		return method
	}
	if f.Params != nil {
		var c = len(f.Params.List)
		if c > 0 {
			method.Params = make([]string, 0, c)
			for _, k := range f.Params.List {
				method.Params = append(method.Params, string(src[k.Type.Pos()-1:k.Type.End()-1]))
			}
		}
	}
	if f.Results != nil {
		var c = len(f.Results.List)
		if c > 0 {
			method.Results = make([]string, 0, c)
			for _, k := range f.Results.List {
				method.Results = append(method.Results, string(src[k.Type.Pos()-1:k.Type.End()-1]))
			}
		}
	}
	return method
}
