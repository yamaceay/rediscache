package modes

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func CreateSettings(filters ...Filter) error {
	ipAddress, _ := getFilter("ipAddress", filters)
	ipHost, ipPort := fromAddress(ipAddress)

	dbAddress, _ := getFilter("dbAddress", filters)
	dbHost, dbPort := fromAddress(dbAddress)

	ttlMinutesString, _ := getFilter("ttlMinutes", filters)
	ttlMinutes, _ := strconv.Atoi(ttlMinutesString)

	os.Truncate("settings.json", 0)

	file, err := os.OpenFile("settings.json", os.O_WRONLY, 0600)
	defer file.Close()

	if err != nil {
		return fmt.Errorf("settings.json couldn't be opened")
	}

	settings := Settings{
		IpHost:     ipHost,
		IpPort:     ipPort,
		DbHost:     dbHost,
		DbPort:     dbPort,
		TTLMinutes: ttlMinutes,
	}

	settingsMarshalled, _ := json.MarshalIndent(settings, "", "  ")
	settingsString := string(settingsMarshalled)
	if _, err := file.WriteString(settingsString); err != nil {
		return fmt.Errorf("settings cannot be written")
	} else {
		fmt.Println(settingsString)
	}

	return nil
}
