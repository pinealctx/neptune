package main

import (
	"fmt"
	"github.com/pinealctx/neptune/pipe/grpcexample/pb"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"net"
	"os"
)

const (
	M32 = 32 * 1024 * 1024
)

func main() {
	var srv = cli.App{
		Name: "hello echo server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "addr",
			},
		},
		Action: srvRun,
	}
	var err = srv.Run(os.Args)
	if err != nil {
		fmt.Println("server error:", err)
	}
}

func srvRun(c *cli.Context) error {
	var ln, err = net.Listen("tcp", c.String("addr"))
	if err != nil {
		return err
	}
	var runner = NewSrvRunner()
	var s = grpc.NewServer(grpc.MaxSendMsgSize(M32), grpc.MaxRecvMsgSize(M32))
	pb.RegisterHelloServiceServer(s, runner)
	return s.Serve(ln)
}
