import sqlite3
import pandas as pd
import numpy as np
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler
from sklearn.ensemble import RandomForestRegressor
from sklearn.metrics import r2_score, mean_squared_error
import matplotlib.pyplot as plt
import os
import warnings
warnings.filterwarnings('ignore')

class WeatherOnlyModel:
    def __init__(self):
        self.model = RandomForestRegressor(
             n_estimators=500,          
            max_depth=15,              
            min_samples_split=4,       
            min_samples_leaf=2,        
            max_features=0.8,          
            random_state=42,
            n_jobs=-1,                 
            bootstrap=True,
            min_impurity_decrease=0.0001  
        )
        self.scaler = StandardScaler()
    
    def load_and_prepare_data(self):
        # Connect to SQLite database
        db_path = os.path.join(os.path.dirname(__file__), '..', '..', 'db', 'app.db')
        conn = sqlite3.connect(db_path)
        
        # Load only weather data
        weather_data = pd.read_sql_query("""
            SELECT 
                year,
                month,
                avg_sunshine_duration_seconds,
                avg_daylight_duration_seconds,
                min_temperature_C,
                avg_temperature_C,
                max_temperature_C,
                avg_solar_irradiance_wm2,
                avg_relative_humidity_percent,
                avg_cloud_cover_percent,
                avg_wind_speed_kmh,
                total_rainfall_mm
            FROM weather_monthly
            ORDER BY year, month
        """, conn)
        
        # Load power generation data 
        power_data = pd.read_sql_query("""
            SELECT 
                year, 
                month,
                MAX(CASE WHEN location_id = (SELECT id FROM locations WHERE name = 'Refinery') THEN actual_kwh END) as total_refinery,
                MAX(CASE WHEN location_id = (SELECT id FROM locations WHERE name = 'Awali') THEN actual_kwh END) as total_awali,
                MAX(CASE WHEN location_id = (SELECT id FROM locations WHERE name = 'UOB') THEN actual_kwh END) as total_UOB,
                MAX(CASE WHEN location_id = (SELECT id FROM locations WHERE name = 'Total System') THEN actual_kwh END) as total_all
            FROM monthly_generation
            GROUP BY year, month
            ORDER BY year, month
        """, conn)
        
        conn.close()
        
        # Merge data
        merged_data = weather_data.merge(power_data, on=['year', 'month'])
        
        # Store dates for plotting
        self.dates = pd.to_datetime(merged_data[['year', 'month']].assign(day=1))
        
        # Weather features only
        X = merged_data[[
            'avg_sunshine_duration_seconds',
            'avg_daylight_duration_seconds',
            'min_temperature_C',
            'avg_temperature_C',
            'max_temperature_C',
            'avg_solar_irradiance_wm2',
            'avg_relative_humidity_percent',
            'avg_cloud_cover_percent',
            'avg_wind_speed_kmh',
            'total_rainfall_mm'
        ]]
        
        y = merged_data[['total_awali', 'total_refinery', 'total_UOB', 'total_all']]
        
        return X, y
    
    def plot_feature_importance(self, feature_names):
        # Create plots directory if it doesn't exist
        current_dir = os.path.dirname(os.path.abspath(__file__))
        plots_folder = os.path.join(current_dir, "plots")
        os.makedirs(plots_folder, exist_ok=True)

        # Create readable feature names mapping
        feature_name_mapping = {
            'avg_sunshine_duration_seconds': 'Sunshine Duration',
            'avg_daylight_duration_seconds': 'Daylight Duration',
            'min_temperature_C': 'Min Temperature',
            'avg_temperature_C': 'Avg Temperature',
            'max_temperature_C': 'Max Temperature',
            'avg_solar_irradiance_wm2': 'Solar Irradiance',
            'avg_relative_humidity_percent': 'Relative Humidity',
            'avg_cloud_cover_percent': 'Cloud Cover',
            'avg_wind_speed_kmh': 'Wind Speed',
            'total_rainfall_mm': 'Rainfall'
        }

        importance_dict = {
            feature: importance 
            for feature, importance in zip(feature_names, self.model.feature_importances_)
        }
        
        # Create DataFrame and sort
        importance_df = pd.DataFrame({
            'feature': importance_dict.keys(),
            'importance': importance_dict.values()
        }).sort_values('importance', ascending=True)
        
        # Map feature names to readable names
        importance_df['feature'] = importance_df['feature'].map(feature_name_mapping)
        
        # Create visualization
        plt.figure(figsize=(12, 8))
        bars = plt.barh(importance_df['feature'], importance_df['importance'])
        
        # Add value labels on the bars
        for bar in bars:
            width = bar.get_width()
            plt.text(width, bar.get_y() + bar.get_height()/2, 
                    f'{width:.3f}', 
                    ha='left', va='center', fontweight='bold')
        
        plt.title('Weather Feature Importance in Power Generation', pad=20)
        plt.xlabel('Relative Importance')
        plt.ylabel('Features')
        plt.grid(True, alpha=0.3)
        plt.tight_layout()
        
        # Save plot to plots folder
        plt.savefig(os.path.join(plots_folder, 'weather_feature_importance.png'), 
                    bbox_inches='tight', dpi=300)
        plt.close()

        # Save feature importance to database
        self.save_feature_importance_to_db(importance_dict)
        
        # Print importance scores
        print("\nWeather Feature Importance Scores:")
        for feature, importance in sorted(importance_dict.items(), key=lambda x: x[1], reverse=True):
            print(f"{feature}: {importance:.4f}")

    def save_feature_importance_to_db(self, importance_dict):
        db_path = os.path.join(os.path.dirname(__file__), '..', '..', 'db', 'app.db')
        conn = sqlite3.connect(db_path)
        cursor = conn.cursor()
        
        # Create readable feature names mapping for database
        feature_name_mapping = {
            'avg_sunshine_duration_seconds': 'Sunshine Duration (hours)',
            'avg_daylight_duration_seconds': 'Daylight Duration (hours)',
            'min_temperature_C': 'Minimum Temperature (°C)',
            'avg_temperature_C': 'Average Temperature (°C)',
            'max_temperature_C': 'Maximum Temperature (°C)',
            'avg_solar_irradiance_wm2': 'Solar Irradiance (W/m²)',
            'avg_relative_humidity_percent': 'Relative Humidity (%)',
            'avg_cloud_cover_percent': 'Cloud Cover (%)',
            'avg_wind_speed_kmh': 'Wind Speed (km/h)',
            'total_rainfall_mm': 'Rainfall (mm)'
        }
        
        try:
            # Clear existing entries
            cursor.execute("DELETE FROM feature_importance")
            
            for feature, importance in importance_dict.items():
                display_name = feature_name_mapping.get(feature, feature)
                rounded_importance = round(importance, 3)
                cursor.execute("""
                    INSERT INTO feature_importance (feature_name, importance_value)
                    VALUES (?, ?)
                """, (display_name, rounded_importance))
            
            conn.commit()
        except sqlite3.Error as e:
            print(f"Database error: {e}")
        finally:
            conn.close()
    
    def train(self):
        X, y = self.load_and_prepare_data()
        
        # Scale features
        X_scaled = self.scaler.fit_transform(X)
        X_scaled = pd.DataFrame(X_scaled, columns=X.columns)
        
        # Split the data
        X_train, X_test, y_train, y_test = train_test_split(
            X_scaled, y, test_size=0.2,
            random_state=42, shuffle=False
        )
        
        print("Training Weather-Only Random Forest Model...")
        self.model.fit(X_train, y_train)
        
        # Make predictions
        y_pred = self.model.predict(X_test)
        
        # Calculate and print metrics
        for i, location in enumerate(['Awali', 'Refinery', 'UOB', 'Total']):
            r2 = r2_score(y_test.iloc[:, i], y_pred[:, i])
            rmse = np.sqrt(mean_squared_error(y_test.iloc[:, i], y_pred[:, i]))
            mape = np.mean(np.abs((y_test.iloc[:, i] - y_pred[:, i]) / y_test.iloc[:, i])) * 100
            
            print(f"\n{location} Results:")
            print(f"R² Score: {r2:.4f}")
            print(f"RMSE: {rmse:.2f}")
            print(f"Error Percentage: {mape:.2f}%")
        
        self.plot_feature_importance(X.columns)
        return self.model

def main():
    model = WeatherOnlyModel()
    model.train()

if __name__ == "__main__":
    main()