package server

import (
	"github.com/harshgupta9473/fi/configs"
	"github.com/harshgupta9473/fi/di"
	"github.com/harshgupta9473/fi/handlers"
	"os"
)

func main() {
	environment, err := configs.LoadEnvironment()
	if err != nil {
		os.Exit(1)
	}
	container, err := di.NewContainer(environment)
	if err != nil {
		os.Exit(1)
	}
	handler := handlers.NewHandler(container.ProductService, container.UserService)

}
