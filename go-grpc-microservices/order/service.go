package order

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type OrderedProduct struct {
	ID string 
	Name  string    
	Description string
	Price float64
	Quantity	uint32
}

type Order struct {
	ID string 	`json:"id"`
	AccountID string 	`json:"account_id"`
	CreatedAt time.Time	 `json:"created_at"`
	TotalPrice float64	  `json:"total_price"`
	Products []OrderedProduct `json:"products"`
}

type OrderService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &OrderService{r}
}

func (s OrderService) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order,error) {
	order:=&Order{
		ID:ksuid.New().String(),
		AccountID:accountID,
		CreatedAt:time.Now().UTC(),
		Products: products,
	}
	order.TotalPrice=0.0
	for _,product:=range products{
		order.TotalPrice+=product.Price*float64(product.Quantity)
	}
	err:=s.repository.PutOrder(ctx,*order)
	if err!=nil{
		return nil,err
	}
	return order,nil
}

func (s OrderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	return s.repository.GetOrdersForAccount(ctx,accountID)
}
