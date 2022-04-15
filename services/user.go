package services

import (
	"IM/dao"
	"IM/utils"
	"errors"
	"gorm.io/gorm"
	"sync"
)

type (
	UserRegistrar struct {
		Username string
		Nickname string
		Password string
	}
	User struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
	Friend struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Alias string `json:"alias"`
	}
	FriendList []*Friend
	UserList   []*User
)

func orderID(id1, id2 uint) (uint, uint, uint, uint) {
	if id1 > id2 {
		return id2, id1, dao.FlagUserBStep, dao.FlagUserAStep
	}
	return id1, id2, dao.FlagUserAStep, dao.FlagUserBStep
}

func (m *User) LoginByPwd(username string, pwd string) (err error) {
	u := dao.UserModel{
		ID:       m.ID,
		Username: username,
	}
	err = u.Query()
	if err != nil {
		return err
	}
	ok := utils.Verify(u.Hash, pwd)
	if !ok {
		return errors.New("password error")
	}
	m.Name = u.Nickname
	return nil
}

func (m *User) LoginByToken(token string) (err error) {
	uid, flag, err := utils.JWTVerify(token)
	if err != nil {
		return err
	}
	if flag != utils.JWTFlagAuthToken {
		return errors.New("token error")
	}
	u := dao.UserModel{ID: uid}
	err = u.Query()
	if err != nil {
		return err
	}
	m.ID = u.ID
	m.Name = u.Nickname
	return nil
}

func (registrar UserRegistrar) Register() (user User, err error) {
	hash, err := utils.Encrypt(registrar.Password)
	if err != nil {
		return User{}, err
	}
	u := dao.UserModel{
		Username: registrar.Username,
		Nickname: registrar.Nickname,
		Hash:     hash,
	}
	err = u.Create()
	if err != nil {
		return User{}, err
	}
	return User{
		ID:   u.ID,
		Name: u.Nickname,
	}, nil
}

func (m *User) Friends() (friends FriendList, errs []error) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		var associations dao.FriendList
		err := associations.Find(
			dao.Where(&dao.FriendModel{UserAID: m.ID}),
			dao.IncludeFlag(dao.FlagFriendGeneralEstablishment),
		)
		if err != nil {
			errs = append(errs, err)
			wg.Done()
			return
		}
		for _, association := range associations {
			wg.Add(1)
			go func(association *dao.FriendModel) {
				user := dao.UserModel{ID: association.UserBID}
				err := user.Query()
				if err != nil {
					errs = append(errs, err)
					wg.Done()
					return
				}
				friend := Friend{
					ID:   user.ID,
					Name: user.Nickname,
				}
				remark := dao.FriendRemarkModel{
					UserId:   m.ID,
					FriendID: user.ID,
				}
				err = remark.Query()
				if err != nil {
					errs = append(errs, err)
				} else {
					friend.Alias = remark.Alias
				}
				friends = append(friends, &friend)
				wg.Done()
			}(association)
		}
		wg.Done()
	}()
	go func() {
		var associations dao.FriendList
		err := associations.Find(
			dao.Where(&dao.FriendModel{UserBID: m.ID}),
			dao.IncludeFlag(dao.FlagFriendGeneralEstablishment),
		)
		if err != nil {
			errs = append(errs, err)
			wg.Done()
			return
		}
		for _, association := range associations {
			wg.Add(1)
			go func(association *dao.FriendModel) {
				user := dao.UserModel{ID: association.UserAID}
				err := user.Query()
				if err != nil {
					errs = append(errs, err)
					wg.Done()
					return
				}
				friend := Friend{
					ID:   user.ID,
					Name: user.Nickname,
				}
				remark := dao.FriendRemarkModel{
					UserId:   m.ID,
					FriendID: user.ID,
				}
				err = remark.Query()
				if err != nil {
					errs = append(errs, err)
				} else {
					friend.Alias = remark.Alias
				}
				friends = append(friends, &friend)
				wg.Done()
			}(association)
		}
		wg.Done()
	}()
	wg.Wait()
	return friends, errs
}

