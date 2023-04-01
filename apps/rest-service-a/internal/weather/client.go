package weather

import (
	"encoding/json"
	"github.com/eldius/golang-observability-poc/apps/rest-service-a/internal/config"
	"net/http"
)

func GetWeather(city string) (*Weather, error) {

	c := http.Client{}

	resp, err := c.Get(config.GetServiceBEndpoint())
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var w Weather

	if err := json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return nil, err
	}

	return &w, nil
}
