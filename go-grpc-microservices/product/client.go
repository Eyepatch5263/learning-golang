package product

import (
	"context"

	"github.com/eyepatch5263/go-grpc-microservices/product/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn *grpc.ClientConn
	service pb.ProductServiceClient
}

func NewClient (url string)(*Client,error){
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil,err
	}
	c:=pb.NewProductServiceClient(conn)
	return &Client{conn,c},nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostProduct(ctx context.Context, name,description string, price float64) (*Product,error) {
	r,err:=c.service.PostProduct(ctx,&pb.PostProductRequest{
		Name: name,
		Description: description,
		Price: price,
	})
	if err!=nil{
		return nil,err
	}
	return &Product{
		ID: r.Product.Id,
		Name:  r.Product.Name,
		Price: r.Product.Price,
		Description: r.Product.Description,
	},nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product,error) {
	r,err:=c.service.GetProduct(ctx,&pb.GetProductRequest{
		Id: id,
	})
	if err!=nil{
		return nil,err
	}
	return &Product{
		ID: r.Product.Id,
		Name:  r.Product.Name,
		Price: r.Product.Price,
		Description: r.Product.Description,
	},nil

}

func (c *Client) GetProducts(ctx context.Context, query string, ids []string, skip uint64, take uint64) ([]Product,error) {
	r,err:=c.service.GetProducts(
		ctx,
		&pb.GetProductsRequest{
			Query: query,
			Ids: ids,
			Skip: skip,
			Take: take,
		},
	)
	if err!=nil{
		return nil,err
	}
	var products []Product
	for _,product:=range r.Products{
		products=append(products, Product{
			ID: product.Id,
			Name: product.Name,
			Price: product.Price,
			Description: product.Description,
		})
	}
	return products,nil
}