func (m *User) Followers() (followers UserList, errs []*error) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		var associations dao.FriendList
		err := associations.Find(
			dao.Where(&dao.FriendModel{UserAID: m.ID}),
			dao.ExcludeFlag(dao.FlagFriendGeneralEstablishment),
			dao.IncludeFlag(dao.FlagMemberRequest<<dao.FlagUserBStep),
		)
		if err != nil {
			errs = append(errs, &err)
			wg.Done()
			return
		}
		for _, association := range associations {
			wg.Add(1)
			go func(association *dao.FriendModel) {
				user := dao.UserModel{ID: association.UserBID}
				err := user.Query()
				if err != nil {
					errs = append(errs, &err)
					wg.Done()
					return
				}
				followers = append(followers, &User{
					ID:   user.ID,
					Name: user.Nickname,
				})
				wg.Done()
			}(association)
		}
		wg.Done()
	}()
	go func() {
		var associations dao.FriendList
		err := associations.Find(
			dao.Where(&dao.FriendModel{UserBID: m.ID}),
			dao.ExcludeFlag(dao.FlagFriendGeneralEstablishment),
			dao.IncludeFlag(dao.FlagMemberRequest<<dao.FlagUserAStep),
		)
		if err != nil {
			errs = append(errs, &err)
			wg.Done()
			return
		}
		for _, association := range associations {
			wg.Add(1)
			go func(association *dao.FriendModel) {
				user := dao.UserModel{ID: association.UserAID}
				err := user.Query()
				if err != nil {
					errs = append(errs, &err)
					wg.Done()
					return
				}
				followers = append(followers, &User{
					ID:   user.ID,
					Name: user.Nickname,
				})
				wg.Done()
			}(association)
		}
		wg.Done()
	}()
	wg.Wait()
	return followers, errs
}

func (m *User) Followings() (followings UserList, errs []*error) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		var associations dao.FriendList
		err := associations.Find(
			dao.Where(&dao.FriendModel{UserAID: m.ID}),
			dao.ExcludeFlag(dao.FlagFriendGeneralEstablishment),
			dao.IncludeFlag(dao.FlagMemberRequest<<dao.FlagUserAStep),
		)
		if err != nil {
			errs = append(errs, &err)
			wg.Done()
			return
		}
		for _, association := range associations {
			wg.Add(1)
			go func(association *dao.FriendModel) {
				user := dao.UserModel{ID: association.UserBID}
				err := user.Query()
				if err != nil {
					errs = append(errs, &err)
					wg.Done()
					return
				}
				followings = append(followings, &User{
					ID:   user.ID,
					Name: user.Nickname,
				})
				wg.Done()
			}(association)
		}
		wg.Done()
	}()
	go func() {
		var associations dao.FriendList
		err := associations.Find(
			dao.Where(&dao.FriendModel{UserBID: m.ID}),
			dao.ExcludeFlag(dao.FlagFriendGeneralEstablishment),
			dao.IncludeFlag(dao.FlagMemberRequest<<dao.FlagUserBStep),
		)
		if err != nil {
			errs = append(errs, &err)
			wg.Done()
			return
		}
		for _, association := range associations {
			wg.Add(1)
			go func(association *dao.FriendModel) {
				user := dao.UserModel{ID: association.UserAID}
				err := user.Query()
				if err != nil {
					errs = append(errs, &err)
					wg.Done()
					return
				}
				followings = append(followings, &User{
					ID:   user.ID,
					Name: user.Nickname,
				})
				wg.Done()
			}(association)
		}
		wg.Done()
	}()
	wg.Wait()
	return followings, errs
}

func (m *User) Follow(friendID uint) error {
	userAID, userBID, userStep, friendStep := orderID(m.ID, friendID)
	association := dao.FriendModel{UserAID: userAID, UserBID: userBID}
	err := association.Query()
	if err == gorm.ErrRecordNotFound {
		association.Flag = dao.FlagFriendDedicateRequest << userStep
		return association.Create()
	} else if err != nil {
		return err
	}
	if association.Flag&dao.FlagFriendGeneralEstablishment != 0 {
		return errors.New("already friend")
	}
	if association.Flag&(dao.FlagFriendDedicateRequest<<userStep) != 0 {
		return errors.New("already followed")
	}
	if association.Flag&(dao.FlagFriendDedicateRequest<<friendStep) != 0 {
		association.Flag |= dao.FlagFriendGeneralEstablishment
		association.Flag &^= dao.FlagMemberRequest << friendStep
	} else {
		association.Flag |= dao.FlagFriendDedicateRequest << userStep
	}
	return association.ForceUpdate()
}

