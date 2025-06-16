package di

import (
	"database/sql"
	"github.com/harshgupta9473/fi/configs"
	"github.com/harshgupta9473/fi/repository"
	"github.com/harshgupta9473/fi/services"
)

type Container struct {
	DB *sql.DB

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
