package modes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

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
