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
	err = c.createProductsTable()
	if err != nil {
		return err
	}
	return nil
}

func (c *Container) createUserAuthTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL
	)`
	_, err := c.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}
	return nil
}

func (c *Container) createProductsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		sku TEXT UNIQUE NOT NULL,
		image_url TEXT,
		description TEXT,
		quantity INTEGER NOT NULL DEFAULT 0,
		price NUMERIC(10, 2) NOT NULL
	)`
	_, err := c.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating products table: %v", err)
	}
	return nil
}
