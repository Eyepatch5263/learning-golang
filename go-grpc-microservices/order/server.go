package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/eyepatch5263/go-grpc-microservices/account"
	"github.com/eyepatch5263/go-grpc-microservices/order/pb"
	"github.com/eyepatch5263/go-grpc-microservices/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)


type grpcServer struct{
	service Service
	accountClient *account.Client
	productClient *product.Client
	pb.UnimplementedOrderServiceServer
}


func ListenGRPC( s Service,accountURL,productURL string,port int) error {
	accountClient,err:=account.NewClient(accountURL)
	if err!=nil{
		return err
	}
	productClient,err:=product.NewClient(productURL)
	if err!=nil{
		accountClient.Close()
		return err
	}
	lis,err:=net.Listen("tcp",fmt.Sprintf(":%d",port))
	if err!=nil{
		accountClient.Close()
		productClient.Close()
		return err
	}
	server:=grpc.NewServer()
	pb.RegisterOrderServiceServer(server,&grpcServer{
		s,
		accountClient,
		productClient,
		pb.UnimplementedOrderServiceServer{},
	})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest)(*pb.PostOrderResponse,error){
	_,err:=s.accountClient.GetAccount(ctx,r.AccountId)
	if err!=nil{
		log.Println("Error getting account:",err)
		return nil,errors.New("account not found")
	}
	productIDs:=[]string{}
	orderedProducts,err:=s.productClient.GetProducts(ctx,"",productIDs,0,0)
	if err!=nil{
		log.Println("Error getting products :",err)
		return nil,errors.New("product not found")
	}
	products:=[]OrderedProduct{}
	for _,p:=range orderedProducts{
		product:=OrderedProduct{
			ID:p.ID,
			Name:p.Name,
			Price:p.Price,
			Quantity:0,
			Description:p.Description,
		}
		for _,rp:=range r.Products{
			if rp.ProductId==p.ID{
				product.Quantity=rp.Quantity
				break
			}
		}
		if product.Quantity!=0{
			products=append(products, product)
		}
	}
	order,err:=s.service.PostOrder(ctx,r.AccountId,products)
	if err!=nil{
		log.Println("Error posting order:",err)
		return nil,errors.New("could not post order")
	}
	orderProto:=&pb.Order{
		Id:order.ID,
		AccountId: order.AccountID,
		TotalPrice: order.TotalPrice,
		Products: []*pb.Order_OrderProduct{},
	}
	orderProto.CreatedAt,_=order.CreatedAt.MarshalBinary()
	for _,p:=range order.Products{
		orderProto.Products=append(orderProto.Products, &pb.Order_OrderProduct{
			Id:p.ID,
			Name:p.Name,
			Description: p.Description,
			Price: p.Price,
			Quantity: p.Quantity,
		})
	}
	return &pb.PostOrderResponse{
		Order: orderProto,
	},nil
}

func (s *grpcServer) GetOrdersForAccount(ctx context.Context,r *pb.GetOrdersForAccountRequest)(*pb.GetOrdersForAccountResponse,error){
	accountOrders,err:=s.service.GetOrdersForAccount(ctx,r.AccountId)
	if err!=nil{
		log.Println(err)
		return nil,err
	}
	productIDMap:=map[string]bool{}
	for _,o:=range accountOrders{
		for _,p:=range o.Products{
			productIDMap[p.ID]=true
		}
	}
	productIDs:=[]string{}
	for id:=range productIDMap{
		productIDs=append(productIDs, id)
	}
	products,err:=s.productClient.GetProducts(ctx,"",productIDs,0,0)
	if err!=nil{
		log.Println("error getting account products:",err)
		return nil,err
	}
	orders:=[]*pb.Order{}
	for _,o:= range accountOrders{
		op:=&pb.Order{
			AccountId: o.AccountID,
			Id:o.ID,
			TotalPrice: o.TotalPrice,
			Products: []*pb.Order_OrderProduct{},
		}
		op.CreatedAt,_=o.CreatedAt.MarshalBinary()
		for _,product:=range o.Products{
			for _,p:=range products{
				if p.ID==product.ID{
					product.Name=p.Name
					product.Description=p.Description
					product.Price=p.Price
					break
				}
			}
			op.Products = append(op.Products, &pb.Order_OrderProduct{
				Id:product.ID,
				Name:product.Name,
				Description: product.Description,
				Price: product.Price,
				Quantity: product.Quantity,
			})
		}
		orders=append(orders, op)
	}
	return &pb.GetOrdersForAccountResponse{
		Orders: orders,
	},nil
}