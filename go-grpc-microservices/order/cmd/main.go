package main

import (
	"log"
	"time"

	"github.com/eyepatch5263/go-grpc-microservices/order"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
	AccountURL string 	`envconfig:"ACCOUNT_SERVICE_URL"`
	ProductURL string 	`envconfig:"PRODUCT_SERVICE_URL"`
}

func main( ){
	var cfg Config
	err:=envconfig.Process("",&cfg)
	if err!=nil{
		log.Fatal(err)
	}

	var r order.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error){
		r, err = order.NewPostgresRepository(cfg.DatabaseURL)
		if err!=nil{
			log.Println("Auth error:",err)
			log.Fatal(err)
		}
		return
	})
	defer r.Close()
	log.Println("Listening on port 8080...")
	s:=order.NewService(r)
	log.Fatal(order.ListenGRPC(s,cfg.AccountURL,cfg.ProductURL,8080))
}