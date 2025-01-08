package main

import (
	"dashboard/config"
	"dashboard/models"
	"dashboard/utils/dashboard"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v\n", err)
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
		log.Fatalf("Dashboard setup failed: %v", err)
	}
}
