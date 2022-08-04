package modes

import (
	"fmt"
	"io"
	"net/http"
)

func StartClient(settings ServerSettings, request ClientRequest, response *string) error {
	ipAddress := settings.IpAddress()
	fullPath := fmt.Sprintf("http://%s/%d/%s", ipAddress, request.Db, request.Method)

	req, err := http.NewRequest("GET", fullPath, nil)
	if err != nil {
		return fmt.Errorf("GET %s - no valid HTTP request: %s", fullPath, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	q := req.URL.Query()
	if request.Key != "" {
		q.Add("key", request.Key)
	}
	if request.Value != "" {
		q.Add("value", request.Value)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return fmt.Errorf("server returns %d error: %s", resp.StatusCode, err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	*response = string(bodyBytes)
	return nil
}
