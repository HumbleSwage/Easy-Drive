package types

type Response struct {
	Code   int         `json:"code" form:"code"`
	Status string      `json:"status" form:"status"`
	Data   interface{} `json:"data" form:"data"`
	Info   string      `json:"info" form:"info"`
	Error  string      `json:"error" form:"error"`
}

type DataList struct {
	Item  interface{} `json:"item"`
	Total int         `json:"total"`
}

type TokenData struct {
	User        interface{} `json:"user"`
	AccessToken string      `json:"access_token"`
}
