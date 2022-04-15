package dao

import (
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	o *gorm.DB
)

func init() {
	if db, err := gorm.Open(sqlite.Open("models/models.db")); err != nil {
		panic(err)
	} else {
		o = db
	}
	err := o.AutoMigrate(
		&UserModel{},
		&GroupModel{},
		&FriendModel{},
		&FriendRemarkModel{},
		&MemberModel{},
	)
	if err != nil {
		panic(err)
	}
}
