package modes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func SaveSettings(settings ServerSettings) error {
	os.Truncate("settings.json", 0)

	file, err := os.OpenFile("settings.json", os.O_WRONLY, 0600)
	defer file.Close()

	if err != nil {
		return fmt.Errorf("settings.json couldn't be opened")
	}

	settingsMarshalled, _ := json.MarshalIndent(settings, "", "  ")
	settingsString := string(settingsMarshalled)
	if _, err := file.WriteString(settingsString); err != nil {
		return fmt.Errorf("settings cannot be written")
	} else {
		// fmt.Println(settingsString)
	}

	return nil
}

func ReadSettings() ServerSettings {
	defaultSettings := ServerSettings{
		DbHost:     "localhost",
		DbPort:     6379,
		IpHost:     "localhost",
		IpPort:     8080,
		TTLMinutes: 10080,
	}
	file, err := os.Open("settings.json")
	defer file.Close()
	if err == nil {
		var settings ServerSettings
		if settingsBytes, err := ioutil.ReadAll(file); err == nil {
			if err := json.Unmarshal(settingsBytes, &settings); err == nil {
				return settings
			}
		}
	}
	return defaultSettings
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
