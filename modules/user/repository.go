package user

import (
	"errors"
	"log"

	"github.com/go-xorm/xorm"
)

type Repository interface {
	Create(p User) (*User, error)
	GetByID(id int64) (*User, error)
	List() ([]User, error)
}
type repository struct {
	db *xorm.Engine
}

func (r *repository) Create(p User) (*User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) GetByID(id int64) (*User, error) {
	user := &User{}
	has, err := r.db.SQL("select * from user where id=?", id).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *repository) List() ([]User, error) {
	user := &User{}
	err := r.db.SQL("select * from user").Find(&user.List)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return user.List, nil
}

func NewRepository(db *xorm.Engine) Repository {
	return &repository{db}
}
