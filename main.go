package main

import (
	"dashboard/config"
	"dashboard/logger"
	"dashboard/models"
	"dashboard/utils/dashboard"
	"os"

	"github.com/joho/godotenv"
)

var splunkLogger = logger.CreateLogger()

func main() {
	splunkLogger.Println("Starting Splunk Dashboard Setup")
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		splunkLogger.Fatalf("Failed to load config: %v", err)
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		splunkLogger.Printf("Warning: Error loading .env file: %v\n", err)
	}

	baseURL := os.Getenv("SPLUNK_BASE_URL")
	apiPath := os.Getenv("SPLUNK_API_PATH")

	// Initialize Splunk and Dashboard Config
	splunkConfig := &models.SplunkConfig{
		BaseURL: baseURL,
		Token:   cfg.SplunkToken,
		ApiPath: apiPath,
	}

	dashboardConfig := &models.DashboardConfig{
		TeamName: "Splunk",
		Index:    "user_management_api_dev",
	}

	// Run dashboard setup
	if err := dashboard.RunDashboardSetup(splunkConfig, dashboardConfig); err != nil {
		splunkLogger.Fatalf("Dashboard setup failed: %v", err)
	}
}
