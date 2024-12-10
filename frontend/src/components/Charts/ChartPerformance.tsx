import ChartOne from "./ChartOne";
import { ApexOptions } from "apexcharts";

interface ChartPerformanceProps {
  data: {
    title: string;
    labels: string[];
    datasets: {
      name: string;
      data: number[];
    }[];
    tooltipSuffix?: string;
  };
}

const ChartPerformance: React.FC<ChartPerformanceProps> = ({ data }) => {
  // Sort datasets by their average values (highest to lowest)
  const sortedDatasets = [...data.datasets].sort((a, b) => {
    const aAvg = a.data.reduce((sum, val) => sum + val, 0) / a.data.length;
    const bAvg = b.data.reduce((sum, val) => sum + val, 0) / b.data.length;
    return bAvg - aAvg;
  });

  const modifiedData = {
    ...data,
    datasets: sortedDatasets.map(dataset => ({
      name: dataset.name,
      data: dataset.data.map(val => val * 100), // Convert to percentage
    })),
  };

  const customOptions: ApexOptions = {
    colors: [
      '#3C50E0', '#80CAEE', '#259AE6', '#375E83'
    ],
    tooltip: {
      shared: true,
      intersect: false,
      y: {
        formatter: (value: number | null) => 
          value === null ? 'No data' : `${value.toFixed(2)}%`
      }
    },
    yaxis: {
      title: {
        text: data.title
      },
      labels: {
        formatter: (value) => `${value.toFixed(1)}%`
      },
      forceNiceScale: true
    }
  };

  return <ChartOne data={modifiedData} customOptions={customOptions} />;
};

export default ChartPerformance;