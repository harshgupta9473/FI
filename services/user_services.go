package services

import (
	"context"
	"fmt"
	"github.com/harshgupta9473/fi/dto"
	"github.com/harshgupta9473/fi/logger"
	"github.com/harshgupta9473/fi/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo repository.UsersRepoIntf
	logger   *logger.Logger
}

type UserServiceIntf interface {
	CreateUserAccount(ctx context.Context, user *dto.User) error
	LoginUser(ctx context.Context, user *dto.User) error
}

func NewUserService(userRepo repository.UsersRepoIntf, logger *logger.Logger) UserServiceIntf {
	return &UserService{
		UserRepo: userRepo,
		logger:   logger,
	}
}

func (u *UserService) CreateUserAccount(ctx context.Context, user *dto.User) error {
	//if user already exists
	userDB, err := u.UserRepo.GetUserByUsername(ctx, user.Username)
	if err != nil {
		u.logger.Error(err, "failed to check if username exists", zap.String("username", user.Username))

		return fmt.Errorf("failed to check if username exists")
	}
	if userDB != nil {
		u.logger.Info("username already exists", zap.String("username", user.Username))
		return fmt.Errorf("username already exists")
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error(err, "failed to encrypt password", zap.String("username", user.Username))
		return fmt.Errorf("error in encrypting the password")
	}
	user.Password = string(encryptedPassword)

	err = u.UserRepo.AddUser(ctx, user)
	if err != nil {
		u.logger.Error(err, "failed to add user", zap.String("username", user.Username))
		return fmt.Errorf("error creating user account")
	}
	u.logger.Info("user account created successfully", zap.String("username", user.Username))
	return nil

}

func (u *UserService) LoginUser(ctx context.Context, user *dto.User) error {
	userfromDB, err := u.UserRepo.GetUserByUsername(ctx, user.Username)
	if err != nil {
		u.logger.Error(err, "failed to fetch user during login", zap.String("username", user.Username))
		return fmt.Errorf("error in getting the user details")
	}
	if userfromDB == nil {
		u.logger.Info("login attempt: user does not exists", zap.String("username", user.Username))
		return fmt.Errorf("invalid user name of password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(userfromDB.Password), []byte(user.Password))
	if err != nil {
		u.logger.Info("invalid login credentials", zap.String("username", user.Username))
		return fmt.Errorf("invalid user name of password")
	}
	u.logger.Info("user logged in successfully", zap.String("username", user.Username))
	return nil
}
