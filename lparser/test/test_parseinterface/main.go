package main

import (
	"fmt"
	"github.com/pinealctx/neptune/lparser"
	"github.com/urfave/cli/v2"
	"go/ast"
	"os"
)

func main() {
	var app = cli.App{
		Name:    "parse go interface",
		Version: "0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "src",
				Usage: "go file to parse",
			},
		},
		Action: parseInterface,
	}
	var err = app.Run(os.Args)
	if err != nil {
		fmt.Println("run error:", err)
	}
}

func parseInterface(c *cli.Context) error {
	var goFile = c.String("src")
	var astFile, src = lparser.AbsRoot(goFile)
	printInterface(astFile, src)
	return nil
}

func printInterface(a *ast.File, src []byte) {
	var interfaceList = lparser.AbsInterfaces(a, src)
	for _, e := range interfaceList {
		fmt.Println(e)
	}
}
