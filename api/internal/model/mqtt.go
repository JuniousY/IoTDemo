package model

type ConnectMsg struct {
	UserName string `json:"username"`
	Ts       int64  `json:"ts"`
	Address  string `json:"ipaddress"`
	ClientID string `json:"clientid"`
	Reason   string `json:"reason"`
}
