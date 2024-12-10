"use client";
import { useEffect, useState } from "react";
import Breadcrumb from "@/components/Breadcrumbs/Breadcrumb";
import DefaultLayout from "@/components/Layouts/DefaultLaout";
import ChartThree from "@/components/Charts/ChartThree";
import { fetchSystemConfigurationData } from "../../../api/system-config";
import { SystemDataItem } from "../../../api/system-config";

const SystemConfigurationPage = () => {
  const [systemData, setSystemData] = useState<{
    installedCapacity: number[];
    numberOfPanels: number[];
    labels: string[];
  }>({
    installedCapacity: [],
    numberOfPanels: [],
    labels: []
  });

  const [totalSystemData, setTotalSystemData] = useState<SystemDataItem | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const data: SystemDataItem[] = await fetchSystemConfigurationData();
        const totalSystemData = data.find((item: SystemDataItem) => item.Name === "Total System");
        const filteredData = data.filter((item: SystemDataItem) => item.Name !== "Total System");
        setSystemData({
          installedCapacity: filteredData.map((item: SystemDataItem) => parseFloat(item.InstalledCapacity.toFixed(1))),
          numberOfPanels: filteredData.map((item: SystemDataItem) => item.NumberOfPanels),
          labels: filteredData.map((item: SystemDataItem) => item.Name)
        });

        if (totalSystemData) {
          totalSystemData.InstalledCapacity = parseFloat(totalSystemData.InstalledCapacity.toFixed(1));
          setTotalSystemData(totalSystemData);
        }
      } catch (error) {
        console.error('Error:', error);
      }
    };
    fetchData();
  }, []);

  return (
    <DefaultLayout>
      <Breadcrumb pageName="System Configuration" />
      <div className="flex flex-col gap-8">
        <ChartThree
          series={systemData.installedCapacity}
          labels={systemData.labels}
          title="Installed Capacity (kw)"
          unit="kW"
        />
        <ChartThree
          series={systemData.numberOfPanels}
          labels={systemData.labels}
          title="Number of Panels"
          unit="panels"
        />
      </div>
    </DefaultLayout>
  );
};

export default SystemConfigurationPage;
