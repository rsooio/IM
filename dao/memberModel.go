package dao

import (
	"time"
)

const (
	FlagMemberRequest = 1 << iota
	FlagMemberInvite
	FlagMemberEstablishment
	FlagMemberAdmin
	FlagMemberOwner
)

type (
	MemberModel struct {
		ID            uint      `gorm:"primarykey"`
		UserID        uint      `gorm:""`
		GroupID       uint      `gorm:""`
		GroupNickname string    `gorm:""`
		Flag          uint64    `gorm:""`
		Version       uint      `gorm:""`
		CreatedAt     time.Time `gorm:""`
		UpdatedAt     time.Time `gorm:""`
	}
	MemberList []*MemberModel
)

func (m *MemberModel) Query() error {
	return o.Where(m).First(m).Error
}

func (m *MemberModel) Create() error {
	return o.Create(m).Error
}

func (m *MemberModel) Update() error {
	return o.Model(m).Updates(m).Error
}

func (m *MemberModel) ForceUpdate() error {
	return o.Model(m).Save(m).Error
}

func (m *MemberModel) Delete() error {
	return o.Delete(m).Error
}

func (lst *MemberList) Find(handlers ...Handler) error {
	db := o
	for _, f := range handlers {
		db = f(db)
	}
	return db.Find(lst).Error
}
