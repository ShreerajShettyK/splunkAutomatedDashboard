// package main

// import (
// 	"crypto/tls"
// 	"encoding/base64"
// 	"fmt"
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
// }

// // DashboardConfig holds the dashboard configuration data.
// type DashboardConfig struct {
// 	TeamName   string
// 	Index      string
// 	PanelQuery string
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

// 	return base64.StdEncoding.EncodeToString([]byte(splunk.Username + ":" + splunk.Password)), nil
// }

// // checkDashboardExists checks if the dashboard already exists
// func checkDashboardExists(client *resty.Client, splunk SplunkConfig, dashboardName string) (bool, error) {
// 	// Construct the URL to check for the dashboard
// 	getDashboardURL := fmt.Sprintf("%s/servicesNS/admin/search/data/ui/views/%s",
// 		strings.TrimSuffix(splunk.BaseURL, "/"),
// 		url.PathEscape(fmt.Sprintf("dashboard_%s", strings.ToLower(dashboardName))))

// 	// Add output_mode=json to get JSON response
// 	getDashboardURL += "?output_mode=json"

// 	resp, err := client.R().
// 		SetHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(splunk.Username+":"+splunk.Password))).
// 		Get(getDashboardURL)

// 	if err != nil {
// 		return false, fmt.Errorf("error checking dashboard existence: %v", err)
// 	}

// 	// For debugging
// 	fmt.Printf("Check dashboard response status: %s\n", resp.Status())
// 	fmt.Printf("Check dashboard response body: %s\n", string(resp.Body()))

// 	return resp.StatusCode() == 200, nil
// }

// // createOrUpdateDashboard handles both creation and updating of dashboards
// func createOrUpdateDashboard(splunk SplunkConfig, dashboard DashboardConfig) error {
// 	client := resty.New().
// 		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
// 		SetTimeout(30 * time.Second).
// 		SetDebug(true)

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

// 	// Create dashboard XML
// 	dashboardXML := fmt.Sprintf(`<?xml version="1.0"?>
// 		<dashboard version="1.1">
// 			<label>%s Dashboard</label>
// 			<row>
// 				<panel>
// 					<title>Log Events</title>
// 					<chart>
// 						<search>
// 							<query>index=%s | stats count by source</query>
// 							<earliest>-24h@h</earliest>
// 							<latest>now</latest>
// 						</search>
// 						<option name="charting.chart">column</option>
// 					</chart>
// 				</panel>
// 			</row>
// 		</dashboard>`, dashboard.TeamName, dashboard.Index)

// 	// Prepare form data
// 	formData := url.Values{}
// 	formData.Set("name", fmt.Sprintf("dashboard_%s", strings.ToLower(dashboard.TeamName)))
// 	formData.Set("eai:data", dashboardXML)
// 	formData.Set("output_mode", "json")

// 	var apiURL string
// 	if exists {
// 		// Update existing dashboard
// 		apiURL = fmt.Sprintf("%s/servicesNS/admin/search/data/ui/views/%s",
// 			strings.TrimSuffix(splunk.BaseURL, "/"),
// 			url.PathEscape(fmt.Sprintf("dashboard_%s", strings.ToLower(dashboard.TeamName))))
// 		fmt.Println("Updating existing dashboard...")
// 	} else {
// 		// Create new dashboard
// 		apiURL = fmt.Sprintf("%s/servicesNS/admin/search/data/ui/views",
// 			strings.TrimSuffix(splunk.BaseURL, "/"))
// 		fmt.Println("Creating new dashboard...")
// 	}

// 	// Send request to create/update the dashboard
// 	resp, err := client.R().
// 		SetHeader("Authorization", "Basic "+authToken).
// 		SetHeader("Content-Type", "application/x-www-form-urlencoded").
// 		SetBody(formData.Encode()).
// 		Post(apiURL)

// 	if err != nil {
// 		return fmt.Errorf("error making request: %v", err)
// 	}

// 	// Print detailed response for debugging
// 	fmt.Printf("Response Status: %s\n", resp.Status())
// 	fmt.Printf("Response Body: %s\n", string(resp.Body()))

// 	if resp.StatusCode() != 200 && resp.StatusCode() != 201 {
// 		return fmt.Errorf("failed to %s dashboard: %s - %s",
// 			map[bool]string{true: "update", false: "create"}[exists],
// 			resp.Status(), string(resp.Body()))
// 	}

// 	fmt.Printf("Dashboard successfully %s for %s\n",
// 		map[bool]string{true: "updated", false: "created"}[exists],
// 		dashboard.TeamName)
// 	return nil
// }

