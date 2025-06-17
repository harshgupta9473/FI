package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/harshgupta9473/fi/dto"
	"github.com/harshgupta9473/fi/logger"
	"go.uber.org/zap"
)

type UsersRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

type UsersRepoIntf interface {
	AddUser(ctx context.Context, user *dto.User) error
	GetUserByUsername(ctx context.Context, username string) (*dto.User, error)
}

func NewUsersRepository(db *sql.DB, logger *logger.Logger) UsersRepoIntf {
	return &UsersRepository{
		db:     db,
		logger: logger,
	}
}

func (u *UsersRepository) AddUser(ctx context.Context, user *dto.User) error {
	query := `
		INSERT INTO users (username, password)
		VALUES ($1, $2)
	`
	_, err := u.db.ExecContext(ctx, query, user.Username, user.Password)
	if err != nil {
		u.logger.Error(err, "failed to insert user",
			zap.String("username", user.Username),
		)
		return err
	}
	u.logger.Info("user added successfully", zap.String("username", user.Username))
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
			u.logger.Info("user not found", zap.String("username", username))
			return nil, nil
		}
		u.logger.Error(err, "failed to get user by username", zap.String("username", username))
		return nil, err
	}
	u.logger.Info("user retrieved successfully", zap.String("username", user.Username))
	return &user, nil
}
