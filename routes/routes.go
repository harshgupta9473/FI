package routes

import (
	"github.com/gorilla/mux"
	"github.com/harshgupta9473/fi/handlers"
)

func SetupRoutes(router *mux.Router, handlers *handlers.Handlers) {
	router.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	router.HandleFunc("/login", handlers.LoginUser).Methods("POST")

	router.HandleFunc("/products", handlers.AddProduct).Methods("POST")
	router.HandleFunc("/products", handlers.GetAllProducts).Methods("GET")
	router.HandleFunc("/products/{id}/quantity", handlers.UpdateProductQuantitty).Methods("PUT")
}
