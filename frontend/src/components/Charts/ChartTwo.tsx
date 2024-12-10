import { ApexOptions } from "apexcharts";
import React from "react";
import ReactApexChart from "react-apexcharts";

interface ChartTwoProps {
  title: string;
  labels: string[];
  data: number[];
  hideXAxis?: boolean;
  tooltipTitle?: string;
}

const ChartTwo: React.FC<ChartTwoProps> = ({ 
  title, 
  labels, 
  data, 
  hideXAxis = false, 
  tooltipTitle = "Feature Importance"
}) => {
  const series = [
    {
      name: tooltipTitle,
      data: data,
    }
  ];

  const options: ApexOptions = {
    colors: ["#5750F1"],
    chart: {
      fontFamily: "Satoshi, sans-serif",
      type: "bar",
      height: 335,
      toolbar: {
        show: false,
      },
      zoom: {
        enabled: false,
      },
    },
    plotOptions: {
      bar: {
        horizontal: true,
        borderRadius: 3,
        columnWidth: "60%",
        distributed: true,
      },
    },
    dataLabels: {
      enabled: true,
      formatter: (val) => Number(val).toFixed(3),
      textAnchor: 'start',
      style: {
        fontSize: '12px',
      },
    },
    grid: {
      show: true,
      xaxis: {
        lines: {
          show: true,
        },
      },
      yaxis: {
        lines: {
          show: false,
        },
      },
    },
    xaxis: {
      categories: labels,
      labels: {
        formatter: (val) => Number(val).toFixed(3),
        show: !hideXAxis,
      },
      axisBorder: {
        show: !hideXAxis,
      },
      axisTicks: {
        show: !hideXAxis,
      },
    },
    yaxis: {
      labels: {
        style: {
          fontSize: '12px',
        },
        show: true,
      },
      axisBorder: {
        show: false,
      },
      axisTicks: {
        show: false,
      },
    },
    legend: {
      show: false,
    },
  };

  return (
    <div className="rounded-[10px] bg-white px-7.5 pt-7.5 shadow-1 dark:bg-gray-dark dark:shadow-card">
      <div className="mb-4">
        <h4 className="text-body-2xlg font-bold text-dark dark:text-white text-center">
          {title}
        </h4>
      </div>
      <div>
        <div id="chartTwo">
          <ReactApexChart
            options={options}
            series={series}
            type="bar"
            height={500}  // Increased height to accommodate all features
          />
        </div>
      </div>
    </div>
  );
};

export default ChartTwo;
