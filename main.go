// package main

// import (
// 	"crypto/tls"
// 	"encoding/base64"
// 	"fmt"
// 	"log"
// 	"net/url"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/go-resty/resty/v2"
// 	"github.com/joho/godotenv"
// )

// // SplunkConfig holds the Splunk instance information.
// type SplunkConfig struct {
// 	BaseURL  string
// 	Username string
// 	Password string
// 	ApiPath  string // New field for API path
// }

// // DashboardConfig holds the dashboard configuration data.
// type DashboardConfig struct {
// 	TeamName string
// 	Index    string
// }

// // getAuthToken gets a session token from Splunk
// func getAuthToken(client *resty.Client, splunk SplunkConfig) (string, error) {
// 	loginURL := fmt.Sprintf("%s/services/auth/login", splunk.BaseURL)

// 	formData := url.Values{}
// 	formData.Set("username", splunk.Username)
// 	formData.Set("password", splunk.Password)
// 	formData.Set("output_mode", "json")

// 	resp, err := client.R().
// 		SetHeader("Content-Type", "application/x-www-form-urlencoded").
// 		SetBody(formData.Encode()).
// 		Post(loginURL)

// 	if err != nil {
// 		return "", fmt.Errorf("authentication request failed: %v", err)
// 	}

// 	if resp.StatusCode() != 200 {
// 		return "", fmt.Errorf("authentication failed: %s - %s", resp.Status(), string(resp.Body()))
// 	}

// 	log.Println("Session token received from splunk")

// 	return base64.StdEncoding.EncodeToString([]byte(splunk.Username + ":" + splunk.Password)), nil
// }

// // checkDashboardExists checks if the dashboard already exists
// func checkDashboardExists(client *resty.Client, splunk SplunkConfig, dashboardName string) (bool, error) {
// 	getDashboardURL := fmt.Sprintf("%s%s/%s",
// 		strings.TrimSuffix(splunk.BaseURL, "/"),
// 		splunk.ApiPath,
// 		url.PathEscape(fmt.Sprintf("dashboard_%s", strings.ToLower(dashboardName))))

// 	getDashboardURL += "?output_mode=json"

// 	resp, err := client.R().
// 		SetHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(splunk.Username+":"+splunk.Password))).
// 		Get(getDashboardURL)

// 	if err != nil {
// 		return false, fmt.Errorf("error checking dashboard existence: %v", err)
// 	}

// 	return resp.StatusCode() == 200, nil
// }

// // loadDashboardTemplate reads the dashboard XML template from file
// func loadDashboardTemplate(filepath string) (string, error) {
// 	content, err := os.ReadFile(filepath)
// 	if err != nil {
// 		return "", fmt.Errorf("error reading dashboard template: %v", err)
// 	}
// 	return string(content), nil
// }

// // createOrUpdateDashboard handles both creation and updating of dashboards
// func createOrUpdateDashboard(splunk SplunkConfig, dashboard DashboardConfig) error {
// 	// client := resty.New().
// 	// 	SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
// 	// 	SetTimeout(30 * time.Second).
// 	// 	SetDebug(true)

// 	client := resty.New().
// 		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
// 		SetTimeout(30 * time.Second)

// 	// Get authentication token
// 	authToken, err := getAuthToken(client, splunk)
// 	if err != nil {
// 		return fmt.Errorf("authentication failed: %v", err)
// 	}

// 	// Check if dashboard exists
// 	exists, err := checkDashboardExists(client, splunk, dashboard.TeamName)
// 	if err != nil {
// 		return fmt.Errorf("error checking dashboard existence: %v", err)
// 	}

// 	template, err := loadDashboardTemplate("dashboard_template.xml")
// 	if err != nil {
// 		return fmt.Errorf("failed to load dashboard template: %v", err)
// 	}

// 	dashboardXML := fmt.Sprintf(template, dashboard.TeamName)

// 	// Prepare form data
// 	formData := url.Values{}
// 	formData.Set("name", fmt.Sprintf("dashboard_%s", strings.ToLower(dashboard.TeamName)))
// 	formData.Set("eai:data", dashboardXML)
// 	formData.Set("output_mode", "json")

// 	dashboardName := fmt.Sprintf("dashboard_%s", strings.ToLower(dashboard.TeamName))

// 	if exists {
// 		log.Println("Updating existing dashboard...")
// 		err = updateDashboard(client, splunk, dashboardName, dashboardXML, authToken)
// 		if err != nil {
// 			return fmt.Errorf("error updating dashboard: %v", err)
// 		}
// 		log.Println("Dashboard updated successfully")
// 	} else {
// 		log.Println("Creating new dashboard...")
// 		err = createDashboard(client, splunk, dashboardName, dashboardXML, authToken)
// 		if err != nil {
// 			return fmt.Errorf("error creating dashboard: %v", err)
// 		}
// 		log.Println("Dashboard created successfully")
// 	}

