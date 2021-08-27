package main

import (
	"context"
	"log"
	"time"

	"github.com/pinealctx/neptune/store/etcd"
	"github.com/urfave/cli/v2"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
)

func createNode(c *cli.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var dRet = eCli.Create(ctx, c.String("key"), []byte(c.String("content")))
	if dRet.Err != nil {
		return dRet.Err
	}
	log.Println("create node:", c.String("key"), " version:", dRet.Revision)
	return nil
}

func deleteNode(c *cli.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var revision = int64(c.Int("dev"))
	var dRet = eCli.Delete(ctx, c.String("key"), revision)
	if dRet.Err != nil {
		return dRet.Err
	}
	log.Println("delete node:", c.String("key"), " modification version:", dRet.Revision)
	return nil
}

func deleteDir(c *cli.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var revision = int64(c.Int("dev"))
	var dRet = eCli.DeleteDir(ctx, c.String("key"), revision)
	if dRet.Err != nil {
		return dRet.Err
	}
	log.Println("delete dir:", c.String("key"), " modification version:", dRet.Revision)
	return nil
}

func updateNode(c *cli.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var revision = int64(c.Int("dev"))
	var dRet = eCli.Put(ctx, c.String("key"), []byte(c.String("content")), revision)
	if dRet.Err != nil {
		return dRet.Err
	}
	log.Println("update node:", c.String("key"), " modification version:", dRet.Revision)
	return nil
}

func getNode(c *cli.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var dRet = eCli.Get(ctx, c.String("key"))
	if dRet.Err != nil {
		return dRet.Err
	}
	log.Println("get node:", c.String("key"), " revision:", dRet.Revision, " data:", string(dRet.Data))
	return nil
}

func getDir(c *cli.Context) error {
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var dRet = eCli.GetDir(ctx, c.String("key"), !c.Bool("full"))
	if dRet.Err != nil {
		return dRet.Err
	}
	log.Println("get dir node:", c.String("key"), " modification version:", dRet.Revision)
	for _, i := range dRet.KVS {
		printItem(i)
	}
	return nil
}

func watchDir(c *cli.Context) error {
	var watcher = etcd.NewWatcher(eCli, c.String("key"), time.Second*3, time.Second*3)
	watcher.StartWatchDir()
	for kEvent := range watcher.DirChan() {
		log.Println("watched:")
		log.Println("version:", kEvent.Revision, " error:", kEvent.Err)
		for i := range kEvent.Events {
			printEvent(kEvent.Events[i])
		}
		log.Println("")
	}
	return nil
}

func printItem(i *mvccpb.KeyValue) {
	log.Println("key:", i.Key, " version:", i.Version, "mod version", i.ModRevision, " value:", string(i.Value))
}

func printEvent(i *clientv3.Event) {
	log.Println("type:", i.Type, " k:", string(i.Kv.Key), " v:", string(i.Kv.Value))
	if i.PrevKv != nil {
		log.Println("previous k:", string(i.PrevKv.Key), "previous v:", string(i.PrevKv.Value))
	}
}
