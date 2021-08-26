package main

import (
	"fmt"
	"github.com/pinealctx/neptune/lparser"
	"github.com/urfave/cli/v2"
	"go/ast"
	"os"
	"reflect"
)

func main() {
	var flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "src",
			Usage: "go file to parse",
		},
		&cli.BoolFlag{
			Name: "raw",
			Usage: "print raw source or not",
		},
	}
	var app = cli.App{
		Name:    "parse go file",
		Version: "0.1",
		Flags: flags,
		Action: parseAstFile,
	}
	var err = app.Run(os.Args)
	if err != nil {
		fmt.Println("run error:", err)
	}
}

func parseAstFile(c *cli.Context) error {
	var goFile = c.String("src")
	var astFile, src = lparser.AbsRoot(goFile)
	printAstFile(astFile, src, c.Bool("raw"))
	return nil
}

func printAstFile(a *ast.File, src []byte, raw bool) {
	fmt.Println("print ast file")
	fmt.Println("doc:")
	printCommentGroup(a.Doc, src, 1, raw)

	fmt.Println("name:")
	printIdent(a.Name, src, 1, raw)

	fmt.Println("decls:")
	printDeclList(a.Decls, src, 1, raw)

	fmt.Println("imports:")
	printImportList(a.Imports, src, 1, raw)

	fmt.Println("comments:")
	printCommentGroupList(a.Comments, src, 1, raw)
}

