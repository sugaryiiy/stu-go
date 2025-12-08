package user

import "time"

type User struct {
	ID        int64     `json:"id" xorm:"id"`
	Name      string    `json:"name" xorm:"name"`
	Email     string    `json:"email" xorm:"email"`
	CreatedAt time.Time `json:"created_at" xorm:"created_at"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated_at"`
	List      []User    `json:"list" xorm:"-"`
}
