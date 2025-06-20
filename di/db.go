package di

import (
	"database/sql"
	"fmt"
	"github.com/harshgupta9473/fi/configs"
	_ "github.com/lib/pq"
)

func initDB(env *configs.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", env.DBConnStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting with database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}
	return db, nil
}
