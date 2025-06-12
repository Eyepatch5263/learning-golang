package order

import (
	"context"
	"time"

	"github.com/eyepatch5263/go-grpc-microservices/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := pb.NewOrderServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostOrder(ctx context.Context,accountID string, products []OrderedProduct)(*Order,error){
	protoProduct:=[]*pb.PostOrderRequest_OrderedProduct{}
	for _,p:=range products{
		protoProduct=append(protoProduct, &pb.PostOrderRequest_OrderedProduct{
			ProductId: p.ID,
			Quantity: p.Quantity,
		})
	}
	r,err:=c.service.PostOrder(ctx,&pb.PostOrderRequest{
		AccountId: accountID,
		Products: protoProduct,
	})
	if err!=nil{
		return nil,err
	}
	newOrder:=r.Order
	newOderCreatedAt:=time.Time{}
	newOderCreatedAt.UnmarshalBinary(newOrder.CreatedAt)
	return &Order{
		ID: newOrder.Id,
		AccountID: newOrder.AccountId,
		TotalPrice: newOrder.TotalPrice,
		Products: products,
	},nil
}

func (c *Client) GetOrdersForAccount(ctx context.Context, accountID string)([]Order,error){
	r,err:=c.service.GetOrdersForAccount(ctx,&pb.GetOrdersForAccountRequest{
		AccountId: accountID,
	})
	if err!=nil{
		return nil,err
	}
	orders:=[]Order{}
	for _,orderProto:=range r.Orders{
		newOrder:=Order{
			ID: orderProto.Id,
			AccountID: orderProto.AccountId,
			TotalPrice: orderProto.TotalPrice,
		}
		newOrder.CreatedAt=time.Time{}
		newOrder.CreatedAt.UnmarshalBinary(orderProto.CreatedAt)
		products:=[]OrderedProduct{}
		for _,p:=range orderProto.Products{
			products=append(products, OrderedProduct{
				ID:p.Id,
				Quantity: p.Quantity,
				Name: p.Name,
				Description: p.Description,
				Price: p.Price,
			})
		}
		newOrder.Products=products
		orders=append(orders, newOrder)
	}
	return orders,nil
}