package main

import (
	"fmt"
	"main/modes"
	"strconv"

	flag "github.com/spf13/pflag"
)

func StartServer(dbHost string, dbPort int, ipHost string, ipPort int, ttlMinutes int) error {
	settings := modes.NewSettings(dbHost, dbPort, ipHost, ipPort, ttlMinutes)
	return modes.StartServer(settings)
}

func StartClient(method string, key string, value string, db int, response *string) error {
	settings := modes.NewSettings(modes.ReadSettings())
	request := modes.NewRequest(method, key, value, db)
	return modes.StartClient(settings, request, response)
}

func main() {
	mode, params := parseArgs()
	if err := handleMode(mode, params); err != nil {
		fmt.Printf("program ended: %s", err)
	}
}

func handleMode(mode string, params map[string]string) error {
	var body string
	if mode == modes.ServerMode {
		if err := StartServer(modes.ReadSettings()); err != nil {
			return fmt.Errorf("server cannot be started: %s", err)
		}
	} else if mode == modes.ClientMode {
		method := params["method"]
		key, value := params["key"], params["value"]
		db, _ := strconv.Atoi(params["db"])

		if err := StartClient(method, key, value, db, &body); err != nil {
			return fmt.Errorf("client cannot be started: %s", err)
		}
	} else {
		return fmt.Errorf("unknown mode %s: select either \"%s\" or \"%s\"", mode, "server", "client")
	}
	fmt.Println(body)
	return nil
}

func parseArgs() (string, map[string]string) {
	mode := flag.StringP("mode", "M", modes.ServerMode, fmt.Sprintf("%s / %s", modes.ServerMode, modes.ClientMode))

	method := flag.StringP("method", "X", "help", fmt.Sprintf("%s | %s | %s", "get", "set", "help"))
	key := flag.StringP("key", "k", "", "")
	value := flag.StringP("value", "v", "", "")
	db := flag.String("db", "0", "")

	flag.Parse()

	params := map[string]string{
		"key":    *key,
		"value":  *value,
		"method": *method,
		"db":     *db,
	}

	return *mode, params
}
