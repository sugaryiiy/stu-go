package user

import "time"

type User struct {
	ID        int64     `xorm:"pk autoincr 'id'" json:"id"`
	Username  string    `xorm:"varchar(50) unique notnull 'username'" json:"username"`
	Password  string    `xorm:"varchar(255) notnull 'password'" json:"-"`
	Salt      string    `xorm:"varchar(50) 'salt'" json:"-"`
	Email     string    `xorm:"varchar(100) 'email'" json:"email,omitempty"`
	Phone     string    `xorm:"varchar(20) 'phone'" json:"phone,omitempty"`
	Status    int       `xorm:"tinyint default 1 'status'" json:"status"`
	CreatedAt time.Time `xorm:"created 'created_at'" json:"createdAt"`
	UpdatedAt time.Time `xorm:"updated 'updated_at'" json:"updatedAt"`
	List      []User    `xorm:"-" json:"list,omitempty"`
}
