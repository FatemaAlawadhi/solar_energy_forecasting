import { axiosInstance } from './config';

interface Generation {
  year: number;
  month: number;
  location_id: number;
  location_name: string;
  actual_kwh: number;
  theoretical_kwh: number;
}

interface Performance {
  year?: number;
  month?: number;
  start_year?: number;
  end_year?: number;
  location_id: number;
  location_name: string;
  performance_ratio: number;
  capacity_factor: number;
  output_per_pv: number;
}

interface PerformanceResponse {
  monthly_generation: Generation[];
  monthly_performance: Performance[];
  yearly_performance: Performance[];
  overall_performance: Performance[];
}

export const fetchPerformanceData = async (): Promise<PerformanceResponse> => {
  try {
    const response = await axiosInstance.get('/performance');
    return response.data;
  } catch (error) {
    console.error('Error fetching performance data:', error);
    throw error;
  }
};
