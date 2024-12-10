package api

import (
	"backend/pkg/calculation"
	structure "backend/pkg/struct"
	"encoding/json"
	"net/http"
)

func EnvironmentalImpact(w http.ResponseWriter, r *http.Request) {
	co2OffsetAwali, co2OffsetRefinery, co2OffsetUOB, totalCO2Offset := calculation.CO2Offset()
	equivalentTreesAwali, equivalentTreesRefinery, equivalentTreesUOB, equivalentTreesTotal := calculation.EquivalentTrees(co2OffsetAwali, co2OffsetRefinery, co2OffsetUOB)

	env := structure.EnvironmentalImpact{
		Co2OffsetAwali:       calculation.FormatCO2Number(co2OffsetAwali), 
		Co2OffsetRefinery:    calculation.FormatCO2Number(co2OffsetRefinery),
		Co2OffsetUOB:         calculation.FormatCO2Number(co2OffsetUOB),
		TotalCO2Offset:       calculation.FormatCO2Number(totalCO2Offset),
		EquivalentTreesAwali: calculation.FormatTreeNumber(equivalentTreesAwali), 
		EquivalentTreesRefinery: calculation.FormatTreeNumber(equivalentTreesRefinery),
		EquivalentTreesUOB:   calculation.FormatTreeNumber(equivalentTreesUOB),
		EquivalentTreesTotal:  calculation.FormatTreeNumber(equivalentTreesTotal),
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(env)
}
