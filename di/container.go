package di

import (
	"database/sql"
	"fmt"
	"github.com/harshgupta9473/fi/configs"
	"github.com/harshgupta9473/fi/logger"
	"github.com/harshgupta9473/fi/repository"
	"github.com/harshgupta9473/fi/services"
)

type Container struct {
	DB *sql.DB

	Logger *logger.Logger

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
	container.Logger, err = logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("error in initializing the logger: %v", err)
	}

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
	c.UserRepository = repository.NewUsersRepository(c.DB, c.Logger)
	c.ProductRepository = repository.NewProductsRepository(c.DB, c.Logger)
}

func (c *Container) initServices() {
	c.UserService = services.NewUserService(c.UserRepository, c.Logger)
	c.ProductService = services.NewProductService(c.ProductRepository, c.Logger)
}
