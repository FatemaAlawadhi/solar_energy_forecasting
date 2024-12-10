package structure

type Generation struct {
    Year        int     `json:"year"`
    Month       int     `json:"month"`
    LocationID  int     `json:"location_id"`
    LocationName string `json:"location_name"`
    ActualKWH   float64 `json:"actual_kwh"`
    TheoreticalKWH float64 `json:"theoretical_kwh"`
}

type Performance struct {
    Year            int     `json:"year,omitempty"`
    Month           int     `json:"month,omitempty"`
    StartYear       int     `json:"start_year,omitempty"`
    EndYear         int     `json:"end_year,omitempty"`
    LocationID      int     `json:"location_id"`
    LocationName    string  `json:"location_name"`
    PerformanceRatio float64 `json:"performance_ratio"`
    CapacityFactor   float64 `json:"capacity_factor"`
    OutputPerPV      float64 `json:"output_per_pv"`
}

type PerformanceResponse struct {
    MonthlyGeneration []Generation  `json:"monthly_generation"`
    MonthlyPerformance []Performance `json:"monthly_performance"`
    YearlyPerformance  []Performance `json:"yearly_performance"`
    OverallPerformance []Performance `json:"overall_performance"`
}
