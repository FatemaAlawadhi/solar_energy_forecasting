package structure

type EnvironmentalImpact struct {
	Co2OffsetAwali    string `json:"co2OffsetAwali"`
	Co2OffsetRefinery string `json:"co2OffsetRefinery"`
	Co2OffsetUOB      string `json:"co2OffsetUOB"`
	TotalCO2Offset    string `json:"totalCO2Offset"`
	EquivalentTreesAwali string `json:"equivalentTreesAwali"`
	EquivalentTreesRefinery string `json:"equivalentTreesRefinery"`
	EquivalentTreesUOB string `json:"equivalentTreesUOB"`
	EquivalentTreesTotal string `json:"equivalentTreesTotal"`
}
