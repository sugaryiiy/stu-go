package utils

import "github.com/go-xorm/xorm"

type repository struct {
	db *xorm.Engine
}
