package models

import "github.com/rsooio/gsync"

var roomList gsync.Map[uint, *Room]

type (
	Room struct {
		ID       uint
		UserList gsync.Map[uint, *User]
	}
)
