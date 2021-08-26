package main

import (
	"errors"
	"github.com/pinealctx/neptune/store/etcd"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var (
	eCli *etcd.Client
)

func main() {
	var app = cli.App{
		Name:    "test etcd cli",
		Version: "0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "url",
			},
			&cli.StringFlag{
				Name:  "root",
				Usage: "root path",
			},
			&cli.StringFlag{
				Name:  "mode",
				Usage: "create/delete/deleteDir/update/get/getDir/watchDir",
			},
			&cli.StringFlag{
				Name:  "key",
				Usage: "key",
			},
			&cli.StringFlag{
				Name:  "content",
				Usage: "content",
			},
			&cli.StringFlag{
				Name:  "file",
				Usage: "read file path",
			},
			&cli.IntFlag{
				Name:  "dev",
				Usage: "data version",
			},
			&cli.BoolFlag{
				Name:  "full",
				Usage: "get full dir or not(including data and version)",
			},
		},
		Action: runCmd,
	}
	var err = app.Run(os.Args)
	if err != nil {
		log.Println("run command error:", err)
	} else {
		log.Println("run command ok")
	}
}

func runCmd(c *cli.Context) error {
	var err error
	eCli, err = etcd.NewClient(c.String("url"), c.String("root"))
	if err != nil {
		return err
	}

	switch c.String("mode") {
	case "create":
		return createNode(c)
	case "delete":
		return deleteNode(c)
	case "deleteDir":
		return deleteDir(c)
	case "update":
		return updateNode(c)
	case "get":
		return getNode(c)
	case "getDir":
		return getDir(c)
	case "watchDir":
		return watchDir(c)
	}
	return errors.New("unsupported.mode")
}