// func main() {
// 	// Load environment variables from .env file
// 	if err := godotenv.Load(); err != nil {
// 		fmt.Printf("Warning: Error loading .env file: %v\n", err)
// 	}

// 	// Get credentials with fallback to default values
// 	username := getEnvWithDefault("SPLUNK_USERNAME", "admin")
// 	password := getEnvWithDefault("SPLUNK_PASSWORD", "")
// 	baseURL := getEnvWithDefault("SPLUNK_BASE_URL", "https://localhost:8089")

// 	if password == "" {
// 		fmt.Println("Error: Splunk password not set in environment variables")
// 		return
// 	}

// 	splunk := SplunkConfig{
// 		BaseURL:  baseURL,
// 		Username: username,
// 		Password: password,
// 	}

// 	dashboard := DashboardConfig{
// 		TeamName:   "Team1",
// 		Index:      "user_management_api_dev",
// 		PanelQuery: "index=user_management_api_dev | stats count by source",
// 	}

// 	if err := createOrUpdateDashboard(splunk, dashboard); err != nil {
// 		fmt.Printf("Error creating/updating dashboard: %v\n", err)
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
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

// SplunkConfig holds the Splunk instance information.
type SplunkConfig struct {
	BaseURL  string
	Username string
	Password string
}

// DashboardConfig holds the dashboard configuration data.
type DashboardConfig struct {
	TeamName   string
	Index      string
	PanelQuery string
}

// getAuthToken gets a session token from Splunk
func getAuthToken(client *resty.Client, splunk SplunkConfig) (string, error) {
	loginURL := fmt.Sprintf("%s/services/auth/login", splunk.BaseURL)

	formData := url.Values{}
	formData.Set("username", splunk.Username)
	formData.Set("password", splunk.Password)
	formData.Set("output_mode", "json")

	resp, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(formData.Encode()).
		Post(loginURL)

	if err != nil {
		return "", fmt.Errorf("authentication request failed: %v", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("authentication failed: %s - %s", resp.Status(), string(resp.Body()))
	}

	return base64.StdEncoding.EncodeToString([]byte(splunk.Username + ":" + splunk.Password)), nil
}

// checkDashboardExists checks if the dashboard already exists
func checkDashboardExists(client *resty.Client, splunk SplunkConfig, dashboardName string) (bool, error) {
	getDashboardURL := fmt.Sprintf("%s/servicesNS/admin/search/data/ui/views/%s",
		strings.TrimSuffix(splunk.BaseURL, "/"),
		url.PathEscape(fmt.Sprintf("dashboard_%s", strings.ToLower(dashboardName))))

	getDashboardURL += "?output_mode=json"

	resp, err := client.R().
		SetHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(splunk.Username+":"+splunk.Password))).
		Get(getDashboardURL)

	if err != nil {
		return false, fmt.Errorf("error checking dashboard existence: %v", err)
	}

	// For debugging
	// fmt.Printf("Check dashboard response status: %s\n", resp.Status())
	// fmt.Printf("Check dashboard response body: %s\n", string(resp.Body()))

	return resp.StatusCode() == 200, nil
}

