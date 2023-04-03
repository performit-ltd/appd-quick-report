package appd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"
)

type AppDetails struct {
	Name     string
	Id       float64
	Metrics  AppMetrics
	Alerting []AppHealthRules
}
type AppMetrics struct {
	NumberOfErrors              int64   `json:"numberOfErrors"`
	ErrorsPerMinute             float64 `json:"errorsPerMinute"`
	NumberOfCalls               int64   `json:"numberOfCalls"`
	CallsPerMinute              float64 `json:"callsPerMinute"`
	AverageResponseTime         float64 `json:"averageResponseTime"`
	NumberOfActiveHealthRules   float64 `json:"NumberOfActiveHealthRules"`
	NumberOfInactiveHealthRules float64 `json:"NumberOfInactiveHealthRules"`
}
type AppHealthRules struct {
	Name   string
	Id     float64
	Active bool
}
type AppStatisticsPayload struct {
	RequestFilter  []int    `json:"requestFilter"`
	TimeRangeStart int64    `json:"timeRangeStart"`
	TimerangeEnd   int64    `json:"timeRangeEnd"`
	SearchFilters  []string `json:"searchFilters"`
	ColumnSorts    []string `json:"columnSorts"`
	ResultColumns  []string `json:"resultColumns"`
	Offset         int      `json:"offset"`
	Limit          int      `json:"limit"`
}

func GetEntitiesFromController(URL string, controllerAccessToken string) ([]AppDetails, error) {

	var jsonArrInterface []interface{}
	var apps []AppDetails

	// Set HTTP request method
	method := "GET"

	// timeout
	timeout := time.Duration(60 * time.Second)

	// http client
	client := &http.Client{}

	// client
	client = &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{},
	}

	// Create a new HTTP request object
	req, err := http.NewRequest(method, URL, nil)

	// Return error if new request creation fails
	if err != nil {
		log.Printf("ERROR - %v", err)
		return nil, err
	}

	// Set HTTP headers
	req.Header.Set("Authorization", "Bearer "+controllerAccessToken)

	// Make the HTTP request to the Controller
	res, err := client.Do(req)

	// If non-http error is returned from response we quit this goroutine
	if err != nil {
		log.Printf("ERROR - %v", err)
		return nil, err
	}

	// Close the body stream to avoid leaks later
	defer res.Body.Close()

	// Read the body into a byte var
	body, err := ioutil.ReadAll(res.Body)

	// If there is an error while reading response body we quit this function
	if err != nil {
		log.Printf("ERROR - %v", err)
		return nil, err
	}

	// If HTTP state from controller is bad we quit this goroutine
	if res.StatusCode != 200 {
		log.Printf("ERROR - Got HTTP state %v while calling %v. %v", res.StatusCode, URL, string(body))
		return nil, errors.New(fmt.Sprint(res.StatusCode))
	}

	// Unmarshal retrieved APM apps list
	json.Unmarshal(body, &jsonArrInterface)

	for i := range jsonArrInterface {
		appdetails := AppDetails{
			Name: jsonArrInterface[i].(map[string]interface{})["name"].(string),
			Id:   jsonArrInterface[i].(map[string]interface{})["id"].(float64),
		}
		apps = append(apps, appdetails)
	}

	return apps, nil

}

func GetControllerAccessToken(clientName string, account string, clientSecret string, controllerUrl string) (error, string) {

	var jsonMap map[string]interface{}

	// Set HTTP request method
	method := "POST"

	// timeout
	timeout := time.Duration(60 * time.Second)

	// Convert HTTP request payload
	authpayload := "grant_type=client_credentials&client_id=" + clientName + "@" + account + "&client_secret=" + clientSecret
	payload := strings.NewReader(authpayload)

	// Set the auth URL
	authurl := controllerUrl + "/api/oauth/access_token"

	// http client
	client := &http.Client{}

	// client
	client = &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{},
	}

	// Create a new HTTP request object
	req, err := http.NewRequest(method, authurl, payload)

	// If there is an error while trying to create new http request object we quit this goroutine
	if err != nil {
		log.Printf("ERROR - %v", err)
		return err, ""
	}

	// Add the needed headers for temporary access token request
	req.Header.Add("Content-Type", "application/vnd.appd.cntrl+protobuf;v=1")

	// Make the HTTP request to the Controller
	res, err := client.Do(req)

	// If non-http error is returned from response we quit this goroutine
	if err != nil {
		log.Printf("ERROR - %v", err)
		return err, ""
	}

	// Close the body stream to avoid leaks later
	defer res.Body.Close()

	// Read the body into a byte var
	body, err := ioutil.ReadAll(res.Body)

	// If there is an error while reading response body we quit this goroutine
	if err != nil {
		log.Printf("ERROR - %v", err)
		return err, ""
	}

	// If HTTP state from controller is bad we quit this goroutine
	if res.StatusCode != 200 {
		log.Printf("ERROR - Got HTTP state %v while waiting for temp access token from Controller.", res.StatusCode)
		return errors.New(fmt.Sprint(res.StatusCode)), ""
	}

	log.Printf("Got temp access token from Controller (http %v).", res.StatusCode)

	// Extract the Controller temporary access token from the body
	// Access tokens are valid for 5 minutes by default
	json.Unmarshal(body, &jsonMap)
	controllerAccessTokenRaw := jsonMap["access_token"]

	// Convert the access token to string
	controllerAccessToken := fmt.Sprint(controllerAccessTokenRaw)

	// Validate temporary access token based on length if more than 100
	// Typical length is around 600
	if len(controllerAccessToken) > 100 {
		log.Printf("Validated temporary access token (length: %v).", len(controllerAccessToken))
	} else {
		log.Printf("WARN - Got a shorter access token from Controller. Expected >100 chars, got %v).", len(controllerAccessToken))
	}

	return nil, controllerAccessToken

}

