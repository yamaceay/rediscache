package modes

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
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
	DbHost     string `yaml:"DB_HOST"`
	DbPort     int    `yaml:"DB_PORT"`
	IpHost     string `yaml:"IP_HOST"`
	IpPort     int    `yaml:"IP_PORT"`
	TTLMinutes int    `yaml:"TTL_MINUTES"`
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

func ReadSettings() ServerSettings {
	settings := ServerSettings{
		DbHost:     "cache",
		DbPort:     6379,
		IpHost:     "localhost",
		IpPort:     8080,
		TTLMinutes: 10080,
	}
	if file, err := os.Open("settings.yml"); err != nil {
		fmt.Printf("settings.yml cannot be opened: %s", err)
	} else if settingsBytes, err := ioutil.ReadAll(file); err != nil {
		fmt.Printf("settings.yml cannot be read: %s", err)
	} else if err := yaml.Unmarshal(settingsBytes, &settings); err != nil {
		fmt.Printf("settings cannot be parsed: %s", err)
	}

	return settings
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
