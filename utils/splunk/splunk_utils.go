package splunk

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"dashboard/models"
)

func CheckDashboardExists(client *resty.Client, splunk models.SplunkConfig, dashboardName string) (bool, error) {
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
