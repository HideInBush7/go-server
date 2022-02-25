package rpcclient

import (
	"context"

	"github.com/HideInBush7/go-server/example/order/pb"
	"google.golang.org/grpc"
)

type defaultOrder struct {
	client grpc.ClientConnInterface
}

func (d *defaultOrder) GetOrderList(ctx context.Context, in *pb.ListRequest) (*pb.OrderResponse, error) {
	cli := pb.NewOrderClient(d.client)
	return cli.GetOrderList(ctx, in)
}

func (d *defaultOrder) GetOrderInfo(ctx context.Context, in *pb.InfoRequest) (*pb.OrderResponse, error) {
	cli := pb.NewOrderClient(d.client)
	return cli.GetOrderInfo(ctx, in)
}

func NewOrder(conn grpc.ClientConnInterface) *defaultOrder {
	return &defaultOrder{
		client: conn,
	}
}