func GetLoginCookies(controllerUrl string, auth string) (error, []*http.Cookie) {

	var (
		cookies []*http.Cookie
	)

	// Login URL
	loginurl := controllerUrl + "/auth?action=login"

	log.Printf("Calling %v for login cookies.", loginurl)

	// Set HTTP client timeout
	timeout := time.Duration(60 * time.Second)

	// http client
	clientAppdController := &http.Client{}

	// client
	clientAppdController = &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{},
	}

	// Get new request object
	req, err := http.NewRequest("GET", loginurl, nil)
	if err != nil {
		fmt.Println(err)
		return err, nil
	}

	// Authorization based on Controller temporary access token (OAuth)
	req.Header.Add("Authorization", "Basic "+auth)

	// Make the call to Controller
	resp, err := clientAppdController.Do(req)
	if err != nil {
		fmt.Println(err)
		return err, nil
	}

	// Close the response stream
	defer resp.Body.Close()

	// Validate if login cookies are returned by Controller, and then extract them
	if resp.StatusCode == 200 && strings.Contains(fmt.Sprint(resp.Cookies()), "JSESSIONID") && strings.Contains(fmt.Sprint(resp.Cookies()), "X-CSRF-TOKEN") {

		// Get the login cookies
		for _, cookie := range resp.Cookies() {

			// Extract the JSESSIONID and X-CSRF-TOKEN and add them to cookies
			if cookie.Name == "JSESSIONID" || cookie.Name == "X-CSRF-TOKEN" {

				cookies = append(cookies, &http.Cookie{
					Name:   cookie.Name,
					Value:  cookie.Value,
					MaxAge: 300,
				})

			}

		}

	} else {

		// Error out if required cookies are not returned
		errmsg := errors.New("ERROR - Couldn't find JSESSIONID/X-CSRF-TOKEN in Controller response (http " + fmt.Sprint(resp.StatusCode) + "). Cookies: " + fmt.Sprint(resp.Cookies()))
		log.Println(errmsg)
		return errmsg, nil

	}

	log.Println("Got login cookies for Controller.")

	return nil, cookies

}

