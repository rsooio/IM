package models

import (
	"IM/services"
	"IM/utils"
	"IM/web"
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type (
	Client struct {
		ID         uint
		ClientName string
		User       *User
		Conn       *websocket.Conn
	}
)

func (c *Client) Join(u *services.User) error {
	if c.ID != 0 {
		return errors.New("client already joined")
	}
	user, ok := userList.Load(u.ID)
	if !ok {
		errs, err := user.Join(u)
		if err != nil {
			return err
		}
		utils.Display(errs)
	}
	c.User = user
	user.clientIDMax++
	c.ID = user.clientIDMax
	user.ClientList.Store(c.ID, c)
	return nil
}

func (c *Client) Leave() {
	c.User.ClientList.Delete(c.ID)
	empty := true
	c.User.ClientList.Range(func(k uint, value *Client) bool {
		empty = false
		return false
	})
	if empty {
		userList.Delete(c.User.Info.ID)
		wg := sync.WaitGroup{}
		c.User.RoomList.Range(func(_ uint, v *Room) bool {
			wg.Add(1)
			go func(room *Room) {
				v.UserList.Delete(c.User.Info.ID)
				empty := true
				v.UserList.Range(func(k uint, value *User) bool {
					empty = false
					return false
				})
				if empty {
					roomList.Delete(v.ID)
				}
				wg.Done()
			}(v)
			return true
		})
	}
}

func (c *Client) JSON(reply *web.Reply) {
	err := c.Conn.WriteJSON(reply)
	if err != nil {
		utils.Display(err)
	}
}

func (c *Client) ERR(err any) {
	c.JSON(&web.Reply{
		Type: "error",
		Data: err,
	})
}

func (c *Client) SendMsg(msg *web.Message) {
	if msg.CID != 0 {
		client, ok := c.User.ClientList.Load(msg.CID)
		if !ok {
			c.ERR("client not found")
			return
		}
		client.JSON(&web.Reply{
			Type: "message",
			Data: msg,
		})
		c.JSON(&web.Reply{
			Type: "success",
			Data: "message sent",
		})
	}
	if msg.RID != 0 {
		room, ok := c.User.RoomList.Load(msg.RID)
		if !ok {
			c.ERR("room not found")
			return
		}
		wg := sync.WaitGroup{}
		room.UserList.Range(func(_ uint, user *User) bool {
			if user.Info.ID != c.User.Info.ID {
				wg.Add(1)
				go func() {
					user.ClientList.Range(func(_ uint, client *Client) bool {
						wg.Add(1)
						go func(client *Client) {
							client.JSON(&web.Reply{
								Type: "message",
								Data: msg,
							})
							wg.Done()
						}(client)
						return true
					})
					wg.Done()
				}()
			}
			return true
		})
		wg.Wait()
		c.JSON(&web.Reply{
			Type: "success",
			Data: "message sent",
		})
	}
	c.JSON(&web.Reply{
		Type: "error",
		Data: "no client or room id",
	})
}
