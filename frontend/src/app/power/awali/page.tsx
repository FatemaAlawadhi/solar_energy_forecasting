"use client";
import { useEffect, useState } from "react";
import Breadcrumb from "@/components/Breadcrumbs/Breadcrumb";
import DefaultLayout from "@/components/Layouts/DefaultLaout";
import dynamic from 'next/dynamic';
import { fetchAwaliPowerGenerationData } from "../../../../api/power";

interface PowerGenerationData {
  lastMonth: string;
  lastYear: string;
  forecast: {
    Year: number;
    Month: number;
    Actual: number;
    Predicted: number;
  }[];
}

const ChartOne = dynamic(() => import("@/components/Charts/ChartOne"), {
  ssr: false,
  loading: () => <div className="h-[400px] flex items-center justify-center">Loading chart...</div>
});

const AwaliPowerGeneration = () => {
  const [powerData, setPowerData] = useState<PowerGenerationData>({
    lastMonth: "0 kWh",
    lastYear: "0 kWh",
    forecast: []
  });

  const [chartData, setChartData] = useState<{
    labels: string[];
    datasets: { name: string; data: (number | null)[] }[];
    title: string;
    tooltipSuffix?: string;
  }>({
    labels: [],
    datasets: [],
    title: "Awali Power Generation Forecast",
    tooltipSuffix: " kWh"
  });

  useEffect(() => {
    const fetchData = async () => {
      try {
        const data = await fetchAwaliPowerGenerationData();
        setPowerData(data);

        if (data.forecast?.length > 0) {
          const sortedData = [...data.forecast].sort((a, b) => {
            if (a.Year !== b.Year) return a.Year - b.Year;
            return a.Month - b.Month;
          });

          const firstNonZeroIndex = sortedData.findIndex(item => item.Actual > 0);
          
          setChartData({
            labels: sortedData.map(item => `${item.Year}-${String(item.Month).padStart(2, '0')}`),
            datasets: [
              {
                name: "Actual Power Generation (kWh)",
                data: sortedData.map(item => item.Actual || null)
              },
              {
                name: "Predicted Power Generation (kWh)",
                data: sortedData.map((item, index) => 
                  index >= firstNonZeroIndex && item.Predicted ? item.Predicted : null
                )
              }
            ],
            title: "Awali Power Generation Forecast",
            tooltipSuffix: " kWh"
          });
        }
      } catch (error) {
        console.error('Error:', error);
      }
    };
    fetchData();
  }, []);

  const statsData = [
    {
      icon: (
        <svg width="26" height="26" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M13 2L3 14H12L11 22L21 10H12L13 2Z" fill="white"/>
        </svg>
      ),
      color: "#3FD97F",
      title: "Last Month",
      value: powerData.lastMonth,
    },
    {
      icon: (
        <svg width="26" height="26" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M13 2L3 14H12L11 22L21 10H12L13 2Z" fill="white"/>
        </svg>
      ),
      color: "#FF9C55",
      title: "Last Year",
      value: powerData.lastYear,
    },
  ];

  return (
    <DefaultLayout>
      <Breadcrumb pageName="Awali Power Generation" />
      <div className="flex flex-col gap-8">
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 md:gap-6 xl:grid-cols-2 2xl:gap-7.5">
          {statsData.map((data, index) => (
            <div
              key={index}
              className="rounded-[10px] bg-white p-6 shadow-1 dark:bg-gray-dark"
            >
              <div
                className="flex h-14.5 w-14.5 items-center justify-center rounded-full"
                style={{ backgroundColor: data.color }}
              >
                {data.icon}
              </div>

              <div className="mt-6 flex items-end justify-between">
                <div className="w-full">
                  <h4 className="mb-1.5 text-heading-6 font-bold text-dark dark:text-white text-center">
                    {data.value}
                  </h4>
                  <span className="text-xl font-medium block text-center">{data.title}</span>
                </div>
              </div>
            </div>
          ))}
        </div>

        <div className="h-[400px] mb-8">
          {chartData.datasets.length > 0 ? (
            <ChartOne data={chartData} />
          ) : (
            <div className="flex items-center justify-center h-full">
              Loading chart data...
            </div>
          )}
        </div>
      </div>
    </DefaultLayout>
  );
};

export default AwaliPowerGeneration;
