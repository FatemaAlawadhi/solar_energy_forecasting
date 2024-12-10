import { ApexOptions } from "apexcharts";
import React from "react";
import ReactApexChart from "react-apexcharts";

interface ChartDualAxisProps {
  data: {
    labels: string[];
    datasets: {
      name: string;
      data: number[];
    }[];
    title: string;
  };
}

const ChartDualAxis: React.FC<ChartDualAxisProps> = ({ data }) => {
  const series = data.datasets.map(dataset => ({
    name: dataset.name,
    data: dataset.data,
  }));

  const options: ApexOptions = {
    legend: {
      show: true,
      position: "top",
      horizontalAlign: "left",
    },
    colors: ["#5750F1", "#0ABEF9"],
    chart: {
      fontFamily: "Satoshi, sans-serif",
      height: 310,
      type: "area",
      toolbar: {
        show: false,
      },
    },
    fill: {
      gradient: {
        opacityFrom: 0.55,
        opacityTo: 0,
      },
    },
    stroke: {
      curve: "smooth",
      width: 2
    },
    markers: {
      size: 0,
    },
    grid: {
      strokeDashArray: 5,
      xaxis: {
        lines: {
          show: false,
        },
      },
      yaxis: {
        lines: {
          show: true,
        },
      },
    },
    dataLabels: {
      enabled: false,
    },
    tooltip: {
      shared: true,
      intersect: false,
      y: {
        formatter: function (value: number | null) {
          return value === null ? 'No data' : `${value.toFixed(2)}`;
        },
      },
    },
    xaxis: {
      type: "category",
      categories: data.labels,
      axisBorder: {
        show: false,
      },
      axisTicks: {
        show: false,
      },
      labels: {
        show: false,
      }
    },
    yaxis: [
      {
        title: {
          text: data.datasets[0].name,
        },
        labels: {
          formatter: (value) => value.toFixed(0)
        },
      },
      {
        opposite: true,
        title: {
          text: data.datasets[1].name,
        },
        labels: {
          formatter: (value) => value.toFixed(0)
        },
      }
    ],
  };

  return (
    <div className="col-span-12 rounded-[10px] bg-white px-7.5 pb-6 pt-7.5 shadow-1 dark:bg-gray-dark dark:shadow-card xl:col-span-7">
      <div className="mb-3.5 flex justify-center">
        <div>
          <h4 className="text-xl font-semibold text-dark dark:text-white">
            {data.title}
          </h4>
        </div>
      </div>
      <div>
        <div className="-ml-4 -mr-5">
          <ReactApexChart
            options={options}
            series={series}
            type="area"
            height={310}
          />
        </div>
      </div>
    </div>
  );
};

export default ChartDualAxis;