package routes

import (
	"github.com/eyepatch5263/go-postgress/middleware"
	"github.com/gorilla/mux"
)

var RegisterRoutes = func(router *mux.Router) {
	router.HandleFunc("/api/stocks/", middleware.GetAllStock).Methods("GET","OPTIONS")
	router.HandleFunc("/api/stocks/{stockId}", middleware.GetStockById).Methods("GET","OPTIONS")
	router.HandleFunc("/api/stocks/", middleware.CreateStock).Methods("POST","OPTIONS")
	router.HandleFunc("/api/stocks/{stockId}", middleware.UpdateStockById).Methods("PUT","OPTIONS")
	router.HandleFunc("/api/stocks/{stockId}", middleware.DeleteStockById).Methods("DELETE","OPTIONS")
	
}