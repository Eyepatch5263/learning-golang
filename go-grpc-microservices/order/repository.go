package order

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, o Order) error 
	GetOrdersForAccount(ctx context.Context, accountID string)([]Order,error)
}

type postgresRepository struct{
	db *sql.DB
}

func NewPostgresRepository(url string)(Repository,error){
	db,err:=sql.Open("postgres",url)
	if err!=nil{
		return nil,err
	}
	err=db.Ping()
	if err!=nil{
		return nil,err
	}
	return &postgresRepository{db:db},nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) PutOrder(ctx context.Context,o Order) error {

	//begins the transaction it can either fail/rollback or commit/finish
	tx,err:=r.db.BeginTx(ctx,nil)
	if err!=nil{
		return err
	}
	defer func (){
		// if error occurs rollback
		if err!=nil{
			tx.Rollback()
			return
		}
		// if no error occurs commit
		err=tx.Commit()
	}()
	// executes the command od inserting the values in the table via taking ctx,query and other args
	tx.ExecContext(
		ctx,
		"INSERT INTO orders(id,created_at,account_id,total_price) VALUES($1,$2,$3,$4)",
		o.ID,
		o.CreatedAt,
		o.AccountID,
		o.TotalPrice,
	)
	if err!=nil{
		return err
	}
	// now insert the values in the ordered_products with the following fields where order_id is pointing to id field in the orders table
	stmt,_:=tx.PrepareContext(ctx,pq.CopyIn("order_products","order_id","product_id","quantity"))
	for _,p:= range o.Products{
		_,err=stmt.ExecContext(ctx,o.ID,p.ID,p.Quantity)
		if err !=nil{
			return err
		}
	}
	_,err=stmt.ExecContext(ctx)
	if err!=nil{
		return err
	}
	stmt.Close()
	return nil

}

func (r *postgresRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order,error) {
	rows,err:=r.db.QueryContext(ctx,
	`
		SELECT
		    o.id, o.created_at,
		    o.account_id,
		    o.total_price::MONEY::NUMERIC::FLOAT8,
		    op.product_id,
		    op.quantity
		FROM
		    orders o
		    JOIN order_products op ON (o.id == op.order_id)
		WHERE
		    o.account_id = $1
		ORDER BY
		    o.id`,
	accountID,
	)
	if err !=nil{
		return nil,err
	}
	defer rows.Close()
	orders:=[]Order{}
	currentOrder:=&Order{}
	lastOrder:=&Order{}
	orderedProduct:=&OrderedProduct{}
	products:=[]OrderedProduct{}

	for rows.Next(){
		if err=rows.Scan(
			&currentOrder.ID,
			&currentOrder.CreatedAt,
			&currentOrder.AccountID,
			&currentOrder.TotalPrice,
			&orderedProduct.ID,
			&orderedProduct.Quantity,
		);err!=nil{
			return nil,err
		}
		if lastOrder.ID!="" && lastOrder.ID!=currentOrder.ID{
			newOrder:=Order{
				ID:lastOrder.ID,
				AccountID: lastOrder.AccountID,
				CreatedAt: lastOrder.CreatedAt,
				TotalPrice: lastOrder.TotalPrice,
				Products: lastOrder.Products,
			}
			orders=append(orders, newOrder)
			products=[]OrderedProduct{}
		}
		products=append(products, OrderedProduct{
			ID: orderedProduct.ID,
			Quantity: orderedProduct.Quantity,
		})
		lastOrder=currentOrder
	}
	if lastOrder!=nil{
		newOrder:=Order{
				ID:lastOrder.ID,
				AccountID: lastOrder.AccountID,
				CreatedAt: lastOrder.CreatedAt,
				TotalPrice: lastOrder.TotalPrice,
				Products: lastOrder.Products,
			}
			orders=append(orders, newOrder)
	}
	if err=rows.Err();err!=nil{
		return nil,err
	}
	return orders,nil
}