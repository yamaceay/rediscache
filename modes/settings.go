package modes

import (
	"fmt"
	"strconv"
	"strings"
)

var ServerMode string = "server"
var ClientMode string = "client"

type ServerSettings struct {
	DbHost     string `json:"dbHost"`
	DbPort     int    `json:"dbPort"`
	IpHost     string `json:"ipHost"`
	IpPort     int    `json:"ipPort"`
	TTLMinutes int    `json:"ttlMinutes"`
}

type ClientRequest struct {
	Method string `json:"method"`
	Key    string `json:"key"`
	Value  string `json:"value"`
	Db     int    `json:"db"`
}

func (s ServerSettings) DbAddress() string {
	return ToAddress(s.DbHost, s.DbPort)
}

func (s ServerSettings) IpAddress() string {
	return ToAddress(s.IpHost, s.IpPort)
}

func ToAddress(host string, port int) string {
	portString := fmt.Sprintf("%d", port)
	address := strings.Join([]string{host, portString}, ":")
	return address
}

func FromAddress(address string) (string, int) {
	hostAndPort := strings.Split(address, ":")
	host, portString := hostAndPort[0], hostAndPort[1]
	port, _ := strconv.Atoi(portString)
	return host, port
}
