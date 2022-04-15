package controllers

import (
	"IM/models"
	"IM/services"
	"IM/utils"
	"IM/web"
	"encoding/json"
)

func msgHandler(msg *[]byte, client *models.Client) {
	var message any
	err := json.Unmarshal(*msg, &message)
	if err != nil {
		client.ERR(err.Error())
		return
	}
	var receive web.Receive
	err = utils.Decode(&message, &receive)
	if err != nil {
		client.ERR(err.Error())
		return
	}
	switch receive.Type {
	case "message":
		go message(&receive, client)
	case "follow":
		go follow(&receive, client)
	case "reject follow":
		go rejectFollow(&receive, client)
	case "cancel follow":
		go cancelFollow(&receive, client)
	case "join group":
		go joinGroup(&receive, client)
	case "invite group":
		go inviteGroup(&receive, client)
	default:
		client.ERR("unknown message type")
	}
}

func message(receive *web.Receive, client *models.Client) {
}

func follow(receive *web.Receive, client *models.Client) {
	var follow web.User
	err := utils.Decode(receive, &follow)
	if err != nil {
		client.ERR(err.Error())
		return
	}
	err = client.User.Info.Follow(follow.UID)
	if err != nil {
		client.ERR(err.Error())
		return
	}
	client.JSON(&web.Reply{
		Type: "success",
		Data: "follow success",
	})
}

func rejectFollow(receive *web.Receive, client *models.Client) {
	var follow web.User
	err := utils.Decode(receive, &follow)
	if err != nil {
		client.ERR(err.Error())
		return
	}
	err = client.User.Info.RejectFollow(follow.UID)
	if err != nil {
		client.ERR(err.Error())
		return
	}
	client.JSON(&web.Reply{
		Type: "success",
		Data: "reject follow success",
	})
}

func cancelFollow(receive *web.Receive, client *models.Client) {
	var follow web.User
	err := utils.Decode(receive, &follow)
	if err != nil {
		client.ERR(err.Error())
		return
	}
	err = client.User.Info.CancelFollow(follow.UID)
	if err != nil {
		client.ERR(err.Error())
		return
	}
	client.JSON(&web.Reply{
		Type: "success",
		Data: "cancel follow success",
	})
}

func joinGroup(receive *web.Receive, client *models.Client) {
	var group web.Group
	err := utils.Decode(receive, &group)
	if err != nil {
		client.ERR(err.Error())
		return
	}
	err = client.User.Info.Join(group.GID)
	if err != nil {
		client.ERR(err.Error())
		return
	}
	client.JSON(&web.Reply{
		Type: "success",
		Data: "join group success",
	})
}

func inviteGroup(receive *web.Receive, client *models.Client) {
	var group web.Group
	err := utils.Decode(receive, &group)
	if err != nil {
		client.ERR(err.Error())
		return
	}
	g := services.Group{ID: group.GID}
	err = g.Invite(client.User.Info.ID)
	if err != nil {
		client.ERR(err.Error())
		return
	}
	client.JSON(&web.Reply{
		Type: "success",
		Data: "invite group send success",
	})
}
