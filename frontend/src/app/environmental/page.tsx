"use client";
import { useEffect, useState } from "react";
import Breadcrumb from "@/components/Breadcrumbs/Breadcrumb";
import DefaultLayout from "@/components/Layouts/DefaultLaout";
import { fetchEnvironmentData } from "../../../api/environment";
import ChartThree from "@/components/Charts/ChartThree";

const EnvironmentalImpact = () => {
  const [environmentData, setEnvironmentData] = useState({
    co2OffsetAwali: "0 kg",
    co2OffsetRefinery: "0 kg",
    co2OffsetUOB: "0 kg",
    totalCO2Offset: "0 kg",
    equivalentTreesAwali: "0 trees",
    equivalentTreesRefinery: "0 trees",
    equivalentTreesUOB: "0 trees",
    equivalentTreesTotal: "0 trees"
  });

  useEffect(() => {
    const fetchData = async () => {
      try {
        const data = await fetchEnvironmentData();
        setEnvironmentData(data);
      } catch (error) {
        console.error('Error:', error);
      }
    };
    fetchData();
  }, []);

  // Update the data arrays with fetched values
  const co2Data = [
    {
      icon: (
        <svg width="26" height="26" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <text x="1" y="17" fill="white" fontSize="12" fontWeight="bold">CO₂</text>
        </svg>
      ),
      color: "#3FD97F",
      title: "Awali",
      value: environmentData.co2OffsetAwali,
    },
    {
      icon: (
        <svg width="26" height="26" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <text x="1" y="17" fill="white" fontSize="12" fontWeight="bold">CO₂</text>
        </svg>
      ),
      color: "#FF9C55",
      title: "Refinery",
      value: environmentData.co2OffsetRefinery,
    },
    {
      icon: (
        <svg width="26" height="26" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <text x="1" y="17" fill="white" fontSize="12" fontWeight="bold">CO₂</text>
        </svg>
      ),
      color: "#8155FF",
      title: "UOB",
      value: environmentData.co2OffsetUOB,
    },
    {
      icon: (
        <svg width="26" height="26" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <text x="1" y="17" fill="white" fontSize="12" fontWeight="bold">CO₂</text>
        </svg>
      ),
      color: "#3C50E0",
      title: "Total",
      value: environmentData.totalCO2Offset,
    },
  ];

  const treeData = [
    {
      icon: (
        <svg width="26" height="26" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M4 12C4 8 7 5 12 5C17 5 20 8 20 12C20 16 17 16 12 16C7 16 4 16 4 12Z" fill="white"/>
          <path d="M10 16L9 19H15L14 16H10Z" fill="white"/>
        </svg>
      ),
      color: "#18BFFF",
      title: "Awali",
      value: environmentData.equivalentTreesAwali,
    },
    {
      icon: (
        <svg width="26" height="26" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M4 12C4 8 7 5 12 5C17 5 20 8 20 12C20 16 17 16 12 16C7 16 4 16 4 12Z" fill="white"/>
          <path d="M10 16L9 19H15L14 16H10Z" fill="white"/>
        </svg>
      ),
      color: "#FF9C55",
      title: "Refinery",
      value: environmentData.equivalentTreesRefinery,
    },
    {
      icon: (
        <svg width="26" height="26" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M4 12C4 8 7 5 12 5C17 5 20 8 20 12C20 16 17 16 12 16C7 16 4 16 4 12Z" fill="white"/>
          <path d="M10 16L9 19H15L14 16H10Z" fill="white"/>
        </svg>
      ),
      color: "#8155FF",
      title: "UOB",
      value: environmentData.equivalentTreesUOB,
    },
    {
      icon: (
        <svg width="26" height="26" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M4 12C4 8 7 5 12 5C17 5 20 8 20 12C20 16 17 16 12 16C7 16 4 16 4 12Z" fill="white"/>
          <path d="M10 16L9 19H15L14 16H10Z" fill="white"/>
        </svg>
      ),
      color: "#3C50E0",
      title: "Total",
      value: environmentData.equivalentTreesTotal,
    },
  ];

  const co2Values = [
    parseFloat(environmentData.co2OffsetAwali.replace(/[^0-9.-]+/g, "")),
    parseFloat(environmentData.co2OffsetRefinery.replace(/[^0-9.-]+/g, "")),
    parseFloat(environmentData.co2OffsetUOB.replace(/[^0-9.-]+/g, ""))
  ];

  const treeValues = [
    parseFloat(environmentData.equivalentTreesAwali.replace(/[^0-9.-]+/g, "")),
    parseFloat(environmentData.equivalentTreesRefinery.replace(/[^0-9.-]+/g, "")),
    parseFloat(environmentData.equivalentTreesUOB.replace(/[^0-9.-]+/g, ""))
  ];

  return (
    <DefaultLayout>
      <Breadcrumb pageName="Environmental Impact" />
      <div className="flex flex-col gap-10">

        <div>
          <h4 className="mb-6 text-xl font-semibold text-black dark:text-white text-center">
            CO2 Offset
          </h4>
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2 md:gap-6 xl:grid-cols-3 2xl:gap-7.5">
            {co2Data.map((data, index) => (
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
        </div>

        <ChartThree
          series={co2Values}
          labels={["Awali", "Refinery", "UOB"]}
          title="CO2 Offset"
          unit="kg"
        />

        <div>
          <h4 className="mb-6 text-xl font-semibold text-black dark:text-white text-center">
            Equivalent Trees
          </h4>
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2 md:gap-6 xl:grid-cols-3 2xl:gap-7.5">
            {treeData.map((data, index) => (
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
        </div>

        <ChartThree
          series={treeValues}
          labels={["Awali", "Refinery", "UOB"]}
          title="Equivalent Trees"
          unit="trees"
        />
      </div>
    </DefaultLayout>
  );
};

export default EnvironmentalImpact;