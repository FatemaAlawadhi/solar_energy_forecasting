package structure

// ForecastResult represents the power generation forecast data structure
type ForecastResult struct {
    Year      int
    Month     int
    Actual    float64
    Predicted float64
}

type PowerGenerationResponse struct {
	LastMonth string                 `json:"lastMonth"`
	LastYear  string                 `json:"lastYear"`
	Forecast  []ForecastResult       `json:"forecast"`
}