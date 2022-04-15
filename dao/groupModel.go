package dao

import (
	"gorm.io/gorm"
	"time"
)

type (
	GroupModel struct {
		ID        uint           `gorm:"primarykey"`
		Name      string         `gorm:""`
		Flag      uint64         `gorm:""`
		Version   uint           `gorm:""`
		CreatedAt time.Time      `gorm:""`
		UpdatedAt time.Time      `gorm:""`
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}
	GroupList []*GroupModel
)

func (m *GroupModel) Query() error {
	return o.Where(m).First(m).Error
}

func (m *GroupModel) Create() error {
	return o.Create(m).Error
}

func (m *GroupModel) Update() error {
	return o.Model(m).Updates(m).Error
}

func (m *GroupModel) ForceUpdate() error {
	return o.Model(m).Save(m).Error
}

func (m *GroupModel) Delete() error {
	return o.Delete(m).Error
}

func (lst *GroupList) Find(handlers ...Handler) error {
	db := o
	for _, f := range handlers {
		db = f(db)
	}
	return db.Find(lst).Error
}