func printFuncDecl(d *ast.FuncDecl, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("print *ast.FuncDecl")
	if d == nil {
		printTabs(tabs)
		fmt.Println("nil")
		return
	}
	printTabs(tabs)
	fmt.Println("doc:")
	printCommentGroup(d.Doc, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("recv:")
	printFieldList(d.Recv, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("name:")
	printIdent(d.Name, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("type:")
	printFuncType(d.Type, src, tabs+1, raw)
}

func printFuncType(d *ast.FuncType, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("print *ast.FuncType")
	if d == nil {
		printTabs(tabs)
		fmt.Println("nil")
		return
	}
	printTabs(tabs)
	fmt.Println("params:")
	printFieldList(d.Params, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("result:")
	printFieldList(d.Results, src, tabs+1, raw)

	printRawBytes(d, src, tabs, raw)
}

func printFieldList(d *ast.FieldList, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("print *ast.FieldList")
	if d == nil {
		printTabs(tabs)
		fmt.Println("nil")
		return
	}
	printTabs(tabs)
	fmt.Println("lens of field list", len(d.List))
	for _, i := range d.List {
		printField(i, src, tabs+1, raw)
	}
}

func printField(d *ast.Field, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("print *ast.Field")
	if d == nil {
		printTabs(tabs)
		fmt.Println("nil")
		return
	}
	printTabs(tabs)
	fmt.Println("doc:")
	printCommentGroup(d.Doc, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("names:")
	printIdentList(d.Names, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("type expr:")
	printExpr(d.Type, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("tag:")
	printBasicList(d.Tag, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("comments:")
	printCommentGroup(d.Comment, src, tabs+1, raw)

	printRawBytes(d, src, tabs, raw)
}

func printExpr(e ast.Expr, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("print ast.Expr")
	if e == nil {
		printTabs(tabs)
		fmt.Println("nil")
		return
	}
	printTabs(tabs)
	fmt.Println("decl:", reflect.TypeOf(e), reflect.ValueOf(e))

	printRawBytes(e, src, tabs, raw)
}

func printGenDecl(d *ast.GenDecl, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("print *ast.GenDecl")
	if d == nil {
		printTabs(tabs)
		fmt.Println("nil")
		return
	}
	printTabs(tabs)
	fmt.Println("doc:")
	printCommentGroup(d.Doc, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("Tok:", d.Tok)
	printSpecList(d.Specs, src, tabs+1, raw)
}

func printSpecList(s []ast.Spec, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("print []ast.Spec")
	printTabs(tabs)
	fmt.Println("lens of ast Spec", len(s))
	for _, i := range s {
		printSpec(i, src, tabs+1, raw)
	}
}

func printSpec(s ast.Spec, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("print ast.Spec")
	if s == nil {
		printTabs(tabs)
		fmt.Println("nil")
		return
	}

	switch v := s.(type) {
	case *ast.ImportSpec:
		printImport(v, src, tabs, raw)
		return
	case *ast.TypeSpec:
		printTypeSpec(v, src, tabs, raw)
		return
	case *ast.ValueSpec:
		printValueSpec(v, src, tabs, raw)
		return
	}

	printTabs(tabs)
	fmt.Println("spec:", reflect.TypeOf(s), reflect.ValueOf(s))

	printRawBytes(s, src, tabs, raw)
}

func printValueSpec(i *ast.ValueSpec, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("doc:")
	printCommentGroup(i.Doc, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("names:")
	printIdentList(i.Names, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("type:")
	printExpr(i.Type, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("values:")
	printExpr(i.Type, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("comment:")
	printCommentGroup(i.Comment, src, tabs+1, raw)
	printRawBytes(i, src, tabs, raw)
}

func printTypeSpec(i *ast.TypeSpec, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("doc:")
	printCommentGroup(i.Doc, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("name:")
	printIdent(i.Name, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("type:")
	printExpr(i.Type, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("comment:")
	printCommentGroup(i.Comment, src, tabs+1, raw)
	printRawBytes(i, src, tabs, raw)
}

func printImportList(is []*ast.ImportSpec, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("print []*ast.ImportSpec")
	printTabs(tabs)
	fmt.Println("lens of imports", len(is))
	for _, i := range is {
		printImport(i, src, tabs+1, raw)
	}
}

func printImport(i *ast.ImportSpec, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("doc:")
	printCommentGroup(i.Doc, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("name:")
	printIdent(i.Name, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("path:")
	printBasicList(i.Path, src, tabs+1, raw)
	printTabs(tabs)
	fmt.Println("comment:")
	printCommentGroup(i.Comment, src, tabs+1, raw)
	printRawBytes(i, src, tabs, raw)
}

func printBasicList(b *ast.BasicLit, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("print st.BasicLit")
	if b == nil {
		printTabs(tabs)
		fmt.Println("nil")
		return
	}
	printTabs(tabs)
	fmt.Println("kind:", b.Kind)
	printTabs(tabs)
	fmt.Println("value:", b.Value)
	printRawBytes(b, src, tabs, raw)
}

func printCommentGroupList(gs []*ast.CommentGroup, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("print []*ast.CommentGroup")
	printTabs(tabs)
	fmt.Println("lens of comments groups", len(gs))
	for _, i := range gs {
		printCommentGroup(i, src, tabs+1, raw)
	}
}

func printCommentGroup(g *ast.CommentGroup, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("print st.CommentGroup")
	if g == nil {
		printTabs(tabs)
		fmt.Println("nil")
		return
	}
	printTabs(tabs)
	fmt.Println("group lens:", len(g.List))
	printTabs(tabs)
	fmt.Println("group txt:", g.Text())

	if raw {
		printTabs(tabs)
		fmt.Println("group value:", string(src[g.Pos()-1:g.End()-1]))
	}

	for _, i := range g.List {
		printTabs(tabs + 1)
		fmt.Println("item text:", i.Text)

		if raw {
			printTabs(tabs + 1)
			fmt.Println("item bytes:", string(src[i.Pos()-1:i.End()-1]))
		}

	}
}

func printDeclList(ds []ast.Decl, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("print []ast.Decl")
	printTabs(tabs)
	fmt.Println("lens of decls", len(ds))
	for _, i := range ds {
		printDecl(i, src, tabs+1, raw)
	}
}

func printDecl(d ast.Decl, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("print ast.Decl")
	if d == nil {
		printTabs(tabs)
		fmt.Println("nil")
		return
	}
	switch v := d.(type) {
	case *ast.GenDecl:
		printGenDecl(v, src, tabs, raw)
		return
	case *ast.FuncDecl:
		printFuncDecl(v, src, tabs, raw)
		return
	}
	printTabs(tabs)
	fmt.Println("decl:", reflect.TypeOf(d), reflect.ValueOf(d))
	printRawBytes(d, src, tabs, raw)
}

func printIdentList(ds []*ast.Ident, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("print []*ast.Ident")
	printTabs(tabs)
	fmt.Println("lens of idents", len(ds))
	for _, i := range ds {
		printIdent(i, src, tabs+1, raw)
	}
}

func printIdent(i *ast.Ident, src []byte, tabs int, raw bool) {
	printTabs(tabs)
	fmt.Println("print ast.Ident")
	if i == nil {
		printTabs(tabs)
		fmt.Println("nil")
		return
	}
	printTabs(tabs)
	fmt.Println("name:", i.Name)
	if i.Obj == nil {
		printTabs(tabs)
		fmt.Println("obj:nil")
		return
	}
	printObject(i.Obj, tabs+1)
	printRawBytes(i, src, tabs, raw)
}

func printObject(o *ast.Object, tabs int) {
	printTabs(tabs)
	fmt.Println("print *ast.Object")
	printTabs(tabs+1)
	fmt.Println("kind:", o.Kind)
	printTabs(tabs+1)
	fmt.Println("name:", o.Name)
	printTabs(tabs+1)
	fmt.Println("decl:", reflect.TypeOf(o.Decl))
	printTabs(tabs+1)
	fmt.Println("data:", reflect.TypeOf(o.Data))
	printTabs(tabs+1)
	fmt.Println("type:", reflect.TypeOf(o.Type))
}

func printRawBytes(s ast.Node, src []byte, tabs int, raw bool) {
	if raw {
		printTabs(tabs)
		fmt.Println("bytes:", string(src[s.Pos()-1:s.End()-1]))
	}
}

func printTabs(tabs int) {
	for i := 0; i < tabs; i++ {
		fmt.Print("\t")
	}
}
