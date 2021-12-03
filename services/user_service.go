package services

import (
	"errors"
	datamodels "shopping_system/dataModels"
	"shopping_system/repositories"

	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOK bool)
	AddUser(user *datamodels.User) (userId int64, err error)
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func NewService(repository repositories.IUserRepository) IUserService {
	return &UserService{repository}
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOK bool) {
	user, err := u.UserRepository.Select(userName)
	if err != nil {
		return
	}
	isOK, _ = ValidatePassword(pwd, user.HashPassword)
	if !isOK {
		return &datamodels.User{}, false
	}
	return
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	pwdByte, errPwd := GeneratePassword(user.HashPassword)
	if errPwd != nil {
		return userId, errPwd
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func ValidatePassword(userPassword string, hashed string) (isOK bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码对比错误！")
	}
	return true, nil
}
