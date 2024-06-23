package main

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	defer cli.Close()
	fmt.Println("connect to etcd success")

	// Put
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// _, err = cli.Put(ctx, "q1mi", "dsb")
	// cancel()
	// if err != nil {
	// 	fmt.Printf("put to etcd failed, err:%v\n", err)

	// 	return
	// }
	// fmt.Println("put to etcd success")

	// Get
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := cli.Get(ctx, "/TikTok-config/favorite")
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
}
