package dao

import (
	"gorm.io/gorm"
	"time"
)

type (
	UserModel struct {
		ID        uint           `gorm:"primarykey"`
		Username  string         `gorm:"unique;index"`
		Nickname  string         `gorm:""`
		Hash      string         `gorm:""`
		Version   uint           `gorm:""`
		CreatedAt time.Time      `gorm:""`
		UpdatedAt time.Time      `gorm:""`
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}
	UserList []*UserModel
)

func (m *UserModel) BeforeUpdate(tx *gorm.DB) error {
	m.Version++
	return nil
}

func (m *UserModel) Query() error {
	return o.Where(m).First(m).Error
}

func (m *UserModel) Create() error {
	return o.Create(m).Error
}

func (m *UserModel) Update() error {
	return o.Model(m).Updates(m).Error
}

func (m *UserModel) ForceUpdate() error {
	return o.Model(m).Save(m).Error
}

func (m *UserModel) Delete() error {
	return o.Delete(m).Error
}

func (lst *UserList) Find(handlers ...Handler) error {
	db := o
	for _, f := range handlers {
		db = f(db)
	}
	return db.Find(lst).Error
}
