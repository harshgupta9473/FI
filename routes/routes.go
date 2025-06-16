package routes

import (
	"github.com/gorilla/mux"
	"github.com/harshgupta9473/fi/handlers"
	"github.com/harshgupta9473/fi/middleware"
	"net/http"
)

func SetupRoutes(router *mux.Router, handlers *handlers.Handlers) {
	router.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	router.HandleFunc("/login", handlers.LoginUser).Methods("POST")

	router.Handle("/products", middleware.AuthMiddleware(http.HandlerFunc(handlers.AddProduct))).Methods("POST")
	router.Handle("/products", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetAllProducts))).Methods("GET")
	router.Handle("/products/{id}/quantity", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateProductQuantitty))).Methods("PUT")
}
