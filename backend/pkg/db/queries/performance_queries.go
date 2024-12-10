package queries

const (
    GetMonthlyGeneration = `
        SELECT 
            mg.year,
            mg.month,
            mg.location_id,
            l.name as location_name,
            mg.actual_kwh,
            mg.theoretical_kwh
        FROM monthly_generation mg
        JOIN locations l ON mg.location_id = l.id
        ORDER BY mg.year, mg.month, mg.location_id
    `

    GetMonthlyPerformance = `
        SELECT 
            mp.year,
            mp.month,
            mp.location_id,
            l.name as location_name,
            mp.performance_ratio,
            mp.capacity_factor,
            mp.output_per_pv
        FROM monthly_performance mp
        JOIN locations l ON mp.location_id = l.id
        ORDER BY mp.year, mp.month, mp.location_id
    `

    GetYearlyPerformance = `
        SELECT 
            yp.year,
            yp.location_id,
            l.name as location_name,
            yp.performance_ratio,
            yp.capacity_factor,
            yp.output_per_pv
        FROM yearly_performance yp
        JOIN locations l ON yp.location_id = l.id
        ORDER BY yp.year, yp.location_id
    `

    GetOverallPerformance = `
        SELECT 
            op.start_year,
            op.end_year,
            op.location_id,
            l.name as location_name,
            op.performance_ratio,
            op.capacity_factor,
            op.output_per_pv
        FROM overall_performance op
        JOIN locations l ON op.location_id = l.id
        ORDER BY op.location_id
    `
)