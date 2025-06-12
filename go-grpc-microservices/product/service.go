package product

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostProduct(ctx context.Context, name string, description string, price float64) (*Product,error)
	GetProduct(ctx context.Context,id string)	(*Product,error)
	GetProducts(ctx context.Context, skip uint64, take uint64)	([]Product,error)
	GetProductByIDs(ctx context.Context, ids []string)	([]Product,error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product,error)
}

type Product struct {
	ID string 	`json:"id"`
	Name string 	`json:"name"`
	Price float64 	`json:"price"`
	Description string 	`json:"description"`
}

type ProductService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &ProductService{repo: r}
}

func (s *ProductService) PostProduct(ctx context.Context, name string, description string, price float64) (*Product,error) {
	p:=&Product{
		ID:  ksuid.New().String(),
		Name: name,
		Price: price,
		Description: description,
	}
	if err:=s.repo.PutProduct(ctx,*p);err!=nil{
		return nil,err
	}
	return p,nil
}

func (s *ProductService) GetProduct(ctx context.Context,id string) (*Product,error) {
	return s.repo.GetProductByID(ctx,id)
}

func (s *ProductService) GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product,error) {
	if take>100 || (skip==0 && take==0){
		take=100
	}
	return s.repo.ListProducts(ctx,skip,take)
}

func (s *ProductService) GetProductByIDs(ctx context.Context, ids []string) ([]Product,error) {
	return s.repo.ListProductsWithIDs(ctx,ids)
}

func (s *ProductService) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product,error) {
	if take>100 || (skip==0 && take==0){
		take=100
	}
	return s.repo.SearchProducts(ctx,query,skip,take)
}