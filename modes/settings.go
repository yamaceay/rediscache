package modes

import (
	"fmt"
	"strconv"
	"strings"
)

var ServerMode string = "server"
var ClientMode string = "client"
var ConfigMode string = "config"

type Settings struct {
	DbHost     string `json:"dbHost"`
	DbPort     int    `json:"dbPort"`
	IpHost     string `json:"ipHost"`
	IpPort     int    `json:"ipPort"`
	TTLMinutes int    `json:"ttlMinutes"`
}

func (s Settings) DbAddress() string {
	return toAddress(s.DbHost, s.DbPort)
}

func (s Settings) IpAddress() string {
	return toAddress(s.IpHost, s.IpPort)
}

func toAddress(host string, port int) string {
	portString := fmt.Sprintf("%d", port)
	address := strings.Join([]string{host, portString}, ":")
	return address
}

func fromAddress(address string) (string, int) {
	hostAndPort := strings.Split(address, ":")
	host, portString := hostAndPort[0], hostAndPort[1]
	port, _ := strconv.Atoi(portString)
	return host, port
}

type Filter struct {
	Key   string
	Value string
}

func getFilter(key string, filters []Filter) (string, error) {
	for _, filter := range filters {
		if filter.Key == key {
			return filter.Value, nil
		}
	}
	return "", fmt.Errorf("no filter found")
}
