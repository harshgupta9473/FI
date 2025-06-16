package server

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/harshgupta9473/fi/configs"
	"github.com/harshgupta9473/fi/di"
	"github.com/harshgupta9473/fi/handlers"
	"github.com/harshgupta9473/fi/routes"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	log.Println("ENV loaded")
	environment, err := configs.LoadEnvironment()
	if err != nil {
		os.Exit(1)
	}
	container, err := di.NewContainer(environment)
	if err != nil {
		os.Exit(1)
	}
	handler := handlers.NewHandler(container.ProductService, container.UserService)
	router := mux.NewRouter()
	routes.SetupRoutes(router, handler)

	s := &http.Server{
		Addr:         ":3001",
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("Listening on port 3001")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal Closing the server:", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(ctx)
}
