import { axiosInstance } from './config';

export interface SystemDataItem {
  Name: string;
  InstalledCapacity: number;
  NumberOfPanels: number;
}

export const fetchSystemConfigurationData = async (): Promise<SystemDataItem[]> => {
  try {
    const response = await axiosInstance.get('http://localhost:8080/api/system-configuration');
    return response.data;
  } catch (error) {
    console.error('Error fetching system configuration data:', error);
    throw error;
  }
};
