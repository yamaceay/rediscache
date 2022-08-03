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
		return nil
	}
	if mode == modes.ClientMode {
		if err := modes.StartClient(settings, filters...); err != nil {
			return fmt.Errorf("client couldn't be started: %s", err)
		}
		return nil
	}
	return fmt.Errorf("no valid mode given")
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
	mode := flag.String("mode", modes.ServerMode, fmt.Sprintf("%s / %s", modes.ServerMode, modes.ClientMode))
	method := flag.String("method", "get", fmt.Sprintf("%s / %s / %s", "get", "set", "delete"))
	key := flag.String("key", "", "")
	value := flag.String("value", "", "")
	db := flag.String("db", "0", "")

	flag.Parse()

	filters := []modes.Filter{
		{Key: "key", Value: *key},
		{Key: "value", Value: *value},
		{Key: "method", Value: *method},
		{Key: "db", Value: *db},
	}

	return *mode, filters
}
