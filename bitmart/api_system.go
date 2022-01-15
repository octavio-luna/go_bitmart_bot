package bitmart

import (
	"encoding/json"
	"time"
)

func (cloudClient *CloudClient) GetSystemTime() (*CloudResponse, string, error) {
	var cloudResponse CloudResponse

	if _, err := cloudClient.Request(GET, API_SYSTEM_TIME_URL, nil, NONE, &cloudResponse); err != nil {
		return nil, "", err
	}

	var now Time
	json.Unmarshal([]byte(cloudResponse.response), &now)
	now.Data.ServerTime = now.Data.ServerTime / 1000
	currentTime := time.Unix(now.Data.ServerTime, 0)

	return &cloudResponse, currentTime.String(), nil
}

func (cloudClient *CloudClient) GetSystemService() (*CloudResponse, error) {
	var cloudResponse CloudResponse

	if _, err := cloudClient.Request(GET, API_SYSTEM_SERVICE_URL, nil, NONE, &cloudResponse); err != nil {
		return nil, err
	}

	return &cloudResponse, nil
}
