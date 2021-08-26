package main

import (
	"context"
	"fmt"
	"github.com/pinealctx/neptune/pipe/grpcexample/pb"
	"github.com/urfave/cli/v2"
	"go.uber.org/atomic"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	M32 = 32 * 1024 * 1024
)

func main() {
	app := &cli.App{
		Name: "hello echo client",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "url",
				Usage: "rpc url",
			},
			&cli.IntFlag{
				Name:  "c",
				Usage: "count",
				Value: 10000,
			},
		},
		Action: run,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("run error:", err)
	}
}

func run(c *cli.Context) error {
	var conn, err = grpc.Dial(c.String("url"), grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(M32), grpc.MaxCallSendMsgSize(M32)))

	if err != nil {
		return err
	}

	var client = pb.NewHelloServiceClient(conn)
	var ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()
	var count = c.Int("c")
	var wg sync.WaitGroup
	wg.Add(count)
	var wrapE = atomic.NewError(nil)

	var t1 = time.Now()
	for i := 0; i < count; i++ {
		go func(index int) {
			defer wg.Done()
			var _, e = client.SayHello(ctx, &pb.Halo{
				Msg: strconv.Itoa(index),
			})
			if e != nil {
				wrapE.Store(e)
			}
		}(i)
	}
	wg.Wait()
	var t2 = time.Now()
	var d = t2.Sub(t1)
	fmt.Println("use time:", d, "average:", d/time.Duration(count))
	return wrapE.Load()
}
