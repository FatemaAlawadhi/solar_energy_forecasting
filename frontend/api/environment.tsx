import { axiosInstance } from './config';

interface EnvironmentData {
  co2OffsetAwali: string;
  co2OffsetRefinery: string;
  co2OffsetUOB: string;
  totalCO2Offset: string;
  equivalentTreesAwali: string;
  equivalentTreesRefinery: string;
  equivalentTreesUOB: string;
  equivalentTreesTotal: string;
}

export const fetchEnvironmentData = async (): Promise<EnvironmentData> => {
  try {
    const response = await axiosInstance.get('/environment-impact');
    return response.data;
  } catch (error) {
    console.error('Error fetching environment data:', error);
    throw error;
  }
};