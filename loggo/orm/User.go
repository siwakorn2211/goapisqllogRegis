package orm

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string
	Passworde string
	Fullname  string
	Avatar    string
}
