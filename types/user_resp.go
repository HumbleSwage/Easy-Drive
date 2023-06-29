package types

type UserLoginResp struct {
	UserId    string      `json:"userId"`
	NickName  string      `json:"nickName"`
	Admin     bool        `json:"admin"`
	Authority bool        `json:"authority"`
	Token     string      `json:"token"`
	Avatar    interface{} `json:"avatar"`
}

type UserSpaceResp struct {
	UseSpace   int64 `json:"useSpace"`
	TotalSpace int64 `json:"totalSpace"`
}
