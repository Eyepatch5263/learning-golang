package main

import (
	"log"
	"time"
	"github.com/eyepatch5263/go-grpc-microservices/product"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main(){
	var cfg Config
	err:=envconfig.Process("",&cfg)
	if err!=nil{
		log.Fatal(err)
	}
	var r product.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error){
		r, err = product.NewElasticRepository(cfg.DatabaseURL)
		if err!=nil{
			log.Println("Auth error:",err)
			log.Fatal(err)
		}
		return
	})
	defer r.Close()
	log.Println("Listening on port 8080...")
	s:=product.NewService(r)
	log.Fatal(product.ListenGRPC(s,8080))
}