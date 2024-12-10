package calculation

import (
	"backend/pkg/db/queries"
	"fmt"
	"strings"
)

func CO2Offset() (float64, float64,float64, float64) {
	totalUOB, totalRefinery, totalAwali := queries.TotalPowerGeneration()

	const carbonIntensityNaturalGas = 400.0 // gCO2/kWh

    // To calculate CO2 offsets in kilograms
    co2OffsetUOB := (totalUOB * carbonIntensityNaturalGas) / 1000.0 // kg CO2
    co2OffsetRefinery := (totalRefinery * carbonIntensityNaturalGas) / 1000.0 // kg CO2
    co2OffsetAwali := (totalAwali * carbonIntensityNaturalGas) / 1000.0 // kg CO2

    // Total CO2 offset
    totalCO2Offset := co2OffsetUOB + co2OffsetRefinery + co2OffsetAwali
	return co2OffsetAwali, co2OffsetRefinery, co2OffsetUOB, totalCO2Offset
}

func EquivalentTrees(co2OffsetAwali, co2OffsetRefinery, co2OffsetUOB float64) (float64, float64, float64, float64) {
	// 21 kg CO2 absorbed by one tree per year
	const co2AbsorptionPerTreePerYear = 21.0
	totalCO2AbsorbedByOneTreeInFiveYears := co2AbsorptionPerTreePerYear * 5

	// To calculate equivalent trees planted
	equivalentTreesAwali := co2OffsetAwali / totalCO2AbsorbedByOneTreeInFiveYears
	equivalentTreesRefinery := co2OffsetRefinery / totalCO2AbsorbedByOneTreeInFiveYears
	equivalentTreesUOB := co2OffsetUOB / totalCO2AbsorbedByOneTreeInFiveYears
	equivalentTreesTotal := equivalentTreesAwali + equivalentTreesRefinery +equivalentTreesUOB

	return equivalentTreesAwali, equivalentTreesRefinery, equivalentTreesUOB, equivalentTreesTotal
}

func FormatCO2Number(num float64) string {
    // To convert to string with 2 decimal places
    str := fmt.Sprintf("%.2f", num)
    
    // To split into integer and decimal parts
    parts := strings.Split(str, ".")
    intPart := parts[0]
    
    // To add commas to integer part
    var result []byte
    for i, j := len(intPart)-1, 0; i >= 0; i-- {
        if j > 0 && j%3 == 0 {
            result = append([]byte{','}, result...)
        }
        result = append([]byte{intPart[i]}, result...)
        j++
    }
    
    // Tp combine with decimal part and add kg
    return fmt.Sprintf("%s kg", string(result))
}

func FormatTreeNumber(num float64) string {
    if num >= 1e9 {
        return fmt.Sprintf("%.0f billion trees", num/1e9)
    } else if num >= 1e6 {
        return fmt.Sprintf("%.0f million trees", num/1e6)
    } else if num >= 1e3 {
        return fmt.Sprintf("%.0f thousand trees", num/1e3)
    }
    return fmt.Sprintf("%.0f trees", num)
}