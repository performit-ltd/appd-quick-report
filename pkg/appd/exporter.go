package appd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func GenerateCSV(profile string, controllerAppsWithDetails []AppDetails) error {

	filename := profile + ".csv"
	csvrecords := [][]string{
		{
			"Application Name",
			"Controller",
			"Number of Calls (last day)",
			"Number of Errors (last day)",
			"Calls per Minute (last day)",
			"Errors per Minute (last day)",
			"Active Alerts (health rules)",
			"Inactive Alerts (health rules)",
			"Alert List (name, id, enabled)"},
	}

	// Remove existing CSV file
	os.Remove(filename)

	for i := range controllerAppsWithDetails {

		app := controllerAppsWithDetails[i]

		csvrecords = append(csvrecords, []string{
			app.Name,
			profile,
			fmt.Sprint(app.Metrics.NumberOfCalls),
			fmt.Sprint(app.Metrics.NumberOfErrors),
			fmt.Sprint(app.Metrics.CallsPerMinute),
			fmt.Sprint(app.Metrics.ErrorsPerMinute),
			fmt.Sprint(app.Metrics.NumberOfActiveHealthRules),
			fmt.Sprint(app.Metrics.NumberOfInactiveHealthRules),
			fmt.Sprint(app.Alerting),
		})
	}

	f, e := os.Create(filename)
	if e != nil {
		log.Println(e)
		return e
	}

	writer := csv.NewWriter(f)
	writer.Comma = ';'
	e = writer.WriteAll(csvrecords)
	if e != nil {
		log.Println(e)
		return e
	}

	return nil

}
