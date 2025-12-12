package user

import (
	"errors"
	"log"
	"stu-go/common"
)

// Service defines user-related business operations.
type service struct {
	repo repository
}

func (s *service) Create(u User) (*User, error) {
	return s.repo.Create(u)
}

func (s *service) List() ([]User, error) {
	return s.repo.List()
}
func (s *service) GetByID(id int64) (*User, error) {
	return s.repo.GetByID(id)
}
func (s *service) DeleteByUserName(username string) error {
	err := s.repo.DeleteByUserName(username)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) Login(user *User) error {
	md5String := common.GetMd5String(user.Password)
	err := s.repo.GetUserByUserName(user)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if user.Password != md5String {
		return errors.New("密码错误请重新登录")
	}
	token, err := common.GenerateToken(user.Id, user.Username)
	if err != nil {
		return err
	}
	user.Token = token
	return nil
}
