SPLUNK DASHBOARD CREATION THROUGH GOLANG:

Add env file whith the below contents:
SPLUNK_BASE_URL="https://127.0.0.1:8089"
SPLUNK_API_PATH="/servicesNS/Nobody/search/data/ui/views"
SECRETS_MANAGER_NAME="testing/splunkToken"
 

Steps to run the code:
go run main.go

Endpoint used to creade, update and delete dashboard
/servicesNS/Nobody/dashboardApp/data/ui/views​

​
Refernece:​
REST api splunk doccumentation
https://docs.splunk.com/Documentation/Splunk/7.2.0/RESTREF/RESTknowledge?_gl=1*5lyxk4*_gcl_au*MTY2MTE2NDE1Ni4xNzI4ODI5MDM1*FPAU*MTY2MTE2NDE1Ni4xNzI4ODI5MDM1*_ga*NDU2NzA4MDU0LjE3Mjg4MjkwMzU.*_ga_5EPM2P39FV*MTczMTMxNDgwOC42OC4xLjE3MzEzMTQ4MjIuNDYuMC45MjMyNTUzMTE.*_fplc*ZDZBQlJUQXM5UjkzY3lLQTMlMkZyZjdBNnlmMUE1bzg2TEc1JTJGc1hMbWc5RUFYMjR1V2lLdDBabjJzUmlYZzJSZXp4VkhzRU8wOUg4OVJKb1JFbWtMMnloYnR4NGRzJTJGVjR3NkdyJTJGeUl5SlBLejJyMWo3RE8lMkJhT0R0a3B1cjRIdyUzRCUzRA..#data.2Fui.2Fviews