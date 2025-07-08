package model

import "encoding/json"

type ConnectMsg struct {
	UserName string `json:"username"`
	Ts       int64  `json:"ts"`
	Address  string `json:"ipaddress"`
	ClientID string `json:"clientid"`
	Reason   string `json:"reason"`
}

type UploadDataMsg struct {
	Topic  string          `json:"topic"`
	Type   string          `json:"type"`
	Device string          `json:"device"`
	Data   json.RawMessage `json:"data"`
}

// SensorData type 为 sensor
type SensorData struct {
	Temperature string `json:"temperature"`
	Humidity    int    `json:"humidity"`
}
