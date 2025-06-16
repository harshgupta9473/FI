package di

import (
	"database/sql"
	"fmt"
	"github.com/harshgupta9473/fi/configs"
	logger2 "github.com/harshgupta9473/fi/logger"
	"github.com/harshgupta9473/fi/repository"
	"github.com/harshgupta9473/fi/services"
)

type Container struct {
	DB *sql.DB

	Logger *logger2.Logger

	TableCreated bool

	UserRepository    repository.UsersRepoIntf
	ProductRepository repository.ProductsRepoIntf

	UserService    services.UserServiceIntf
	ProductService services.ProductServiceIntf
}

func NewContainer(env *configs.Config) (*Container, error) {
	db, err := initDB(env)
	if err != nil {
		return nil, err
	}

	container := &Container{
		DB: db,
	}

	logger, err := logger2.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("error in initializing the logger: %v", err)
	}
	container.Logger = logger

	err = container.CreateAllTables()
	if err != nil {
		return nil, err
	}
	container.TableCreated = true
	container.initRepositories()
	container.initServices()
	return container, nil
}

func (c *Container) initRepositories() {
	c.UserRepository = repository.NewUsersRepository(c.DB)
	c.ProductRepository = repository.NewProductsRepository(c.DB)
}

func (c *Container) initServices() {
	c.UserService = services.NewUserService(c.UserRepository)
	c.ProductService = services.NewProductService(c.ProductRepository)
}
