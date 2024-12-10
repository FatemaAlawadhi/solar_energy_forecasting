package structure

type HourlyData struct {
	Time                   []string  `json:"time"`
	Temperature2m          []float64 `json:"temperature_2m"`
	RelativeHumidity2m     []float64 `json:"relative_humidity_2m"`
	CloudCover             []float64 `json:"cloud_cover"`
	WindSpeed10m           []float64 `json:"wind_speed_10m"`
	DirectNormalIrradiance []float64 `json:"direct_normal_irradiance"`
}

type DailyData struct {
	Time    []string  `json:"time"`
	Sunrise           []string  `json:"sunrise"`
    Sunset            []string  `json:"sunset"`
    DaylightDuration  []float64 `json:"daylight_duration"`
    SunshineDuration  []float64 `json:"sunshine_duration"`
	RainSum []float64 `json:"rain_sum"`
}

type APIResponse struct {
	Hourly HourlyData `json:"hourly"`
	Daily  DailyData  `json:"daily"`
}

type WeatherImpactData struct {
	Year                     int     `json:"year"`
	Month                    int     `json:"month"`
	AvgSunshineDuration     int     `json:"avgSunshineDuration"`
	AvgDaylightDuration     int     `json:"avgDaylightDuration"`
	MinTemperature          float64 `json:"minTemperature"`
	AvgTemperature          float64 `json:"avgTemperature"`
	MaxTemperature          float64 `json:"maxTemperature"`
	AvgSolarIrradiance      float64 `json:"avgSolarIrradiance"`
	AvgRelativeHumidity     float64 `json:"avgRelativeHumidity"`
	AvgCloudCover           float64 `json:"avgCloudCover"`
	AvgWindSpeed            float64 `json:"avgWindSpeed"`
	CumulativeRainfall      float64 `json:"cumulativeRainfall"`
	TotalPowerGeneration    float64 `json:"totalPowerGeneration"`
}