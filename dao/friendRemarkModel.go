package dao

import "time"

type (
	FriendRemarkModel struct {
		ID        uint      `gorm:"primarykey"`
		UserId    uint      `gorm:""`
		FriendID  uint      `gorm:""`
		Alias     string    `gorm:""`
		State     uint64    `gorm:""`
		Version   uint      `gorm:""`
		CreatedAt time.Time `gorm:""`
		UpdatedAt time.Time `gorm:""`
	}
	FriendRemarkList []*FriendRemarkModel
)

func (m *FriendRemarkModel) Query() error {
	return o.Where(m).First(m).Error
}

func (m *FriendRemarkModel) Create() error {
	return o.Create(m).Error
}

func (m *FriendRemarkModel) Update() error {
	return o.Model(m).Updates(m).Error
}

func (m *FriendRemarkModel) ForceUpdate() error {
	return o.Model(m).Save(m).Error
}

func (m *FriendRemarkModel) Delete() error {
	return o.Delete(m).Error
}

func (lst *FriendRemarkList) Find(handlers ...Handler) error {
	for _, f := range handlers {
		o = f(o)
	}
	return o.Find(lst).Error
}
