package rpcserver

import (
	"context"
	"fmt"
	"net"

	"github.com/HideInBush7/go-server/example/order/pb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type OrderServer struct {
	pb.UnimplementedOrderServer
}

type order struct {
	Name   string
	Id     string
	UserId string
}

var orders = []order{
	{"order 用户1", "1", "1"},
	{"order 用户1", "11", "1"},
	{"order 用户2", "2", "2"},
	{"order 用户2", "22", "2"},
}

func (o *OrderServer) GetOrderList(c context.Context, req *pb.ListRequest) (*pb.OrderResponse, error) {
	fmt.Println(`GetOrderList()`)
	var res []order
	for _, v := range orders {
		if v.UserId == req.UserId {
			res = append(res, v)
		}
	}
	return &pb.OrderResponse{
		Data: fmt.Sprintf("%+v", res),
	}, nil
}

func (o *OrderServer) GetOrderInfo(c context.Context, req *pb.InfoRequest) (*pb.OrderResponse, error) {
	for _, v := range orders {
		if v.Id == req.Id {
			return &pb.OrderResponse{
				Data: fmt.Sprintf("%+v", v),
			}, nil
		}
	}
	return &pb.OrderResponse{
		Data: "nil",
	}, nil
}

func NewServer() {
	lis, err := net.Listen("tcp", net.JoinHostPort(viper.GetString(`ip`), viper.GetString(`port`)))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterOrderServer(s, &OrderServer{})
	fmt.Println("[RPC] Server start... Listenning->", net.JoinHostPort(viper.GetString(`ip`), viper.GetString(`port`)))

	err = s.Serve(lis)

	if err != nil {
		panic(err)
	}
}