func (m *User) RejectFollow(friendID uint) error {
	userAID, userBID, _, friendStep := orderID(m.ID, friendID)
	association := dao.FriendModel{UserAID: userAID, UserBID: userBID}
	err := association.Query()
	if err != nil {
		return err
	}
	if association.Flag&(dao.FlagFriendDedicateRequest<<friendStep) == 0 {
		return errors.New("not followed")
	}
	association.Flag &^= dao.FlagFriendDedicateRequest << friendStep
	return association.ForceUpdate()
}

func (m *User) CancelFollow(friendID uint) error {
	userAID, userBID, userStep, _ := orderID(m.ID, friendID)
	association := dao.FriendModel{UserAID: userAID, UserBID: userBID}
	err := association.Query()
	if err != nil {
		return err
	}
	if association.Flag&(dao.FlagFriendDedicateRequest<<userStep) == 0 {
		return errors.New("not followed")
	}
	association.Flag &^= dao.FlagFriendDedicateRequest << userStep
	return association.ForceUpdate()
}

func (m *User) RemoveFriend(friendID uint) error {
	userAID, userBID, _, _ := orderID(m.ID, friendID)
	association := dao.FriendModel{UserAID: userAID, UserBID: userBID}
	err := association.Query()
	if err != nil {
		return err
	}
	if association.Flag&dao.FlagFriendGeneralEstablishment == 0 {
		return errors.New("not friends")
	}
	association.Flag &^= dao.FlagFriendGeneralEstablishment
	return association.ForceUpdate()
}

func (m *User) Groups() (groups GroupList, errs []error) {
	wg := sync.WaitGroup{}
	var associations dao.MemberList
	err := associations.Find(
		dao.Where(dao.MemberModel{UserID: m.ID}),
		dao.IncludeFlag(dao.FlagMemberEstablishment),
	)
	if err != nil {
		errs = append(errs, err)
		return
	}
	for _, association := range associations {
		wg.Add(1)
		go func(association *dao.MemberModel) {
			group := dao.GroupModel{ID: association.GroupID}
			err := group.Query()
			if err != nil {
				errs = append(errs, err)
			}
			groups = append(groups, &Group{
				ID:   group.ID,
				Name: group.Name,
			})
			wg.Done()
		}(association)
	}
	wg.Wait()
	return
}

func (m *User) Joining() (groups GroupList, errs []error) {
	wg := sync.WaitGroup{}
	var associations dao.MemberList
	err := associations.Find(
		dao.Where(dao.MemberModel{UserID: m.ID}),
		dao.IncludeFlag(dao.FlagMemberRequest),
	)
	if err != nil {
		errs = append(errs, err)
		return
	}
	for _, association := range associations {
		wg.Add(1)
		go func(association *dao.MemberModel) {
			group := dao.GroupModel{ID: association.GroupID}
			err := group.Query()
			if err != nil {
				errs = append(errs, err)
			}
			groups = append(groups, &Group{
				ID:   group.ID,
				Name: group.Name,
			})
			wg.Done()
		}(association)
	}
	wg.Wait()
	return
}

func (m *User) Invites() (groups GroupList, errs []error) {
	wg := sync.WaitGroup{}
	var associations dao.MemberList
	err := associations.Find(
		dao.Where(dao.MemberModel{UserID: m.ID}),
		dao.IncludeFlag(dao.FlagMemberInvite),
	)
	if err != nil {
		errs = append(errs, err)
		return
	}
	for _, association := range associations {
		wg.Add(1)
		go func(association *dao.MemberModel) {
			group := dao.GroupModel{ID: association.GroupID}
			err := group.Query()
			if err != nil {
				errs = append(errs, err)
			}
			groups = append(groups, &Group{
				ID:   group.ID,
				Name: group.Name,
			})
			wg.Done()
		}(association)
	}
	wg.Wait()
	return
}

func (m *User) Join(groupID uint) error {
	association := dao.MemberModel{UserID: m.ID, GroupID: groupID}
	err := association.Query()
	if err == gorm.ErrRecordNotFound {
		association.Flag = dao.FlagMemberRequest
		return association.Create()
	} else if err != nil {
		return err
	}
	if association.Flag&dao.FlagMemberRequest != 0 {
		return errors.New("already request")
	}
	if association.Flag&dao.FlagMemberEstablishment != 0 {
		return errors.New("already joined")
	}
	if association.Flag&dao.FlagMemberInvite != 0 {
		association.Flag |= dao.FlagMemberEstablishment
		association.Flag &^= dao.FlagMemberInvite
	} else {
		association.Flag |= dao.FlagMemberRequest
	}
	return association.ForceUpdate()
}
