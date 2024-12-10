import sqlite3
import pandas as pd
import numpy as np
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler
from sklearn.metrics import r2_score
from sklearn.linear_model import LinearRegression
from sklearn.ensemble import RandomForestRegressor, AdaBoostRegressor, GradientBoostingRegressor
from sklearn.neural_network import MLPRegressor
from sklearn.neighbors import KNeighborsRegressor
from sklearn.svm import SVR
from sklearn.tree import DecisionTreeRegressor
import lightgbm as lgb
import xgboost as xgb
from pygam import GAM
from sklearn.multioutput import MultiOutputRegressor
import matplotlib.pyplot as plt
import seaborn as sns
import os
import matplotlib
matplotlib.use('TkAgg')

def load_and_prepare_data():
    # Connect to SQLite database
    db_path = os.path.join(os.path.dirname(__file__), '..', '..', 'db', 'app.db')
    conn = sqlite3.connect(db_path)
    
    # Load monthly weather data (features)
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
            mg.year, 
            mg.month,
            MAX(CASE WHEN l.name = 'Refinery' THEN mg.actual_kwh END) as total_refinery,
            MAX(CASE WHEN l.name = 'Awali' THEN mg.actual_kwh END) as total_awali,
            MAX(CASE WHEN l.name = 'UOB' THEN mg.actual_kwh END) as total_UOB,
            MAX(CASE WHEN l.name = 'Total System' THEN mg.actual_kwh END) as total_system
        FROM monthly_generation mg
        JOIN locations l ON mg.location_id = l.id
        GROUP BY mg.year, mg.month
        ORDER BY mg.year, mg.month
    """, conn)
    
    conn.close()
    
    # Merge weather and power generation data
    merged_data = weather_data.merge(power_data, on=['year', 'month'])
    
    # Separate features and targets
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
    
    y = merged_data[['total_awali', 'total_refinery', 'total_UOB', 'total_system']]
    
    return X, y

def evaluate_models(X_train, X_test, y_train, y_test):
    models = {
        'Linear Regression': LinearRegression(),
        'Random Forest': RandomForestRegressor(),
        'Neural Network': MLPRegressor(),
        'KNN': KNeighborsRegressor(),
        'SVR': MultiOutputRegressor(SVR()),
        'AdaBoost': MultiOutputRegressor(AdaBoostRegressor()),
        'Decision Tree': DecisionTreeRegressor(),
        'LGBM': MultiOutputRegressor(lgb.LGBMRegressor()),
        'XGBoost': MultiOutputRegressor(xgb.XGBRegressor()),
        'Gradient Boosting': MultiOutputRegressor(GradientBoostingRegressor()),
        'GAM': MultiOutputRegressor(GAM())
    }
    
    results = {}
    
    for name, model in models.items():
        print(f"Training {name}...")
        model.fit(X_train, y_train)
        y_pred = model.predict(X_test)
        
        r2_scores = {}
        for i, col in enumerate(y_test.columns):
            r2 = r2_score(y_test.iloc[:, i], y_pred[:, i])
            r2_scores[col] = r2
        
        results[name] = r2_scores
    
    return results

def find_best_model(results):
    model_avg_scores = {}
    for model_name, scores in results.items():
        avg_score = sum(scores.values()) / len(scores)
        model_avg_scores[model_name] = avg_score
    
    best_model = max(model_avg_scores.items(), key=lambda x: x[1])
    
    print("\nModels Ranked by Average R² Score:")
    print("=" * 50)
    for model_name, avg_score in sorted(model_avg_scores.items(), key=lambda x: x[1], reverse=True):
        print(f"{model_name}: {avg_score:.4f}")
    
    print("\nBest Performing Model:")
    print("=" * 50)
    print(f"Model: {best_model[0]}")
    print(f"Average R² Score: {best_model[1]:.4f}")
    print("\nDetailed R² Scores for Best Model:")
    for target, score in results[best_model[0]].items():
        print(f"{target}: {score:.4f}")
    
    return best_model[0]

def visualize_results(results):
    plots_folder = 'backend/pkg/model/monthly/plots'
    os.makedirs(plots_folder, exist_ok=True)
    
    # Define metrics and their display names
    metrics = ['total_awali', 'total_refinery', 'total_UOB', 'total_system']
    location_names = ['Awali', 'Refinery', 'UOB', 'Total System']
    
    plt.style.use('default')
    fig, axes = plt.subplots(len(metrics), 2, figsize=(15, 4 * len(metrics)), facecolor='white')
    fig.suptitle('Monthly Model Performance Comparison', fontsize=14, y=1.02)
    
    for i, (metric, location) in enumerate(zip(metrics, location_names)):
        metric_scores = [results[model][metric] for model in results.keys()]
        model_names = list(results.keys())
        
        # Split into positive and negative scores
        pos_scores = [score if score >= 0 else 0 for score in metric_scores]
        neg_scores = [score if score < 0 else 0 for score in metric_scores]
        
        # Add location title centered between the plots
        axes[i, 0].text(1.0, 1.15, location, 
                       horizontalalignment='center',
                       transform=axes[i, 0].transAxes,
                       fontsize=12,
                       fontweight='bold')
        
        # Plot negative scores (left plot)
        axes[i, 0].barh(model_names, neg_scores)
        axes[i, 0].set_title('Negative R² Scores')
        axes[i, 0].set_xlabel('R² Score')
        axes[i, 0].grid(True, alpha=0.3)
        axes[i, 0].set_xlim(min(neg_scores) * 1.1, 0)  # Add 10% margin
        
        # Plot positive scores (right plot)
        axes[i, 1].barh(model_names, pos_scores)
        axes[i, 1].set_title('Positive R² Scores')
        axes[i, 1].set_xlabel('R² Score')
        axes[i, 1].grid(True, alpha=0.3)
        axes[i, 1].set_xlim(0, max(pos_scores) * 1.1)  # Add 10% margin
    
    plt.subplots_adjust(
        top=0.95,
        bottom=0.05,
        left=0.3,
        right=0.95,
        hspace=0.6,  
        wspace=0.4
    )
    
    output_path = os.path.join(plots_folder, 'monthly_model_comparison.png')
    plt.savefig(output_path, 
                bbox_inches='tight', 
                dpi=150,
                facecolor='white',
                edgecolor='none',
                transparent=False)
    plt.close()

    print(f"\nPlot saved to: {output_path}")

def main():
    # Load data
    X, y = load_and_prepare_data()
    
    # Scale features
    scaler = StandardScaler()
    X_scaled = scaler.fit_transform(X)
    X_scaled = pd.DataFrame(X_scaled, columns=X.columns)
    
    # Split data
    X_train, X_test, y_train, y_test = train_test_split(
        X_scaled, y, test_size=0.2, random_state=42, shuffle=False  # shuffle=False to maintain time order
    )
    
    # Evaluate models
    results = evaluate_models(X_train, X_test, y_train, y_test)
    
    # Print all results
    print("\nAll Models Performance (R² scores):")
    print("=" * 50)
    for model_name, scores in results.items():
        print(f"\n{model_name}:")
        for target, r2 in scores.items():
            print(f"{target}: {r2:.4f}")
    
    # Find and print the best model
    best_model = find_best_model(results)
    
    # Visualize results
    visualize_results(results)

if __name__ == "__main__":
    main()