// 	return setDashboardPermissions(client, splunk, dashboardName, authToken)
// }

// // updateDashboard updates an existing Splunk dashboard
// func updateDashboard(client *resty.Client, splunk SplunkConfig, dashboardName, dashboardXML, authToken string) error {
// 	apiURL := fmt.Sprintf("%s%s/%s",
// 		strings.TrimSuffix(splunk.BaseURL, "/"),
// 		splunk.ApiPath,
// 		url.PathEscape(dashboardName))

// 	formData := url.Values{}
// 	// formData.Set("name", dashboardName)
// 	formData.Set("eai:data", dashboardXML)
// 	formData.Set("output_mode", "json")

// 	resp, err := client.R().
// 		SetHeader("Authorization", "Basic "+authToken).
// 		SetHeader("Content-Type", "application/x-www-form-urlencoded").
// 		SetBody(formData.Encode()).
// 		Post(apiURL)

// 	if err != nil {
// 		return fmt.Errorf("error making request to update dashboard: %v", err)
// 	}

// 	if resp.StatusCode() != 200 {
// 		return fmt.Errorf("failed to update dashboard: %s - %s", resp.Status(), string(resp.Body()))
// 	}

// 	return nil
// }

// // createDashboard creates a new Splunk dashboard
// func createDashboard(client *resty.Client, splunk SplunkConfig, dashboardName, dashboardXML, authToken string) error {
// 	apiURL := fmt.Sprintf("%s%s",
// 		strings.TrimSuffix(splunk.BaseURL, "/"),
// 		splunk.ApiPath)

// 	formData := url.Values{}
// 	formData.Set("name", dashboardName)
// 	formData.Set("eai:data", dashboardXML)
// 	formData.Set("output_mode", "json")

// 	resp, err := client.R().
// 		SetHeader("Authorization", "Basic "+authToken).
// 		SetHeader("Content-Type", "application/x-www-form-urlencoded").
// 		SetBody(formData.Encode()).
// 		Post(apiURL)

// 	if err != nil {
// 		return fmt.Errorf("error making request to create dashboard: %v", err)
// 	}

// 	if resp.StatusCode() != 201 {
// 		return fmt.Errorf("failed to create dashboard: %s - %s", resp.Status(), string(resp.Body()))
// 	}

// 	return nil
// }

// // setDashboardPermissions sets read and write permissions for the dashboard
// func setDashboardPermissions(client *resty.Client, splunk SplunkConfig, dashboardName, authToken string) error {
// 	permissionsURL := fmt.Sprintf("%s%s/%s/acl",
// 		strings.TrimSuffix(splunk.BaseURL, "/"),
// 		splunk.ApiPath,
// 		url.PathEscape(dashboardName))

// 	formData := url.Values{}
// 	formData.Set("sharing", "app")
// 	formData.Set("owner", "admin")
// 	formData.Set("perms.read", "*")
// 	formData.Set("perms.write", "admin")

// 	resp, err := client.R().
// 		SetHeader("Authorization", "Basic "+authToken).
// 		SetHeader("Content-Type", "application/x-www-form-urlencoded").
// 		SetBody(formData.Encode()).
// 		Post(permissionsURL)

// 	if err != nil {
// 		return fmt.Errorf("error setting dashboard permissions: %v", err)
// 	}

// 	if resp.StatusCode() != 200 {
// 		return fmt.Errorf("failed to set dashboard permissions: %s - %s", resp.Status(), string(resp.Body()))
// 	}

// 	return nil
// }

// func main() {
// 	// Load environment variables from .env file
// 	if err := godotenv.Load(); err != nil {
// 		log.Printf("Warning: Error loading .env file: %v\n", err)
// 	}

// 	// Get credentials with fallback to default values
// 	username := getEnvWithDefault("SPLUNK_USERNAME", "admin")
// 	password := getEnvWithDefault("SPLUNK_PASSWORD", "")
// 	baseURL := getEnvWithDefault("SPLUNK_BASE_URL", "https://localhost:8089")
// 	apiPath := getEnvWithDefault("SPLUNK_API_PATH", "/servicesNS/admin/search/data/ui/views") // New environment variable

// 	if password == "" {
// 		log.Println("Error: Splunk password not set in environment variables")
// 		return
// 	}

// 	splunk := SplunkConfig{
// 		BaseURL:  baseURL,
// 		Username: username,
// 		Password: password,
// 		ApiPath:  apiPath,
// 	}

