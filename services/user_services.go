package services

import (
	"context"
	"fmt"
	"github.com/harshgupta9473/fi/dto"
	"github.com/harshgupta9473/fi/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo repository.UsersRepoIntf
}

type UserServiceIntf interface {
	CreateUserAccount(ctx context.Context, user *dto.User) error
	LoginUser(ctx context.Context, user *dto.User) error
}

func NewUserService(userRepo repository.UsersRepoIntf) UserServiceIntf {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (u *UserService) CreateUserAccount(ctx context.Context, user *dto.User) error {
	//if user already exists
	userDB, err := u.UserRepo.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return fmt.Errorf("error in checking the username status i.e. if it exists or not")
	}
	if userDB != nil {
		return fmt.Errorf("username already exists")
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error in encrypting the password")
	}
	user.Password = string(encryptedPassword)

	err = u.UserRepo.AddUser(ctx, user)
	if err != nil {
		return fmt.Errorf("error creating user account")
	}

	return nil

}

func (u *UserService) LoginUser(ctx context.Context, user *dto.User) error {
	userfromDB, err := u.UserRepo.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return fmt.Errorf("error in getting the user details")
	}
	if userfromDB == nil {
		return fmt.Errorf("invalid user name of password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(userfromDB.Password), []byte(user.Password))
	if err != nil {
		return fmt.Errorf("invalid user name of password")
	}
	return nil
}
