package models

import (
	"IM/services"
	"errors"
	"github.com/rsooio/gsync"
	"sync"
)

var userList gsync.Map[uint, *User]

const (
	userStep uint = iota
	groupStep
	allStep
)

type (
	User struct {
		Info        *services.User
		ClientList  gsync.Map[uint, *Client]
		RoomList    gsync.Map[uint, *Room]
		clientIDMax uint
	}
)

func calcRoomID(baseID, step uint) uint {
	return baseID*allStep - step
}

func (u *User) Join(user *services.User) (errList []error, err error) {
	if u.Info != nil {
		return errList, errors.New("user already joined")
	}
	wg := sync.WaitGroup{}
	u.Info = user
	userList.Store(user.ID, u)
	friends, errs := user.Friends()
	errList = append(errList, errs...)
	if len(friends) > 0 {
		wg.Add(1)
		go func() {
			for _, friend := range friends {
				wg.Add(1)
				go func(friend *services.Friend) {
					roomID := calcRoomID(friend.ID, userStep)
					room, ok := roomList.Load(roomID)
					if !ok {
						room = &Room{
							ID:       roomID,
							UserList: gsync.Map[uint, *User],
						}
					}
					u.RoomList.Store(roomID, room)
					room.UserList.Store(u.Info.ID, u)
					wg.Done()
				}(friend)
			}
			wg.Done()
		}()
	}
	groups, errs := user.Groups()
	errList = append(errList, errs...)
	if len(groups) > 0 {
		wg.Add(1)
		go func() {
			for _, group := range groups {
				wg.Add(1)
				go func(group *services.Group) {
					roomID := calcRoomID(group.ID, groupStep)
					room, ok := roomList.Load(roomID)
					if !ok {
						room = &Room{
							ID:       roomID,
							UserList: gsync.Map[uint, *User],
						}
					}
					u.RoomList.Store(roomID, room)
					room.UserList.Store(u.Info.ID, u)
					wg.Done()
				}(group)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return
}
