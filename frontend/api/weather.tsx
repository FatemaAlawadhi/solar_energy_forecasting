import axios from 'axios';

interface WeatherImpactData {
  year: number;
  month: number;
  avgSunshineDuration: number;
  avgDaylightDuration: number;
  minTemperature: number;
  avgTemperature: number;
  maxTemperature: number;
  avgSolarIrradiance: number;
  avgRelativeHumidity: number;
  avgCloudCover: number;
  avgWindSpeed: number;
  cumulativeRainfall: number;
  totalPowerGeneration: number;
}

interface FeatureImportance {
  featureName: string;
  importanceValue: number;
}

interface WeatherResponse {
  weatherData: WeatherImpactData[];
  featureImportance: FeatureImportance[];
}

const API_URL = 'http://localhost:8080/api';

export const fetchWeatherData = async (): Promise<WeatherResponse> => {
  try {
    const response = await axios.get(`${API_URL}/weather-impact`);
    return response.data;
  } catch (error) {
    console.error('Error fetching weather data:', error);
    throw error;
  }
};