// 	dashboard := DashboardConfig{
// 		TeamName: "Team1",
// 		Index:    "user_management_api_dev",
// 	}

// 	if err := createOrUpdateDashboard(splunk, dashboard); err != nil {
// 		log.Printf("Error creating/updating dashboard: %v\n", err)
// 		return
// 	}
// }

// // getEnvWithDefault returns environment variable value or default if not set
// func getEnvWithDefault(key, defaultValue string) string {
// 	if value := os.Getenv(key); value != "" {
// 		return value
// 	}
// 	return defaultValue
// }

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
	ApiPath string // New field for API path
}

// DashboardConfig holds the dashboard configuration data.
type DashboardConfig struct {
	TeamName string
	Index    string
}

// checkDashboardExists checks if the dashboard already exists
func checkDashboardExists(client *resty.Client, splunk SplunkConfig, dashboardName string) (bool, error) {
	getDashboardURL := fmt.Sprintf("%s%s/%s",
		strings.TrimSuffix(splunk.BaseURL, "/"),
		splunk.ApiPath,
		url.PathEscape(fmt.Sprintf("dashboard_%s", strings.ToLower(dashboardName))))

	getDashboardURL += "?output_mode=json"

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+splunk.Token). // Use token-based authorization
		Get(getDashboardURL)

	if err != nil {
		return false, fmt.Errorf("error checking dashboard existence: %v", err)
	}

	return resp.StatusCode() == 200, nil
}

// loadDashboardTemplate reads the dashboard XML template from file
func loadDashboardTemplate(filepath string) (string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("error reading dashboard template: %v", err)
	}
	return string(content), nil
}

// createOrUpdateDashboard handles both creation and updating of dashboards
func createOrUpdateDashboard(splunk SplunkConfig, dashboard DashboardConfig) error {
	client := resty.New().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetTimeout(30 * time.Second)

	// Check if dashboard exists
	exists, err := checkDashboardExists(client, splunk, dashboard.TeamName)
	if err != nil {
		return fmt.Errorf("error checking dashboard existence: %v", err)
	}

	template, err := loadDashboardTemplate("dashboard_template.xml")
	if err != nil {
		return fmt.Errorf("failed to load dashboard template: %v", err)
	}

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

// updateDashboard updates an existing Splunk dashboard
func updateDashboard(client *resty.Client, splunk SplunkConfig, dashboardName, dashboardXML string) error {
	apiURL := fmt.Sprintf("%s%s/%s",
		strings.TrimSuffix(splunk.BaseURL, "/"),
		splunk.ApiPath,
		url.PathEscape(dashboardName))

	formData := url.Values{}
	formData.Set("eai:data", dashboardXML)
	formData.Set("output_mode", "json")

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+splunk.Token). // Use token-based authorization
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

// createDashboard creates a new Splunk dashboard
func createDashboard(client *resty.Client, splunk SplunkConfig, dashboardName, dashboardXML string) error {
	apiURL := fmt.Sprintf("%s%s",
		strings.TrimSuffix(splunk.BaseURL, "/"),
		splunk.ApiPath)

	formData := url.Values{}
	formData.Set("name", dashboardName)
	formData.Set("eai:data", dashboardXML)
	formData.Set("output_mode", "json")

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+splunk.Token). // Use token-based authorization
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

// setDashboardPermissions sets read and write permissions for the dashboard
func setDashboardPermissions(client *resty.Client, splunk SplunkConfig, dashboardName string) error {
	permissionsURL := fmt.Sprintf("%s%s/%s/acl",
		strings.TrimSuffix(splunk.BaseURL, "/"),
		splunk.ApiPath,
		url.PathEscape(dashboardName))

	formData := url.Values{}
	formData.Set("sharing", "app")
	formData.Set("owner", "admin")
	formData.Set("perms.read", "*")
	formData.Set("perms.write", "admin")

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+splunk.Token). // Use token-based authorization
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
	token := os.Getenv("SPLUNK_TOKEN") // Using token from environment variables
	apiPath := os.Getenv("SPLUNK_API_PATH")

	if token == "" {
		log.Println("Error: Splunk token not set in environment variables")
		return
	}

	// splunk := SplunkConfig{
	// 	BaseURL: baseURL,
	// 	Token:   token,
	// 	ApiPath: apiPath,
	// }

	splunk := SplunkConfig{
		BaseURL: baseURL,
		Token:   cfg.SplunkToken,
		ApiPath: apiPath,
	}

	dashboard := DashboardConfig{
		TeamName: "Team1",
		Index:    "user_management_api_dev",
	}

	if err := createOrUpdateDashboard(splunk, dashboard); err != nil {
		log.Printf("Error creating/updating dashboard: %v\n", err)
		return
	}
}
