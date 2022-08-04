package modes

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var ServerMode string = "server"
var ClientMode string = "client"

type ClientRequest struct {
	Method string `json:"method"`
	Key    string `json:"key"`
	Value  string `json:"value"`
	Db     int    `json:"db"`
}
type ServerSettings struct {
	DbHost     string `json:"dbHost"`
	DbPort     int    `json:"dbPort"`
	IpHost     string `json:"ipHost"`
	IpPort     int    `json:"ipPort"`
	TTLMinutes int    `json:"ttlMinutes"`
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

func NewSettings(dbHost string, dbPort int, ipHost string, ipPort int, ttlMinutes int) ServerSettings {
	return ServerSettings{
		DbHost:     dbHost,
		DbPort:     dbPort,
		IpHost:     ipHost,
		IpPort:     ipPort,
		TTLMinutes: ttlMinutes,
	}
}

func ReadSettings() (string, int, string, int, int) {
	dbHost := getenv("DB_HOST", "cache")
	dbPort, _ := strconv.Atoi(getenv("DB_PORT", "6379"))
	ipHost := getenv("IP_HOST", "localhost")
	ipPort, _ := strconv.Atoi(getenv("IP_PORT", "8080"))
	ttlMinutes, _ := strconv.Atoi(getenv("TTL_MINUTES", "10080"))
	return dbHost, dbPort, ipHost, ipPort, ttlMinutes
}

func NewRequest(method string, key string, value string, db int) ClientRequest {
	return ClientRequest{
		Method: method,
		Key:    key,
		Value:  value,
		Db:     db,
	}
}

func getenv(key string, defvalue string) string {
	if value, yes := os.LookupEnv(key); yes {
		return value
	} else {
		return defvalue
	}
}
