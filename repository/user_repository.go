package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/harshgupta9473/fi/dto"
)

type UsersRepository struct {
	db *sql.DB
}

type UsersRepoIntf interface {
	AddUser(ctx context.Context, user *dto.User) error
	GetUserByUsername(ctx context.Context, username string) (*dto.User, error)
}

func NewUsersRepository(db *sql.DB) UsersRepoIntf {
	return &UsersRepository{
		db: db,
	}
}

func (u *UsersRepository) AddUser(ctx context.Context, user *dto.User) error {
	query := `
		INSERT INTO users (username, password)
		VALUES ($1, $2)
	`
	_, err := u.db.ExecContext(ctx, query, user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u *UsersRepository) GetUserByUsername(ctx context.Context, username string) (*dto.User, error) {
	query := `
		SELECT id, username, password
		FROM users
		WHERE username = $1
	`

	var user dto.User
	err := u.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Println(err)
		return nil, err
	}

	return &user, nil
}
