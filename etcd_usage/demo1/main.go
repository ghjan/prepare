package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"prepare/etcd_usage"
)

func main() {
	var (
		client    *clientv3.Client
		err       error
		ctx       context.Context
		presponse *clientv3.PutResponse
		gresponse *clientv3.GetResponse
		kv        clientv3.KV
		dresponse *clientv3.DeleteResponse
	)

	//建立连接
	endpoints := []string{"192.168.1.133:2379",}
	if client, err = etcd_usage.GetEtcdClient(endpoints); err != nil {
		fmt.Println(err)
		return
	}
	ctx = context.TODO()
	kv = clientv3.NewKV(client)
	presponse, err = kv.Put(ctx, "/cron/jobs/job1", "hello", clientv3.WithPrevKV())
	//presponse, err = kv.Put(ctx, "/cron/jobs/job1", "hello")
	fmt.Printf("presponse.Header.Revision:%d\n", presponse.Header.Revision)
	fmt.Printf("presponse.Header.RaftTerm:%d\n", presponse.Header.RaftTerm)
	fmt.Printf("presponse.PrevKv:%s\n", presponse.PrevKv)

	presponse, err = kv.Put(ctx, "/cron/jobs/job1", "hello2", clientv3.WithPrevKV())
	//presponse, err = kv.Put(ctx, "/cron/jobs/job1", "hello2")

	fmt.Printf("presponse.Header.Revision:%d\n", presponse.Header.Revision)
	fmt.Printf("presponse.Header.RaftTerm:%d\n", presponse.Header.RaftTerm)
	fmt.Printf("presponse.PrevKv:%s\n", presponse.PrevKv)

	gresponse, err = kv.Get(ctx, "/cron/jobs/job1")
	fmt.Println(gresponse.Count)
	fmt.Println(gresponse.Kvs)
	gresponse, err = kv.Get(ctx, "/cron/jobs/", clientv3.WithPrefix())
	fmt.Println(gresponse.Count)
	fmt.Println(gresponse.Kvs)

	dresponse, err = kv.Delete(ctx, "/cron/jobs/job1", clientv3.WithPrevKV())
	fmt.Printf("dresponse.Deleted:%d\n", dresponse.Deleted)
	fmt.Printf("dresponse.PrevKvs:%s\n", dresponse.PrevKvs)

}