func GetAllAppsSummaryStats(controllerUrl string, login []*http.Cookie, appsinfo []AppDetails, startTime int64, endTime int64) (error, []AppDetails) {

	var payload AppStatisticsPayload
	var appDetailsWithMetrics []AppDetails
	method := "POST"

	// Set the base flow map URL
	listappsurl := controllerUrl + "/controller/restui/v1/app/list/ids"

	// Get timerange
	payload.TimerangeEnd = endTime
	payload.TimeRangeStart = startTime

	// Get resultColumns
	payload.ResultColumns = []string{
		"APP_OVERALL_HEALTH",
		"CALLS",
		"CALLS_PER_MINUTE",
		"AVERAGE_RESPONSE_TIME",
		"ERROR_PERCENT",
		"ERRORS",
		"ERRORS_PER_MINUTE",
	}

	// Get limits
	payload.Limit = -1

	for i := range appsinfo {

		// Attach this app id to []int
		payload.RequestFilter = append(payload.RequestFilter, int(appsinfo[i].Id))

	}

	// Make it JSON
	JSONpayload, _ := json.Marshal(&payload)

	// strings.Reader
	data := strings.NewReader(string(JSONpayload))

	// http client
	client := &http.Client{}

	// Create HTTP request
	req, err := http.NewRequest(method, listappsurl, data)
	if err != nil {
		return err, appDetailsWithMetrics
	}

	// Add the login cookies and headers to request
	for i := range login {

		// Add the login cookies to request
		req.AddCookie(login[i])

		// Add the X-CSRF-TOKEN as a request header
		if login[i].Name == "X-CSRF-TOKEN" {
			req.Header.Add(login[i].Name, login[i].Value)
		}

	}

	// More headers
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("Accept", "application/json, text/plain, */*")

	// Make the call
	res, err := client.Do(req)
	if err != nil || res.StatusCode != 200 {

		err := errors.New(fmt.Sprint(err) + fmt.Sprint(res))
		return err, appDetailsWithMetrics

	}

	// Close stream
	defer res.Body.Close()

	// empty map
	var allstats map[string]interface{}

	// response body to byte
	body, err := ioutil.ReadAll(res.Body)

	// JSON
	json.Unmarshal(body, &allstats)

	// Get all apps stats
	var appstats interface{}
	for i := range allstats["data"].([]interface{}) {

		appstats = allstats["data"].([]interface{})[i]

		for ii := range appsinfo {
			if appstats.(map[string]interface{})["name"].(string) == appsinfo[ii].Name {

				// Number of Calls (convert from scientific notation)
				str := fmt.Sprint(appstats.(map[string]interface{})["numberOfCalls"].(float64))
				float, _, _ := big.ParseFloat(str, 10, 64, big.ToNearestEven)
				noc, _ := float.Int64()
				appsinfo[ii].Metrics.NumberOfCalls = noc

				// Number of Errors (convert from scientific notation)
				str = fmt.Sprint(appstats.(map[string]interface{})["numberOfErrors"].(float64))
				float, _, _ = big.ParseFloat(str, 10, 64, big.ToNearestEven)
				noe, _ := float.Int64()
				appsinfo[ii].Metrics.NumberOfErrors = noe

				// Average Response Time
				appsinfo[ii].Metrics.AverageResponseTime = appstats.(map[string]interface{})["averageResponseTime"].(float64)

				// Calls per Minute
				appsinfo[ii].Metrics.CallsPerMinute = appstats.(map[string]interface{})["callsPerMinute"].(float64)

				// Errors per Minute
				appsinfo[ii].Metrics.ErrorsPerMinute = appstats.(map[string]interface{})["errorsPerMinute"].(float64)

			}
		}

	}

	return nil, appsinfo
}

func GetHealthRules(controllerUrl string, token string, appsinfo []AppDetails) (error, []AppDetails) {

	var jsonArrInterface []interface{}

	method := "GET"

	// timeout
	timeout := time.Duration(60 * time.Second)

	// http client
	client := &http.Client{}

	// client
	client = &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{},
	}

	for i := range appsinfo {

		// Set the HR url
		hrurl := controllerUrl + "/controller/alerting/rest/v1/applications/" + fmt.Sprint(appsinfo[i].Id) + "/health-rules"

		// Create a new HTTP request object
		req, err := http.NewRequest(method, hrurl, nil)

		// Return error if new request creation fails
		if err != nil {
			log.Printf("ERROR - %v", err)
			return err, appsinfo
		}

		// Set HTTP headers
		req.Header.Set("Authorization", "Bearer "+token)

		// Make the HTTP request to the Controller
		res, err := client.Do(req)

		// If non-http error is returned from response we quit this goroutine
		if err != nil {
			log.Printf("ERROR - %v", err)
			return err, appsinfo
		}

		// Close the body stream to avoid leaks later
		defer res.Body.Close()

		// Read the body into a byte var
		body, err := ioutil.ReadAll(res.Body)

		// If there is an error while reading response body we quit this function
		if err != nil {
			log.Printf("ERROR - %v", err)
			return err, appsinfo
		}

		// If HTTP state from controller is bad we quit this goroutine
		if res.StatusCode != 200 {
			log.Printf("ERROR - Got HTTP state %v while calling %v. %v", res.StatusCode, hrurl, string(body))
			return errors.New(fmt.Sprint(res.StatusCode)), appsinfo
		}

		// Unmarshal retrieved APM apps list
		json.Unmarshal(body, &jsonArrInterface)

		// Iterate hrs
		for ii := range jsonArrInterface {

			alertInfo := AppHealthRules{
				Name:   jsonArrInterface[ii].(map[string]interface{})["name"].(string),
				Id:     jsonArrInterface[ii].(map[string]interface{})["id"].(float64),
				Active: jsonArrInterface[ii].(map[string]interface{})["enabled"].(bool),
			}
			appsinfo[i].Alerting = append(appsinfo[i].Alerting, alertInfo)

		}

		active := 0
		inactive := 0
		for ii := range appsinfo[i].Alerting {

			if appsinfo[i].Alerting[ii].Active == true {

				active++

			} else {

				inactive++

			}

		}
		appsinfo[i].Metrics.NumberOfActiveHealthRules = float64(active)
		appsinfo[i].Metrics.NumberOfInactiveHealthRules = float64(inactive)

	}

	return nil, appsinfo

}
