package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sivanovie/appd-stats/pkg/appd"
	"github.com/sivanovie/appd-stats/pkg/conf"
	report "github.com/sivanovie/appd-stats/pkg/excel"
)

func main() {

	// LOGGER
	logname := "appd-stats.log"
	_, err := conf.SetLogger(logname)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	// CONF
	conf := conf.LoadConf()

	// PER CONTROLLER
	for i := range conf.Stats {

		// VARS
		controller := conf.Stats[i].Name
		url := conf.Stats[i].Url
		client := conf.Stats[i].Client
		secret := conf.Stats[i].Secret
		account := conf.Stats[i].Account
		auth := conf.Stats[i].Auth
		reportName := conf.Stats[i].Report.Name
		reportSubtitle := conf.Stats[i].Report.Subtitle
		reportHeaderB2 := conf.Stats[i].Report.Header.B2
		reportHeaderB3 := conf.Stats[i].Report.Header.B3
		reportHeaderB4 := conf.Stats[i].Report.Header.B4
		reportHeaderB5 := conf.Stats[i].Report.Header.B5
		scope := conf.Stats[i].Report.Scope
		team := conf.Stats[i].Report.Team
		description := conf.Stats[i].Report.Description
		timerangePref := strings.ToLower(conf.Stats[i].Report.Timerange)

		// Set time range

		// end
		reportTimeEnd := time.Now().UnixMilli()

		// start
		var reportTimeStart int64

		// day in milliseconds
		day := int64(86400000)

		if timerangePref == "last 1 day" {

			reportTimeStart = reportTimeEnd - day

		} else if timerangePref == "last 1 week" {

			reportTimeStart = reportTimeEnd - (7 * day)

		} else if timerangePref == "last 2 weeks" {

			reportTimeStart = reportTimeEnd - (14 * day)

		} else if timerangePref == "last 1 month" {

			reportTimeStart = reportTimeEnd - (30 * day)

		} else if timerangePref == "last 3 months" {

			reportTimeStart = reportTimeEnd - (90 * day)

		} else if timerangePref == "last 6 months" {

			reportTimeStart = reportTimeEnd - (180 * day)

		} else if timerangePref == "last 1 year" {

			reportTimeStart = reportTimeEnd - (360 * day)

		} else {

			log.Printf("Unsupported timerange %v.", timerangePref)
			break

		}

		// LOGIN
		err, logincookies := appd.GetLoginCookies(url, auth)
		if err != nil {
			log.Println("Couldn't login to Controller.")
		}

		// TOKEN
		log.Printf("Fetching temp token for %v.", controller)

		err, token := appd.GetControllerAccessToken(client, account, secret, url)
		if err != nil {
			log.Println("Couldn't retrieve access token.")
		}

		// ALL APPS
		appsurl := url + "/controller/rest/applications?output=json"
		apps, err := appd.GetEntitiesFromController(appsurl, token)
		if err != nil {
			log.Println("Couldn't get all apps for controller.")
		}

		// Get total number of calls and other summary stats
		err, appsWithMetrics := appd.GetAllAppsSummaryStats(url, logincookies, apps, reportTimeStart, reportTimeEnd)
		if err != nil {
			log.Println(err)
		}

		err, appsWithMetricsAndHrs := appd.GetHealthRules(url, token, appsWithMetrics)
		if err != nil {
			log.Println(err)
		}

		err = report.BuildExcelReport(
			appsWithMetricsAndHrs,
			controller,
			time.UnixMilli(reportTimeStart).Format(time.RFC3339),
			time.UnixMilli(reportTimeEnd).Format(time.RFC3339),
			reportName,
			scope,
			team,
			description,
			url,
			reportHeaderB2,
			reportHeaderB3,
			reportHeaderB4,
			reportHeaderB5,
			reportSubtitle)

	}

}
