package modes

import "fmt"

type Settings struct {
	DbHost     string `json:"dbHost"`
	DbPort     int    `json:"dbPort"`
	IpHost     string `json:"ipHost"`
	IpPort     int    `json:"ipPort"`
	TTLMinutes int    `json:"ttlMinutes"`
}

type Filter struct {
	Key   string
	Value string
}

var ServerMode string = "server"
var ClientMode string = "client"

func (s Settings) DbAddress() string {
	return fmt.Sprintf("%s:%d", s.DbHost, s.DbPort)
}

func (s Settings) IpAddress() string {
	return fmt.Sprintf("%s:%d", s.IpHost, s.IpPort)
}
