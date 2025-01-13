package dashboard

import (
	"crypto/tls"
	"dashboard/models"
	"dashboard/utils/splunk"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

func RunDashboardSetup(cfg *models.SplunkConfig, dashboardConfig *models.DashboardConfig) error {
	client := resty.New().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		SetTimeout(30 * time.Second)

	// Check if the dashboard exists
	exists, err := splunk.CheckDashboardExists(client, *cfg, dashboardConfig.TeamName)
	if err != nil {
		return fmt.Errorf("error checking dashboard existence: %v", err)
	}

	// Load the dashboard template
	template, err := os.ReadFile("dashboard_template_json.xml")
	if err != nil {
		return fmt.Errorf("error loading dashboard template: %v", err)
	}

	dashboardXML := fmt.Sprintf(string(template), dashboardConfig.TeamName)
	dashboardName := fmt.Sprintf("dashboard_%s", strings.ToLower(dashboardConfig.TeamName))

	// Create or update the dashboard
	if exists {
		log.Println("Updating existing dashboard...")
		err = UpdateDashboard(client, *cfg, dashboardName, dashboardXML)
	} else {
		log.Println("Creating new dashboard...")
		err = CreateDashboard(client, *cfg, dashboardName, dashboardXML)
	}

	if err != nil {
		return fmt.Errorf("error creating/updating dashboard: %v", err)
	}

	// Set permissions
	if err = SetDashboardPermissions(client, *cfg, dashboardName); err != nil {
		return fmt.Errorf("error setting dashboard permissions: %v", err)
	}

	log.Println("Dashboard setup complete")
	return nil
}
