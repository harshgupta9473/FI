package di

import (
	"fmt"
)

func (c *Container) CreateAllTables() error {
	var err error
	err = c.createUserAuthTable()
	if err != nil {
		return err
	}
	err = c.CreateProductsTable()
	if err != nil {
		return err
	}
	return nil
}

func (c *Container) createUserAuthTable() error {
	query := `CREATE users table if not exists(
	username VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    )`
	_, err := c.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creaeting users table: %v", err)
	}
	return nil
}

func (c *Container) CreateProductsTable() error {
	query := `CREATE TABLE products IF NOT EXISTS product (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    sku TEXT UNIQUE NOT NULL,
    image_url TEXT,
    description TEXT,
    quantity INTEGER NOT NULL DEFAULT 0,
    price NUMERIC(10, 2) NOT NULL
   );`
	_, err := c.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating products table: %v", err)
	}
	return nil
}
