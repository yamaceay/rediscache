package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/modes"

	flag "github.com/spf13/pflag"
)

func main() {
	settings := setupSettings()
	mode, filters := parseArgs()
	if err := handleMode(mode, settings, filters...); err != nil {
		panic(err)
	}
}

func handleMode(mode string, settings modes.Settings, filters ...modes.Filter) error {
	if mode == modes.ServerMode {
		if err := modes.StartServer(settings, filters...); err != nil {
			return fmt.Errorf("server couldn't be started: %s", err)
		}
	} else if mode == modes.ClientMode {
		if err := modes.StartClient(settings, filters...); err != nil {
			return fmt.Errorf("client couldn't be started: %s", err)
		}
	} else if mode == modes.ConfigMode {
		if err := modes.CreateSettings(filters...); err != nil {
			return fmt.Errorf("settings.json couldn't be configured: %s", err)
		}
	} else {
		return fmt.Errorf("no valid mode given")
	}
	return nil
}

func setupSettings() modes.Settings {
	var settings modes.Settings
	if bytes, err := ioutil.ReadFile("settings.json"); err == nil {
		if err := json.Unmarshal(bytes, &settings); err != nil {
			// fmt.Println("settings cannot be read")
		}
	} else {
		settings = modes.Settings{
			IpHost:     "localhost",
			IpPort:     8080,
			DbHost:     "localhost",
			DbPort:     6379,
			TTLMinutes: 10080,
		}
		// fmt.Println("settings.json couldn't be found")
	}
	return settings
}

func parseArgs() (string, []modes.Filter) {
	mode := flag.StringP("mode", "M", modes.ServerMode, fmt.Sprintf("%s / %s", modes.ServerMode, modes.ClientMode))

	// Config
	dbAddress := flag.String("dbAddress", "localhost:6379", fmt.Sprintf("<database address>"))
	ipAddress := flag.String("ipAddress", "localhost:8080", fmt.Sprintf("<application address>"))
	ttlMinutes := flag.String("ttlMinutes", "10080", fmt.Sprintf("<time to live in minutes>"))

	// Client
	method := flag.StringP("method", "X", "", fmt.Sprintf("%s | %s | %s | %s", "", "get", "set", "delete"))
	key := flag.StringP("key", "k", "", "")
	value := flag.StringP("value", "v", "", "")
	db := flag.String("db", "0", "")

	flag.Parse()

	filters := []modes.Filter{
		{Key: "key", Value: *key},
		{Key: "value", Value: *value},
		{Key: "method", Value: *method},
		{Key: "db", Value: *db},

		{Key: "dbAddress", Value: *dbAddress},
		{Key: "ipAddress", Value: *ipAddress},
		{Key: "ttlMinutes", Value: *ttlMinutes},
	}

	return *mode, filters
}
