package main

import (
	"fmt"
	"github.com/pinealctx/neptune/lparser"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	var app = cli.App{
		Name:    "parse go imports",
		Version: "0.1",
		Flags:  []cli.Flag{
			&cli.StringFlag{
				Name:  "src",
				Usage: "go file to parse",
			},
		},
		Action: parseImports,
	}
	var err = app.Run(os.Args)
	if err != nil {
		fmt.Println("run error:", err)
	}
}

func parseImports(c *cli.Context) error {
	var goFile = c.String("src")
	var astFile, _ = lparser.AbsRoot(goFile)
	fmt.Println(lparser.AbsImports(astFile))
	return nil
}
