package weather

type Weather struct {
	CloudPct    int     `json:"cloud_pct,omitempty"`
	Temp        int     `json:"temp,omitempty"`
	FeelsLike   int     `json:"feels_like,omitempty"`
	Humidity    int     `json:"humidity,omitempty"`
	MinTemp     int     `json:"min_temp,omitempty"`
	MaxTemp     int     `json:"max_temp,omitempty"`
	WindSpeed   float64 `json:"wind_speed,omitempty"`
	WindDegrees int     `json:"wind_degrees,omitempty"`
	Sunrise     int     `json:"sunrise,omitempty"`
	Sunset      int     `json:"sunset,omitempty"`
}
