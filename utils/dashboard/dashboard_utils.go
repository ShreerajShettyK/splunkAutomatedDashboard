package dashboard

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"dashboard/models"
)

func CreateDashboard(client *resty.Client, splunk models.SplunkConfig, dashboardName, dashboardXML string) error {
	apiURL := fmt.Sprintf("%s%s", strings.TrimSuffix(splunk.BaseURL, "/"), splunk.ApiPath)

	formData := url.Values{}
	formData.Set("name", dashboardName)
	formData.Set("eai:data", strings.TrimSpace(dashboardXML))
	formData.Set("output_mode", "json")

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+splunk.Token).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(formData.Encode()).
		Post(apiURL)

	if err != nil {
		return fmt.Errorf("error creating dashboard: %v", err)
	}

	if resp.StatusCode() != 201 {
		return fmt.Errorf("failed to create dashboard: %s - %s", resp.Status(), string(resp.Body()))
	}

	return nil
}

func UpdateDashboard(client *resty.Client, splunk models.SplunkConfig, dashboardName, dashboardXML string) error {
	apiURL := fmt.Sprintf("%s%s/%s", strings.TrimSuffix(splunk.BaseURL, "/"), splunk.ApiPath, url.PathEscape(dashboardName))

	formData := url.Values{}
	formData.Set("eai:data", strings.TrimSpace(dashboardXML))
	formData.Set("output_mode", "json")

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+splunk.Token).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(formData.Encode()).
		Post(apiURL)

	if err != nil {
		return fmt.Errorf("error updating dashboard: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to update dashboard: %s - %s", resp.Status(), string(resp.Body()))
	}

	return nil
}

func SetDashboardPermissions(client *resty.Client, splunk models.SplunkConfig, dashboardName string) error {
	permissionsURL := fmt.Sprintf("%s%s/%s/acl", strings.TrimSuffix(splunk.BaseURL, "/"), splunk.ApiPath, url.PathEscape(dashboardName))

	formData := url.Values{}
	formData.Set("sharing", "app")
	formData.Set("owner", "leps")
	formData.Set("perms.read", "*")
	formData.Set("perms.write", "leps")

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+splunk.Token).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(formData.Encode()).
		Post(permissionsURL)

	if err != nil {
		return fmt.Errorf("error setting permissions: %v", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to set permissions: %s - %s", resp.Status(), string(resp.Body()))
	}

	log.Println("Permissions set successfully")
	return nil
}
