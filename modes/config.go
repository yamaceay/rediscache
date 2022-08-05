package modes

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type ServerSettings struct {
	RedisHost     string `yaml:"REDIS_HOST"`
	RedisPort     int    `yaml:"REDIS_PORT"`
	RedisPassword string `yaml:"REDIS_PASSWORD"`

	ServerHost string `yaml:"SERVER_HOST"`
	ServerPort int    `yaml:"SERVER_PORT"`

	TTLMinutes int `yaml:"TTL_MINUTES"`
}

func (s ServerSettings) RedisAddress() string {
	return ToAddress(s.RedisHost, s.RedisPort)
}

func (s ServerSettings) ServerAddress() string {
	return ToAddress(s.ServerHost, s.ServerPort)
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

func NewSettings(redisHost string, redisPort int, redisPassword string, serverHost string, serverPort int, ttlMinutes int) ServerSettings {
	return ServerSettings{
		RedisHost:     redisHost,
		RedisPort:     redisPort,
		RedisPassword: redisPassword,

		ServerHost: serverHost,
		ServerPort: serverPort,

		TTLMinutes: ttlMinutes,
	}
}

func ReadSettings() ServerSettings {
	settings := ServerSettings{
		RedisHost:     "cache",
		RedisPort:     6379,
		RedisPassword: "",

		ServerHost: "localhost",
		ServerPort: 8080,

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
