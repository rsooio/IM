package services

import (
	"IM/dao"
	"errors"
	"gorm.io/gorm"
	"sync"
)

type (
	Group struct {
		ID   uint
		Name string
	}
	Member struct {
		ID       uint
		Name     string
		Nickname string
	}
	MemberList []*Member
	GroupList  []*Group
)

func (m *Group) Query() error {
	group := dao.GroupModel{ID: m.ID}
	err := group.Query()
	if err != nil {
		return err
	}
	m.Name = group.Name
	return nil
}

func (m *Group) Auth(userID uint) bool {
	association := dao.MemberModel{UserID: userID, GroupID: m.ID}
	err := association.Query()
	if err != nil {
		return false
	}
	if association.Flag&(dao.FlagMemberAdmin|dao.FlagMemberOwner) != 0 {
		return true
	}
	return false
}

func (m *Group) Members() (members MemberList, errs []error) {
	wg := sync.WaitGroup{}
	var associations dao.MemberList
	err := associations.Find(
		dao.Where(dao.MemberModel{GroupID: m.ID}),
		dao.IncludeFlag(dao.FlagMemberEstablishment),
	)
	if err != nil {
		errs = append(errs, err)
		return
	}
	for _, association := range associations {
		wg.Add(1)
		go func(association *dao.MemberModel) {
			user := dao.UserModel{ID: association.UserID}
			err := user.Query()
			if err != nil {
				errs = append(errs, err)
			}
			members = append(members, &Member{
				ID:       user.ID,
				Name:     user.Nickname,
				Nickname: association.GroupNickname,
			})
			wg.Done()
		}(association)
	}
	wg.Wait()
	return
}

func (m *Group) Invite(memberID uint) error {
	association := dao.MemberModel{UserID: memberID, GroupID: m.ID}
	err := association.Query()
	if err == gorm.ErrRecordNotFound {
		association.Flag = dao.FlagMemberInvite
		return association.Create()
	} else if err != nil {
		return err
	}
	if association.Flag&dao.FlagMemberInvite != 0 {
		return errors.New("already invited")
	}
	if association.Flag&dao.FlagMemberEstablishment != 0 {
		return errors.New("already joined")
	}
	if association.Flag&dao.FlagMemberRequest != 0 {
		association.Flag |= dao.FlagMemberEstablishment
		association.Flag &^= dao.FlagMemberRequest
	} else {
		association.Flag |= dao.FlagMemberInvite
	}
	return association.ForceUpdate()
}

func (m *Group) Create(ownerID uint) error {
	return dao.Transaction(func(tx *gorm.DB) error {
		group := dao.GroupModel{Name: m.Name}
		err := tx.Create(&group).Error
		if err != nil {
			return err
		}
		err = tx.Create(&dao.MemberModel{
			UserID:  ownerID,
			GroupID: group.ID,
			Flag:    dao.FlagMemberOwner | dao.FlagMemberEstablishment,
		}).Error
		if err != nil {
			return err
		}
		m.ID = group.ID
		return nil
	})
}
