package web

type (
	Receive struct {
		Type string `json:"type"`
		Data any    `json:"data"`
	}
	User struct {
		UID uint `json:"uid" validate:"required"`
	}
	Group struct {
		GID uint `json:"gid" validate:"required"`
	}
	Message struct {
		RID  uint   `json:"rid"`
		CID  uint   `json:"cid"`
		Text string `json:"text" validate:"required"`
	}
)
