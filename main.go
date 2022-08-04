package main

import (
	"fmt"
	"main/modes"
	"strconv"

	flag "github.com/spf13/pflag"
)

func StartServer(dbAddress string, ipAddress string, ttlMinutes int) error {
	dbHost, dbPort := modes.FromAddress(dbAddress)
	ipHost, ipPort := modes.FromAddress(ipAddress)

	settings := modes.ServerSettings{
		DbHost:     dbHost,
		DbPort:     dbPort,
		IpHost:     ipHost,
		IpPort:     ipPort,
		TTLMinutes: ttlMinutes,
	}

	if err := modes.SaveSettings(settings); err != nil {
		return fmt.Errorf("settings cannot be saved: %s", err)
	}

	return modes.StartServer(settings)
}

func StartClient(method string, key string, value string, db int, response *string) error {
	settings := modes.ReadSettings()
	request := modes.ClientRequest{
		Method: method,
		Key:    key,
		Value:  value,
		Db:     db,
	}
	return modes.StartClient(settings, request, response)
}

func main() {
	mode, params := parseArgs()
	if err := handleMode(mode, params); err != nil {
		fmt.Printf("program ended: %s", err)
	}
}

func handleMode(mode string, params map[string]string) error {
	if mode == modes.ServerMode {
		dbAddress := params["dbAddress"]
		ipAddress := params["ipAddress"]
		ttlMinutes, _ := strconv.Atoi(params["ttlMinutes"])

		if err := StartServer(dbAddress, ipAddress, ttlMinutes); err != nil {
			return fmt.Errorf("server cannot be started: %s", err)
		}
	} else if mode == modes.ClientMode {
		method := params["method"]
		key, value := params["key"], params["value"]
		db, _ := strconv.Atoi(params["db"])
		var body string

		if err := StartClient(method, key, value, db, &body); err != nil {
			return fmt.Errorf("client cannot be started: %s", err)
		}
		fmt.Println(body)
	} else {
		return fmt.Errorf("unknown mode")
	}
	return nil
}

func parseArgs() (string, map[string]string) {
	mode := flag.StringP("mode", "M", modes.ServerMode, fmt.Sprintf("%s / %s", modes.ServerMode, modes.ClientMode))

	// Server
	dbAddress := flag.String("dbAddress", "localhost:6379", fmt.Sprintf("<database address>"))
	ipAddress := flag.String("ipAddress", "localhost:8080", fmt.Sprintf("<application address>"))
	ttlMinutes := flag.String("ttlMinutes", "10080", fmt.Sprintf("<time to live in minutes>"))

	// Client
	method := flag.StringP("method", "X", "", fmt.Sprintf("%s | %s | %s", "", "get", "set"))
	key := flag.StringP("key", "k", "", "")
	value := flag.StringP("value", "v", "", "")
	db := flag.String("db", "0", "")

	flag.Parse()

	params := map[string]string{
		"key":    *key,
		"value":  *value,
		"method": *method,
		"db":     *db,

		"dbAddress":  *dbAddress,
		"ipAddress":  *ipAddress,
		"ttlMinutes": *ttlMinutes,
	}

	return *mode, params
}
