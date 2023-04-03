package report

import (
	"fmt"

	"github.com/sivanovie/appd-stats/pkg/appd"
	"github.com/xuri/excelize/v2"
)

const (
	SheetName = "Controller Applications Report"
)

func BuildExcelReport(
	appsdetails []appd.AppDetails,
	profile string,
	timeRangeStart string,
	timeRangeEnd string,
	reportName string,
	scope string,
	team string,
	description string,
	controllerURL string,
	b2 string,
	b3 string,
	b4 string,
	b5 string,
	reportSubtitle string) error {

	var (
		err        error
		reportData = [][]interface{}{}
	)

	// Convert unix ts to human date

	// Create file
	f := excelize.NewFile()

	// Select default spreadsheet and rename
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		return err
	}

	// Select active sheet
	f.SetActiveSheet(index)

	// Rename active sheet
	f.SetSheetName("Sheet1", SheetName)

	// Set column width
	err = f.SetColWidth(SheetName, "A", "A", 6)
	err = f.SetColWidth(SheetName, "H", "H", 6)
	err = f.SetColWidth(SheetName, "B", "B", 30)
	err = f.SetColWidth(SheetName, "C", "C", 20)
	err = f.SetColWidth(SheetName, "D", "D", 20)
	err = f.SetColWidth(SheetName, "E", "E", 20)
	err = f.SetColWidth(SheetName, "F", "F", 20)
	err = f.SetColWidth(SheetName, "G", "G", 6)

	// Height of 1st row
	err = f.SetRowHeight(SheetName, 1, 12)

	// Merge 1st row cells
	err = f.MergeCell(SheetName, "A1", "H1")

	// Height of 2nd row
	err = f.SetRowHeight(SheetName, 2, 25)

	// Merge cells for B2 header
	err = f.MergeCell(SheetName, "B2", "D2")

	// Styling and font of B2 header
	style, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 20, Color: "6d64e8"}})
	err = f.SetCellStyle(SheetName, "B2", "D2", style)

	// Add value (B2 header)
	err = f.SetSheetRow(SheetName, "B2", &[]interface{}{b2})

	// Merge cells for B3 header
	err = f.MergeCell(SheetName, "B3", "D3")

	// Add value (B3 header)
	err = f.SetSheetRow(SheetName, "B3", &[]interface{}{b3})

	// Merge cells for B4 header
	err = f.MergeCell(SheetName, "B4", "D4")

	// Add value (B4 header)
	err = f.SetSheetRow(SheetName, "B4", &[]interface{}{b4})

	// Styling and font of B5 header
	style, err = f.NewStyle(&excelize.Style{Font: &excelize.Font{Color: "666666"}})
	err = f.MergeCell(SheetName, "B5", "D5")
	err = f.SetCellStyle(SheetName, "B5", "D5", style)

	// Add value (B5 header)
	err = f.SetSheetRow(SheetName, "B5", &[]interface{}{b5})

	// Styling and font of report name
	style, err = f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 32, Color: "2B4492", Bold: true}})
	err = f.MergeCell(SheetName, "B7", "G7")
	err = f.SetCellStyle(SheetName, "B7", "G7", style)

	// Add value (report name)
	err = f.SetSheetRow(SheetName, "B7", &[]interface{}{reportName})

	// Styling and font of report subtitle
	style, err = f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 13, Color: "E25184", Bold: true}})
	err = f.MergeCell(SheetName, "B8", "C8")
	err = f.SetCellStyle(SheetName, "B8", "C8", style)

	// Add value (report subtitle)
	err = f.SetSheetRow(SheetName, "B8", &[]interface{}{reportSubtitle})

	// Styling and font of section names for timerange and scope
	style, err = f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 13, Bold: true}})
	err = f.SetCellStyle(SheetName, "B10", "G10", style)

	// Add value (section names for timerange and scope)
	err = f.SetSheetRow(SheetName, "B10", &[]interface{}{"From", "", "Until", "", "Scope"})

	// Merge cells for section names (timerange and scope)
	err = f.MergeCell(SheetName, "B10", "C10")
	err = f.MergeCell(SheetName, "D10", "E10")
	err = f.MergeCell(SheetName, "F10", "G10")

	// Styling and font of section values for timerange and scope
	style, err = f.NewStyle(&excelize.Style{Font: &excelize.Font{Color: "666666"}})
	err = f.SetCellStyle(SheetName, "B11", "G11", style)

	// Add values (timerange and scope)
	err = f.SetSheetRow(SheetName, "B11", &[]interface{}{timeRangeStart, "", timeRangeEnd, "", scope})

	// Merge cells for section values (timerange and scope)
	err = f.MergeCell(SheetName, "B11", "C11")
	err = f.MergeCell(SheetName, "D11", "E11")
	err = f.MergeCell(SheetName, "F11", "G11")

	// Styling and font of section names for team and description
	style, err = f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 13, Bold: true}})
	err = f.SetCellStyle(SheetName, "B13", "G13", style)

	// Add value (section names for team and description)
	err = f.SetSheetRow(SheetName, "B13", &[]interface{}{"Team", "", "Description"})

	// Merge cells for section names (team and description)
	err = f.MergeCell(SheetName, "B13", "C13")
	err = f.MergeCell(SheetName, "D13", "E13")

	// Styling and font of section values for team and description
	style, err = f.NewStyle(&excelize.Style{Font: &excelize.Font{Color: "666666"}})
	err = f.SetCellStyle(SheetName, "B14", "G14", style)

	// Add values (team and description)
	err = f.SetSheetRow(SheetName, "B14", &[]interface{}{team, "", description})

	// Merge cells for section values (team and description)
	err = f.MergeCell(SheetName, "B14", "C14")
	err = f.MergeCell(SheetName, "D14", "E14")

	// Attach extracted apps details to []interface{}
	for i := range appsdetails {

		app := appsdetails[i]

		// Create the individual app record
		s := []interface{}{
			app.Name,
			app.Metrics.NumberOfErrors,
			app.Metrics.NumberOfCalls,
			app.Metrics.NumberOfActiveHealthRules,
			app.Metrics.NumberOfInactiveHealthRules}

		// Attach to all
		reportData = append(reportData, s)

	}

	// Styling and font of table column names for main table with data
	style, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 13, Bold: true, Color: "2B4492"},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	err = f.SetCellStyle(SheetName, "B17", "G17", style)

	// Table column names
	err = f.SetSheetRow(SheetName, "B17", &[]interface{}{"Application", "Number of Errors", "Number of Calls", "Enabled Alerts", "Disabled Alerts"})

	//err = f.MergeCell(SheetName, "D17", "E17")

	// Height of table column names row
	err = f.SetRowHeight(SheetName, 17, 32)

	// We start inserting actual table data from row 18
	startRow := 18

	// Start iterating through subsequent rows
	// For each row insert one line of the []interface{} app statistics data (essentially one app)
	for i := startRow; i < (len(reportData) + startRow); i++ {

		var fill string
		if i%2 == 0 {
			fill = "F3F3F3"
		} else {
			fill = "FFFFFF"
		}

		// Set styling and font for each table row
		style, err = f.NewStyle(&excelize.Style{
			Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{fill}},
			Font:      &excelize.Font{Color: "666666"},
			Alignment: &excelize.Alignment{Vertical: "center"},
		})
		err = f.SetCellStyle(SheetName, fmt.Sprintf("B%d", i), fmt.Sprintf("F%d", i), style)

		// Add row data
		err = f.SetSheetRow(SheetName, fmt.Sprintf("B%d", i), &reportData[i-18])

		/*
			// Add actual data
			err = f.SetCellRichText(SheetName, fmt.Sprintf("C%d", i), []excelize.RichTextRun{
				{
					Text: reportData[i-18][1].(string),
					Font: &excelize.Font{
						Bold: true}},
			})*/

		//err = f.MergeCell(SheetName, fmt.Sprintf("D%d", i), fmt.Sprintf("E%d", i))
		// Set row height
		err = f.SetRowHeight(SheetName, i, 18)

	}

	err = f.SaveAs(profile + ".xlsx")
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
