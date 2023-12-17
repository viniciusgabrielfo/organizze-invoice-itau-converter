package main

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/viniciusgabrielfo/organizze-invoice-itau-converter/internal"
)

var invoicePath string

func init() {
	flag.StringVar(&invoicePath, "file", "", "itaú invoice path to consume")
	flag.Parse()
}

func main() {
	l := slog.Default()

	l.Info("Starting invoice-itau-consumer...")

	entries, err := internal.GetEntriesFromItauInvoice(invoicePath)
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
