import { axiosInstance, API_URL } from './config';

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

export const fetchPowerGenerationData = async (): Promise<PowerGenerationData> => {
  try {
    const response = await axiosInstance.get('/total-power-generation');
    return response.data;
  } catch (error) {
    console.error('Error fetching power generation data:', error);
    throw error;
  }
};

export const fetchUOBPowerGenerationData = async (): Promise<PowerGenerationData> => {
  try {
    const response = await axiosInstance.get('/uob-power-generation');
    return response.data;
  } catch (error) {
    console.error('Error fetching UOB power generation data:', error);
    throw error;
  }
};

export const fetchAwaliPowerGenerationData = async (): Promise<PowerGenerationData> => {
  try {
    const response = await axiosInstance.get('/awali-power-generation');
    return response.data;
  } catch (error) {
    console.error('Error fetching Awali power generation data:', error);
    throw error;
  }
};

export const fetchRefineryPowerGenerationData = async (): Promise<PowerGenerationData> => {
  try {
    const response = await axiosInstance.get('/refinery-power-generation');
    return response.data;
  } catch (error) {
    console.error('Error fetching Refinery power generation data:', error);
    throw error;
  }
};