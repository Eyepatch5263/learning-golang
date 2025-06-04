package main

import (
	"fmt"
	"log"
	"net/http"
   	_ "github.com/lib/pq"
	"github.com/eyepatch5263/go-postgress/routes"
	"github.com/gorilla/mux"
)

func main(){
	// Initialize the router
	r:=mux.NewRouter()	
	// This will register all the routes defined in the routes package
	routes.RegisterRoutes(r)
	fmt.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000",r))
}