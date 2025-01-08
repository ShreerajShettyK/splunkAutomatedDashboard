package main

import (
	"crypto/tls"
	"dashboard/config"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

// SplunkConfig holds the Splunk instance information.
type SplunkConfig struct {
	BaseURL string
	Token   string
	ApiPath string
}

// DashboardConfig holds the dashboard configuration data.
type DashboardConfig struct {
	TeamName string
	Index    string
}

func createOrUpdateDashboard(splunk SplunkConfig, dashboard DashboardConfig) error {
	client := resty.New().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetTimeout(30 * time.Second)

	exists, err := checkDashboardExists(client, splunk, dashboard.TeamName)
	if err != nil {
		return fmt.Errorf("error checking dashboard existence: %v", err)
	}

	template, err := loadDashboardTemplate("dashboard_template_json.xml")
	if err != nil {
		return fmt.Errorf("failed to load dashboard template: %v", err)
	}

	// Format the template with the team name
	dashboardXML := fmt.Sprintf(template, dashboard.TeamName)
	dashboardName := fmt.Sprintf("dashboard_%s", strings.ToLower(dashboard.TeamName))

	if exists {
		log.Println("Updating existing dashboard...")
		err = updateDashboard(client, splunk, dashboardName, dashboardXML)
		if err != nil {
			return fmt.Errorf("error updating dashboard: %v", err)
		}
		log.Println("Dashboard updated successfully")
	} else {
		log.Println("Creating new dashboard...")
		err = createDashboard(client, splunk, dashboardName, dashboardXML)
		if err != nil {
			return fmt.Errorf("error creating dashboard: %v", err)
		}
		log.Println("Dashboard created successfully")
	}

	return setDashboardPermissions(client, splunk, dashboardName)
}

func checkDashboardExists(client *resty.Client, splunk SplunkConfig, dashboardName string) (bool, error) {
	getDashboardURL := fmt.Sprintf("%s%s/%s",
		strings.TrimSuffix(splunk.BaseURL, "/"),
		splunk.ApiPath,
		url.PathEscape(fmt.Sprintf("dashboard_%s", strings.ToLower(dashboardName))))

	getDashboardURL += "?output_mode=json"

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+splunk.Token).
		Get(getDashboardURL)

	if err != nil {
		return false, fmt.Errorf("error checking dashboard existence: %v", err)
	}

	return resp.StatusCode() == 200, nil
}

func loadDashboardTemplate(filepath string) (string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("error reading dashboard template: %v", err)
	}
	return string(content), nil
}

func createDashboard(client *resty.Client, splunk SplunkConfig, dashboardName, dashboardXML string) error {
	apiURL := fmt.Sprintf("%s%s",
		strings.TrimSuffix(splunk.BaseURL, "/"),
		splunk.ApiPath)

	// Create form data similar to the curl command
	formData := url.Values{}
	formData.Set("name", dashboardName)
	formData.Set("eai:data", strings.TrimSpace(dashboardXML)) // Trim any extra whitespace
	formData.Set("output_mode", "json")

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+splunk.Token).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(formData.Encode()).
		Post(apiURL)

	if err != nil {
		return fmt.Errorf("error making request to create dashboard: %v", err)
	}

	if resp.StatusCode() != 201 {
		return fmt.Errorf("failed to create dashboard: %s - %s", resp.Status(), string(resp.Body()))
	}

	return nil
}

func updateDashboard(client *resty.Client, splunk SplunkConfig, dashboardName, dashboardXML string) error {
	apiURL := fmt.Sprintf("%s%s/%s",
		strings.TrimSuffix(splunk.BaseURL, "/"),
		splunk.ApiPath,
		url.PathEscape(dashboardName))

	formData := url.Values{}
	formData.Set("eai:data", strings.TrimSpace(dashboardXML)) // Trim any extra whitespace
	formData.Set("output_mode", "json")

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+splunk.Token).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(formData.Encode()).
		Post(apiURL)

	if err != nil {
		return fmt.Errorf("error making request to update dashboard: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to update dashboard: %s - %s", resp.Status(), string(resp.Body()))
	}

	return nil
}

func setDashboardPermissions(client *resty.Client, splunk SplunkConfig, dashboardName string) error {
	permissionsURL := fmt.Sprintf("%s%s/%s/acl",
		strings.TrimSuffix(splunk.BaseURL, "/"),
		splunk.ApiPath,
		url.PathEscape(dashboardName))

	formData := url.Values{}
	formData.Set("sharing", "app")
	formData.Set("owner", "admin")
	formData.Set("perms.read", "splunkviewingrole")
	formData.Set("perms.write", "admin")

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+splunk.Token).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(formData.Encode()).
		Post(permissionsURL)

	if err != nil {
		return fmt.Errorf("error setting dashboard permissions: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to set dashboard permissions: %s - %s", resp.Status(), string(resp.Body()))
	}

	log.Println("Dashboard permissions set")
	return nil
}

func main() {
	// Load the Splunk token from Secrets Manager
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v\n", err)
	}

	baseURL := os.Getenv("SPLUNK_BASE_URL")
	apiPath := os.Getenv("SPLUNK_API_PATH")

	splunk := SplunkConfig{
		BaseURL: baseURL,
		Token:   cfg.SplunkToken,
		ApiPath: apiPath,
	}

	dashboard := DashboardConfig{
		TeamName: "Splunk",
		Index:    "user_management_api_dev",
	}

	if err := createOrUpdateDashboard(splunk, dashboard); err != nil {
		log.Printf("Error creating/updating dashboard: %v\n", err)
		return
	}
}