// createOrUpdateDashboard handles both creation and updating of dashboards
func createOrUpdateDashboard(splunk SplunkConfig, dashboard DashboardConfig) error {
	client := resty.New().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetTimeout(30 * time.Second).
		SetDebug(true)

	// Get authentication token
	authToken, err := getAuthToken(client, splunk)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	// Check if dashboard exists
	exists, err := checkDashboardExists(client, splunk, dashboard.TeamName)
	if err != nil {
		return fmt.Errorf("error checking dashboard existence: %v", err)
	}

	// Create dashboard XML with the new panels
	dashboardXML := fmt.Sprintf(`<?xml version="1.0"?>
		<dashboard version="1.1">
			<label>%s Dashboard</label>
			<row>
				<panel>
					<title>Time Range Selector</title>
					<input type="time" token="timeRange">
						<label>Time Range</label>
						<default>
							<earliest>-24h</earliest>
							<latest>now</latest>
						</default>
					</input>
				</panel>
			</row>
			<row>
				<panel>
					<title>Login Success VS. Failure</title>
					<chart>
						<search>
							<query>index="user_management_api_dev" uri="/users/login" | eval login_status=if(response_code=200, "Success", "Failure") | stats count by login_status | eval login_status=if(login_status=="Success", "A_Success", "B_Failure") | sort login_status | eval login_status=replace(login_status, "A_", "") | eval login_status=replace(login_status, "B_", "")</query>
							<earliest>$timeRange.earliest$</earliest>
							<latest>$timeRange.latest$</latest>
						</search>
					</chart>
				</panel>
				<panel>
					<title>Response Codes Distribution</title>
					<chart>
						<search>
							<query>index="user_management_api_dev" | stats count by response_code</query>
							<earliest>$timeRange.earliest$</earliest>
							<latest>$timeRange.latest$</latest>
						</search>
					</chart>
				</panel>
				</row>
			<row>
				<panel>
					<title>Number of API Hits</title>
					<chart>
						<search>
							<query>index="user_management_api_dev" | stats count as API_Hits</query>
							<earliest>$timeRange.earliest$</earliest>
							<latest>$timeRange.latest$</latest>
						</search>
					</chart>
				</panel>
				<panel>
					<title>Most Active Endpoint</title>
					<chart>
						<search>
							<query>index="user_management_api_dev" method=* uri=* | stats count by uri | sort -count | head 1</query>
							<earliest>$timeRange.earliest$</earliest>
							<latest>$timeRange.latest$</latest>
						</search>
					</chart>
				</panel>
				</row>
			<row>
				<panel>
					<title>Average Response Time by URI</title>
					<chart>
						<search>
							<query>index="user_management_api_dev" | stats avg(response_time) as avg_response_time by uri</query>
							<earliest>$timeRange.earliest$</earliest>
							<latest>$timeRange.latest$</latest>
						</search>
					</chart>
				</panel>
			</row>
		</dashboard>`, dashboard.TeamName)

	// Prepare form data
	formData := url.Values{}
	formData.Set("name", fmt.Sprintf("dashboard_%s", strings.ToLower(dashboard.TeamName)))
	formData.Set("eai:data", dashboardXML)
	formData.Set("output_mode", "json")

	var apiURL string
	if exists {
		// Update existing dashboard
		apiURL = fmt.Sprintf("%s/servicesNS/admin/search/data/ui/views/%s",
			strings.TrimSuffix(splunk.BaseURL, "/"),
			url.PathEscape(fmt.Sprintf("dashboard_%s", strings.ToLower(dashboard.TeamName))))
		fmt.Println("Updating existing dashboard...")
		resp, err := client.R().
			SetHeader("Authorization", "Basic "+authToken).
			SetHeader("Content-Type", "application/x-www-form-urlencoded").
			SetBody(formData.Encode()).
			Post(apiURL)

		if err != nil {
			return fmt.Errorf("error making request to update dashboard: %v", err)
		}

		if resp.StatusCode() != 200 {
			return fmt.Errorf("failed to update dashboard: %s - %s", resp.Status(), string(resp.Body()))
		}
	} else {
		// Create new dashboard
		apiURL = fmt.Sprintf("%s/servicesNS/admin/search/data/ui/views",
			strings.TrimSuffix(splunk.BaseURL, "/"))
		fmt.Println("Creating new dashboard...")
		resp, err := client.R().
			SetHeader("Authorization", "Basic "+authToken).
			SetHeader("Content-Type", "application/x-www-form-urlencoded").
			SetBody(formData.Encode()).
			Post(apiURL)

		if err != nil {
			return fmt.Errorf("error making request to create dashboard: %v", err)
		}

		if resp.StatusCode() != 201 {
			return fmt.Errorf("failed to create dashboard: %s - %s", resp.Status(), string(resp.Body()))
		}
	}

	fmt.Printf("Dashboard successfully %s for %s\n",
		map[bool]string{true: "updated", false: "created"}[exists],
		dashboard.TeamName)
	return nil
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Get credentials with fallback to default values
	username := getEnvWithDefault("SPLUNK_USERNAME", "admin")
	password := getEnvWithDefault("SPLUNK_PASSWORD", "")
	baseURL := getEnvWithDefault("SPLUNK_BASE_URL", "https://localhost:8089")

	if password == "" {
		fmt.Println("Error: Splunk password not set in environment variables")
		return
	}

	splunk := SplunkConfig{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
	}

	dashboard := DashboardConfig{
		TeamName:   "Team1",
		Index:      "user_management_api_dev",
		PanelQuery: "index=user_management_api_dev | head 5",
	}

	if err := createOrUpdateDashboard(splunk, dashboard); err != nil {
		fmt.Printf("Error creating/updating dashboard: %v\n", err)
		return
	}
}

// getEnvWithDefault returns environment variable value or default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
