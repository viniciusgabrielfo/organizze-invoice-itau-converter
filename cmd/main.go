package main

import (
	"flag"
	"fmt"
	"log/slog"
	"time"

	"github.com/viniciusgabrielfo/organizze-invoice-itau-converter/internal"
)

var (
	invoicePath string
	startDate   string
	endDate     string
)

func init() {
	flag.StringVar(&invoicePath, "file", "", "itaú invoice path to consume")
	flag.StringVar(&startDate, "start-date", "", "only consume from start-date (02/01/2006)")
	flag.StringVar(&endDate, "end-date", "", "only consume until end-date (02/01/2006)")
	flag.Parse()
}

func main() {
	l := slog.Default()

	l.Info("Starting invoice-itau-consumer...")

	itauImportConfigs := &internal.ItauImportConfigs{}

	if startDate != "" {
		tStartDate, err := time.Parse("02/01/2006", startDate)
		if err != nil {
			panic(err)
		}

		itauImportConfigs.StartDate = tStartDate
	}

	if endDate != "" {
		tEndDate, err := time.Parse("02/01/2006", endDate)
		if err != nil {
			panic(err)
		}

		itauImportConfigs.EndDate = tEndDate
	}

	entries, err := internal.GetEntriesFromItauInvoice(itauImportConfigs, invoicePath)
	if err != nil {
		l.Error(err.Error())
		return
	}

	l.Info(fmt.Sprintf("Itaú's invoice successfully consumed from %s!", invoicePath))
	l.Info(fmt.Sprintf("Starting to generate %s...", internal.OrganizzeSheetName))

	if err := internal.GenerateOrganizzeXLXSSheet(entries); err != nil {
		l.Error(err.Error())
		return
	}

	l.Info(fmt.Sprintf("Finished, %s was generated!", internal.OrganizzeSheetName))
	l.Info(fmt.Sprintf("Please run 'make run unoconv' to convert %s to organizze-entries-to-import.xls", internal.OrganizzeSheetName))
}
