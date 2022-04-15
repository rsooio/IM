package controllers

import (
	"IM/models"
	"IM/services"
	"IM/utils"
	"IM/web"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func WSConnector(c *gin.Context) {
	user, ok := wsAuthenticator(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, web.Reply{
			Type: "error",
			Data: "Unauthorized",
		})
		return
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, web.Reply{
			Type: "error",
			Data: "upgrade error",
		})
		return
	}
	client := models.Client{
		ClientName: c.GetString("client-name"),
		Conn:       ws,
	}
	err = client.Join(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, web.Reply{
			Type: "error",
			Data: "join error",
		})
		return
	}
	wsListener(&client)
}

func wsAuthenticator(c *gin.Context) (user services.User, ok bool) {
	err := user.LoginByToken(c.Request.Header.Get("token"))
	if err == nil {
		return user, true
	}
	user.ID = c.GetUint("uid")
	err = user.LoginByPwd(c.GetString("username"), c.GetString("password"))
	if err == nil {
		return user, true
	}
	return services.User{}, false
}

func wsListener(client *models.Client) {
	mt, message, err := client.Conn.ReadMessage()
	if err != nil || mt >= 1000 {
		if err != nil {
			utils.Display(err)
		}
		client.Leave()
		err := client.Conn.Close()
		if err != nil {
			utils.Display(err)
		}
		return
	}
	go wsListener(client)
	msgHandler(&message, client)
}
