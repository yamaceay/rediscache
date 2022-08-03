package modes

import (
	"fmt"
	"io"
	"net/http"
)

var CLIENT_SETTINGS Settings

func StartClient(settings Settings, filters ...Filter) error {
	CLIENT_SETTINGS = settings

	method, _ := getFilter("method", filters)
	key, _ := getFilter("key", filters)
	value, _ := getFilter("value", filters)
	db, _ := getFilter("db", filters)

	ipAddress := CLIENT_SETTINGS.IpAddress()
	fullPath := fmt.Sprintf("http://%s/%s/%s", ipAddress, db, method)

	req, err := http.NewRequest("GET", fullPath, nil)
	if err != nil {
		return fmt.Errorf("client cannot create a request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	q := req.URL.Query()
	q.Add("key", key)
	q.Add("value", value)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("server returned error: %w", err)
	}
	defer resp.Body.Close()
	bodyString, _ := io.ReadAll(resp.Body)
	fmt.Println(string(bodyString))
	return nil
}

func getFilter(key string, filters []Filter) (string, error) {
	for _, filter := range filters {
		if filter.Key == key {
			return filter.Value, nil
		}
	}
	return "", fmt.Errorf("no filter found")
}
