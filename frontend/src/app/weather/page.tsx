"use client";
import { useEffect, useState } from "react";
import Breadcrumb from "@/components/Breadcrumbs/Breadcrumb";
import DefaultLayout from "@/components/Layouts/DefaultLaout";
import { fetchWeatherData } from "../../../api/weather";
import dynamic from 'next/dynamic';

const ChartTwo = dynamic(() => import("@/components/Charts/ChartTwo"), {
  ssr: false,
  loading: () => <div className="h-[400px] flex items-center justify-center">Loading chart...</div>
});

const ChartDualAxis = dynamic(() => import("@/components/Charts/ChartDualAxis"), {
  ssr: false,
  loading: () => <div className="h-[400px] flex items-center justify-center">Loading chart...</div>
});

const WeatherImpact = () => {
    const [weatherCharts, setWeatherCharts] = useState([
        {
            title: "Power Generation vs Average Sunshine Duration",
            labels: [] as string[],
            datasets: [
                { name: "Power Generation (kWh)", data: [] as number[] },
                { name: "Average Sunshine Duration (s)", data: [] as number[] }
            ]
        },
        {
            title: "Power Generation vs Temperature",
            labels: [] as string[],
            datasets: [
                { name: "Power Generation (kWh)", data: [] as number[] },
                { name: "Minimum Temperature (°C)", data: [] as number[] },
                { name: "Average Temperature (°C)", data: [] as number[] },
                { name: "Maximum Temperature (°C)", data: [] as number[] }
            ]
        },
        {
            title: "Power Generation vs Solar Irradiance",
            labels: [] as string[],
            datasets: [
                { name: "Power Generation (kWh)", data: [] as number[] },
                { name: "Average Solar Irradiance (W/m²)", data: [] as number[] }
            ]
        },
        {
            title: "Power Generation vs Relative Humidity",
            labels: [] as string[],
            datasets: [
                { name: "Power Generation (kWh)", data: [] as number[] },
                { name: "Average Relative Humidity (%)", data: [] as number[] }
            ]
        },
        {
            title: "Power Generation vs Cloud Cover",
            labels: [] as string[],
            datasets: [
                { name: "Power Generation (kWh)", data: [] as number[] },
                { name: "Average Cloud Cover (%)", data: [] as number[] }
            ]
        },
        {
            title: "Power Generation vs Wind Speed",
            labels: [] as string[],
            datasets: [
                { name: "Power Generation (kWh)", data: [] as number[] },
                { name: "Average Wind Speed (km/h)", data: [] as number[] }
            ]
        },
        {
            title: "Power Generation vs Rainfall",
            labels: [] as string[],
            datasets: [
                { name: "Power Generation (kWh)", data: [] as number[] },
                { name: "Cumulative Rainfall (mm)", data: [] as number[] }
            ]
        }
    ]);
    const [featureImportance, setFeatureImportance] = useState({
        labels: [] as string[],
        values: [] as number[]
    });

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetchWeatherData();
                const { weatherData, featureImportance } = response;

                const sortedData = weatherData.sort((a, b) =>
                    a.year === b.year ? a.month - b.month : a.year - b.year
                );

                setFeatureImportance({
                    labels: featureImportance.map(item => item.featureName),
                    values: featureImportance.map(item => item.importanceValue)
                });

                const labels = sortedData.map(item =>
                    `${item.year}-${item.month.toString().padStart(2, '0')}`
                );

                setWeatherCharts(charts => charts.map(chart => {
                    const powerData = sortedData.map(item => Math.round(item.totalPowerGeneration));

                    switch (chart.title) {
                        case "Power Generation vs Average Sunshine Duration":
                            return {
                                ...chart,
                                labels,
                                datasets: [
                                    { name: "Power Generation (kWh)", data: powerData },
                                    {
                                        name: "Average Sunshine Duration (s)",
                                        data: sortedData.map(item => Math.round(item.avgSunshineDuration))
                                    }
                                ]
                            };
                        case "Power Generation vs Temperature":
                            return {
                                ...chart,
                                labels,
                                datasets: [
                                    { name: "Power Generation (kWh)", data: powerData },
                                    {
                                        name: "Average Temperature (°C)",
                                        data: sortedData.map(item => Math.round(item.avgTemperature))
                                    }
            
                                ]
                            };
                        case "Power Generation vs Solar Irradiance":
                            return {
                                ...chart,
                                labels,
                                datasets: [
                                    { name: "Power Generation (kWh)", data: powerData },
                                    {
                                        name: "Average Solar Irradiance (W/m²)",
                                        data: sortedData.map(item => Math.round(item.avgSolarIrradiance))
                                    }
                                ]
                            };
                        case "Power Generation vs Relative Humidity":
                            return {
                                ...chart,
                                labels,
                                datasets: [
                                    { name: "Power Generation (kWh)", data: powerData },
                                    {
                                        name: "Average Relative Humidity (%)",
                                        data: sortedData.map(item => Math.round(item.avgRelativeHumidity))
                                    }
                                ]
                            };
                        case "Power Generation vs Cloud Cover":
                            return {
                                ...chart,
                                labels,
                                datasets: [
                                    { name: "Power Generation (kWh)", data: powerData },
                                    {
                                        name: "Average Cloud Cover (%)",
                                        data: sortedData.map(item => Math.round(item.avgCloudCover))
                                    }
                                ]
                            };
                        case "Power Generation vs Wind Speed":
                            return {
                                ...chart,
                                labels,
                                datasets: [
                                    { name: "Power Generation (kWh)", data: powerData },
                                    {
                                        name: "Average Wind Speed (km/h)",
                                        data: sortedData.map(item => Math.round(item.avgWindSpeed))
                                    }
                                ]
                            };
                        case "Power Generation vs Rainfall":
                            return {
                                ...chart,
                                labels,
                                datasets: [
                                    { name: "Power Generation (kWh)", data: powerData },
                                    {
                                        name: "Cumulative Rainfall (mm)",
                                        data: sortedData.map(item => Math.round(item.cumulativeRainfall))
                                    }
                                ]
                            };
                        default:
                            return chart;
                    }
                }));
            } catch (error) {
                console.error('Error:', error);
            }
        };
        fetchData();
    }, []);

    return (
        <DefaultLayout>
            <Breadcrumb pageName="Weather Impact Analysis" />

            <div className="mb-8">
                <ChartTwo
                    title="Feature Importance Analysis"
                    labels={featureImportance.labels}
                    data={featureImportance.values}
                    hideXAxis={false}
                    tooltipTitle="Feature Importance"
                />
            </div>

            <div className="grid gap-16 pb-16">
                {weatherCharts.map((chartData, index) => (
                    <div key={index} className="h-[400px]">
                        <ChartDualAxis data={chartData} />
                    </div>
                ))}
            </div>
        </DefaultLayout>
    );
};

export default WeatherImpact;