package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/graphql-go/handler"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountURL 	string `envconfig:"ACCOUNT_SERVICE_URL"`
	ProductURL 	string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL 	string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	var cfg AppConfig
	err:=envconfig.Process("",&cfg)
	if err!=nil{
		log.Fatal(err)
	}
	s,err:=NewGraphQLServer(cfg.AccountURL,cfg.ProductURL,cfg.OrderURL)
	if err!=nil{
		log.Fatal(err)
	}
	http.Handle("/graphql",handler.New(&handler.Config{
		Schema: s.ToExecutableTableSchema(),
		Pretty: true,
		GraphiQL: true,
	}))
	http.Handle("/playground",playground.Handler("eyepatch","/graphql"))
	log.Fatal(http.ListenAndServe(":8000",nil))
}