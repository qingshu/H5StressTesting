package main

type Result struct {
	State bool   `json:"s"`
	Cmd   string `json:"cmd"`
}

//登陆回调
type LoginResult struct {
	Result
	M struct {
		Guid       string `json:"guid"`
		Session    int64  `json:"session"`
		Createrole int    `json:"createrole"`
		Guidestep  string `json:"guidestep"`
		Uistep     int    `json:"uistep"`
	}
}
