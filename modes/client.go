package modes

import (
	"fmt"
	"io"
	"net/http"
)

func StartClient(settings ServerSettings, request ClientRequest) error {
	ipAddress := settings.IpAddress()
	fullPath := fmt.Sprintf("http://%s/%d/%s", ipAddress, request.Db, request.Method)

	req, err := http.NewRequest("GET", fullPath, nil)
	if err != nil {
		return fmt.Errorf("client cannot create a request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	q := req.URL.Query()
	q.Add("key", request.Key)
	q.Add("value", request.Value)
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
