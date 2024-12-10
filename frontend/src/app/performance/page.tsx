"use client";
import { useEffect, useState } from "react";
import Breadcrumb from "@/components/Breadcrumbs/Breadcrumb";
import DefaultLayout from "@/components/Layouts/DefaultLaout";
import { fetchPerformanceData } from "../../../api/performance";
import dynamic from 'next/dynamic';

const ChartPerformance = dynamic(() => import("@/components/Charts/ChartPerformance"), {
  ssr: false,
  loading: () => <div className="h-[400px] flex items-center justify-center">Loading chart...</div>
});

const ChartOne = dynamic(() => import("@/components/Charts/ChartOne"), {
  ssr: false,
  loading: () => <div className="h-[400px] flex items-center justify-center">Loading chart...</div>
});

const ChartTwo = dynamic(() => import("@/components/Charts/ChartTwo"), {
  ssr: false,
  loading: () => <div className="h-[400px] flex items-center justify-center">Loading chart...</div>
});

const Performance = () => {
  const [performanceData, setPerformanceData] = useState<{
    overall: {
      pr: { labels: string[]; data: number[] };
      cf: { labels: string[]; data: number[] };
      outputPerPV: { labels: string[]; data: number[] };
    };
    monthly: {
      generation: {
        labels: string[];
        datasets: { name: string; data: number[] }[];
      };
      pr: {
        labels: string[];
        datasets: { name: string; data: number[] }[];
      };
      cf: {
        labels: string[];
        datasets: { name: string; data: number[] }[];
      };
      outputPerPV: {
        labels: string[];
        datasets: { name: string; data: number[] }[];
      };
    };
    yearly: {
      pr: {
        labels: string[];
        datasets: { name: string; data: number[] }[];
      };
      cf: {
        labels: string[];
        datasets: { name: string; data: number[] }[];
      };
      outputPerPV: {
        labels: string[];
        datasets: { name: string; data: number[] }[];
      };
    };
  }>({
    overall: {
      pr: { labels: [], data: [] },
      cf: { labels: [], data: [] },
      outputPerPV: { labels: [], data: [] }
    },
    monthly: {
      generation: { labels: [], datasets: [] },
      pr: { labels: [], datasets: [] },
      cf: { labels: [], datasets: [] },
      outputPerPV: { labels: [], datasets: [] }
    },
    yearly: {
      pr: { labels: [], datasets: [] },
      cf: { labels: [], datasets: [] },
      outputPerPV: { labels: [], datasets: [] }
    }
  });

  useEffect(() => {
    const fetchData = async () => {
      try {
        const data = await fetchPerformanceData();
        
        // Process overall performance data
        const overallPR = {
          labels: data.overall_performance.map(item => item.location_name),
          data: data.overall_performance.map(item => item.performance_ratio)
        };
        const overallCF = {
          labels: data.overall_performance.map(item => item.location_name),
          data: data.overall_performance.map(item => item.capacity_factor)
        };
        const overallOutputPerPV = {
          labels: data.overall_performance.map(item => item.location_name),
          data: data.overall_performance.map(item => item.output_per_pv)
        };

        // Process monthly data
        const monthlyLabels = Array.from(new Set(data.monthly_generation.map(
          item => `${item.year}-${String(item.month).padStart(2, '0')}`
        ))).sort();

        const monthlyGeneration = {
          labels: monthlyLabels,
          datasets: [
            {
              name: "Actual Generation",
              data: data.monthly_generation.map(item => item.actual_kwh)
            },
            {
              name: "Theoretical Generation",
              data: data.monthly_generation.map(item => item.theoretical_kwh)
            }
          ]
        };

        const monthlyPR = {
          labels: monthlyLabels,
          datasets: Array.from(new Set(data.monthly_performance.map(item => item.location_name))).map(location => ({
            name: location,
            data: monthlyLabels.map(label => {
              const [year, month] = label.split('-').map(Number);
              const matchingData = data.monthly_performance.find(item => 
                item.location_name === location && 
                item.year === year &&
                item.month === month
              );
              return matchingData ? matchingData.performance_ratio : 0;
            })
          }))
        };

        const monthlyCF = {
          labels: monthlyLabels,
          datasets: Array.from(new Set(data.monthly_performance.map(item => item.location_name))).map(location => ({
            name: location,
            data: monthlyLabels.map(label => {
              const matchingData = data.monthly_performance.find(item => 
                item.location_name === location && 
                `${item.year}-${String(item.month).padStart(2, '0')}` === label
              );
              return matchingData ? matchingData.capacity_factor : 0;
            })
          }))
        };

        const monthlyOutputPerPV = {
          labels: monthlyLabels,
          datasets: Array.from(new Set(data.monthly_performance.map(item => item.location_name))).map(location => ({
            name: location,
            data: monthlyLabels.map(label => {
              const matchingData = data.monthly_performance.find(item => 
                item.location_name === location && 
                `${item.year}-${String(item.month).padStart(2, '0')}` === label
              );
              return matchingData ? matchingData.output_per_pv : 0;
            })
          }))
        };

        // Process yearly data
        const yearlyLabels = Array.from(new Set(data.yearly_performance.map(item => (item.year ?? 0).toString())));

        const yearlyPR = {
          labels: yearlyLabels,
          datasets: data.yearly_performance.reduce((acc, item) => {
            const locationIndex = acc.findIndex(d => d.name === item.location_name);
            if (locationIndex === -1) {
              acc.push({
                name: item.location_name,
                data: [item.performance_ratio]
              });
            } else {
              acc[locationIndex].data.push(item.performance_ratio);
            }
            return acc;
          }, [] as { name: string; data: number[] }[])
        };

        const yearlyCF = {
          labels: yearlyLabels,
          datasets: data.yearly_performance.reduce((acc, item) => {
            const locationIndex = acc.findIndex(d => d.name === item.location_name);
            if (locationIndex === -1) {
              acc.push({
                name: item.location_name,
                data: [item.capacity_factor]
              });
            } else {
              acc[locationIndex].data.push(item.capacity_factor);
            }
            return acc;
          }, [] as { name: string; data: number[] }[])
        };

        const yearlyOutputPerPV = {
          labels: yearlyLabels,
          datasets: data.yearly_performance.reduce((acc, item) => {
            const locationIndex = acc.findIndex(d => d.name === item.location_name);
            if (locationIndex === -1) {
              acc.push({
                name: item.location_name,
                data: [item.output_per_pv]
              });
            } else {
              acc[locationIndex].data.push(item.output_per_pv);
            }
            return acc;
          }, [] as { name: string; data: number[] }[])
        };

        // Set the processed data
        setPerformanceData({
          overall: {
            pr: overallPR,
            cf: overallCF,
            outputPerPV: overallOutputPerPV
          },
          monthly: {
            generation: monthlyGeneration,
            pr: monthlyPR,
            cf: monthlyCF,
            outputPerPV: monthlyOutputPerPV
          },
          yearly: {
            pr: yearlyPR,
            cf: yearlyCF,
            outputPerPV: yearlyOutputPerPV
          }
        });

      } catch (error) {
        console.error('Error:', error);
      }
    };
    fetchData();
  }, []);

  return (
    <DefaultLayout>
      <Breadcrumb pageName="Performance Analysis" />
      
      {/* Overall Performance Section */}
      <div className="mb-8">
        <h2 className="text-2xl font-bold mb-4">Overall Performance (5 Years)</h2>
        <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
          <ChartTwo 
            title="Performance Ratio" 
            labels={performanceData.overall.pr.labels}
            data={performanceData.overall.pr.data}
            hideXAxis={true}
            tooltipTitle="Performance Ratio"
          />
          <ChartTwo 
            title="Capacity Factor" 
            labels={performanceData.overall.cf.labels}
            data={performanceData.overall.cf.data}
            hideXAxis={true}
            tooltipTitle="Capacity Factor"
          />
          <ChartTwo 
            title="Output per PV" 
            labels={performanceData.overall.outputPerPV.labels}
            data={performanceData.overall.outputPerPV.data}
            hideXAxis={true}
            tooltipTitle="Power generation/pv (kwh)"
          />
        </div>
      </div>

      {/* Monthly Performance Section */}
      <div className="mb-8">
        <h2 className="text-2xl font-bold mb-4">Monthly Performance</h2>
        
        {/* Generation Analysis Charts */}
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 mb-8">
          {performanceData.monthly.generation.datasets[0]?.data && 
            performanceData.overall.pr.labels.map((location, index) => {
              // Filter data for this specific location
              const locationData = {
                labels: performanceData.monthly.generation.labels.filter((_, i) => 
                  i % performanceData.overall.pr.labels.length === index
                ),
                datasets: [
                  {
                    name: "Actual Generation",
                    data: performanceData.monthly.generation.datasets[0].data.filter((_, i) => 
                      i % performanceData.overall.pr.labels.length === index
                    )
                  },
                  {
                    name: "Theoretical Generation",
                    data: performanceData.monthly.generation.datasets[1].data.filter((_, i) => 
                      i % performanceData.overall.pr.labels.length === index
                    )
                  }
                ],
                title: `Theoretical vs Actual Generation - ${location}`
              };

              return (
                <ChartOne 
                  key={location}
                  data={locationData}
                />
              );
            })
          }
        </div>
      </div>

      {/* Yearly Performance Section */}
      <div className="mb-8">
        <h2 className="text-2xl font-bold mb-4">Yearly Performance</h2>
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          {/* Performance Ratio & Capacity Factor Charts */}
          <ChartPerformance data={{
            labels: performanceData.yearly.pr.labels,
            datasets: performanceData.yearly.pr.datasets,
            title: "Yearly Performance Ratio",
            tooltipSuffix: "%"
          }} />
          
          <ChartPerformance data={{
            labels: performanceData.yearly.cf.labels,
            datasets: performanceData.yearly.cf.datasets,
            title: "Yearly Capacity Factor",
            tooltipSuffix: "%"
          }} />
        </div>

        {/* Output per PV Chart */}
        <div className="grid grid-cols-1 gap-4 mt-4">
          <ChartOne data={{
            labels: performanceData.yearly.outputPerPV.labels,
            datasets: performanceData.yearly.outputPerPV.datasets,
            title: "Yearly Output per PV",
            tooltipSuffix: " kWh"
          }} />
        </div>
      </div>
    </DefaultLayout>
  );
};

export default Performance;