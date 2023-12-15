package main

import (
	"flag"
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

	l.Info("Invoice Itaú successfully consumed!")
	l.Info("Starting to generate fatura.xlsx...")

	if err := internal.GenerateOrganizzeXLXSSheet(entries); err != nil {
		l.Error(err.Error())
		return
	}

	l.Info("Finished, fatura.xlsx was generated! Please run 'make run unoconv' to convert fatura.xlsx to fatura.xls")
}
