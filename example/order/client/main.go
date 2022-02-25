package main

import (
	"context"
	"fmt"
	"time"

	"github.com/HideInBush7/go-server/config"
	"github.com/HideInBush7/go-server/etcd"
	"github.com/HideInBush7/go-server/example/order/pb"
	"github.com/HideInBush7/go-server/example/order/rpcclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

func main() {
	config.Init()
	c, err := etcd.DefaultClient()
	if err != nil {
		panic(err)
	}
	b, err := etcd.NewBuilder(c)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(`etcd:///order.rpc`,
		grpc.WithResolvers(b),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(time.Second * 3)
		c := rpcclient.NewOrder(conn)
		res, err := c.GetOrderList(context.Background(), &pb.ListRequest{
			UserId: "1",
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(res.String())
	}
}
