package services

import (
	"errors"
	"github.com/l1nkkk/shopSystem/demo/product/datamodels"
	"github.com/l1nkkk/shopSystem/demo/product/repositories"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	// IsPwdSuccess 检查账号密码是否匹配
	IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool)

	// AddUser 添加用户
	AddUser(user *datamodels.User) (userId int64, err error)
}

func NewService(repository repositories.IUserRepository) IUserService {
	return &UserService{repository}
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamodels.User, isOk bool) {

	// 1. 通过userName 获取 user
	user, err := u.UserRepository.Select(userName)

	if err != nil {
		return
	}
	// 2. 密码比对
	isOk, _ = ValidatePassword(pwd, user.HashPassword)

	if !isOk {
		return &datamodels.User{}, false
	}

	return
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	// 1. 对密码进行hash
	pwdByte, errPwd := GeneratePassword(user.HashPassword)
	if errPwd != nil {
		return userId, errPwd
	}
	user.HashPassword = string(pwdByte)

	// 2. 插入数据库
	return u.UserRepository.Insert(user)
}

func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func ValidatePassword(userPassword string, hashed string) (isOK bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码比对错误！")
	}
	return true, nil

}
