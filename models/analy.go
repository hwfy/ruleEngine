package models

import (
	"context"
	"time"

	"github.com/astaxie/beego"
	"github.com/coreos/etcd/clientv3"
	"github.com/hprose/hprose-golang/rpc"
)

type (
	Analy struct {
		Run func([]byte) (string, error)
	}
	//provided to the document display
	data struct {
		name string
		data []string
	}
)

func GetResult(data []byte) (string, error) {
	etcdclient, err := clientv3.New(clientv3.Config{
		Endpoints:   beego.AppConfig.Strings("etcdAddrs"),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return "", err
	}
	defer etcdclient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	resp, err := etcdclient.Get(ctx, "ruleAnaly")
	cancel()
	if err != nil {
		return "", err
	}

	client := rpc.NewClient(string(resp.Kvs[0].Value))

	var analy *Analy
	client.UseService(&analy)

	return analy.Run(data)
}
