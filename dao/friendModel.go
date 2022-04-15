package dao

import (
	"time"
)

const (
	FlagFriendGeneralEstablishment uint64 = 1 << iota
)

const (
	FlagFriendDedicateRequest uint64 = 1 << iota
	FlagFriendDedicateDelete
	FlagFriendDedicateBlock
)

const (
	FlagUserAStep = 32
	FlagUserBStep = 48
)

type (
	FriendModel struct {
		ID        uint      `gorm:"primarykey"`
		UserAID   uint      `gorm:"index"`
		UserBID   uint      `gorm:"index"`
		Flag      uint64    `gorm:""`
		Version   uint      `gorm:""`
		CreatedAt time.Time `gorm:""`
		UpdatedAt time.Time `gorm:""`
	}
	FriendList []*FriendModel
)

func (m *FriendModel) Query() error {
	return o.Where(m).First(m).Error
}

func (m *FriendModel) Create() error {
	return o.Create(m).Error
}

func (m *FriendModel) Update() error {
	return o.Model(m).Updates(m).Error
}

func (m *FriendModel) ForceUpdate() error {
	return o.Model(m).Save(m).Error
}

func (m *FriendModel) Delete() error {
	return o.Delete(m).Error
}

func (lst *FriendList) Find(handlers ...Handler) error {
	db := o
	for _, f := range handlers {
		db = f(db)
	}
	return db.Find(lst).Error
}
