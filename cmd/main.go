package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/Rhymond/go-money"
	"github.com/viniciusgabrielfo/xls"
	"github.com/xuri/excelize/v2"
)

type Entry struct {
	Date        string
	Description string
	Value       *money.Money
}

var invoicePath string

func init() {
	flag.StringVar(&invoicePath, "file", "", "itaú invoice path to consume")
	flag.Parse()
}

func main() {

	l := slog.Default()

	l.Info("Starting invoice-itau-consumer...")

	entries, err := getEntriesFromItauInvoice(invoicePath)
	if err != nil {
		l.Error(err.Error())
		return
	}

	l.Info("Invoice Itaú successfully consumed!")
	l.Info("Starting to generate fatura.xlsx...")

	if err := generateXLXSSheet(entries); err != nil {
		l.Error(err.Error())
		return
	}

	l.Info("Finished, fatura.xlsx was generated! Please run 'make run unoconv' to convert fatura.xlsx to fatura.xls")
}

func getEntriesFromItauInvoice(filePath string) ([]Entry, error) {
	f, err := xls.Open(filePath, "utf-8")
	if err != nil {
		return nil, err
	}

	sheet := f.GetSheet(0)

	if sheet == nil {
		return nil, errors.New("invalid sheet")
	}

	var isEntry bool

	entries := make([]Entry, 0)

	for i := 0; i <= int(sheet.MaxRow); i++ {
		row := sheet.Row(i)
		if row == nil {
			if isEntry {
				isEntry = false
			}
			continue
		}

		col1 := row.Col(0)
		col2 := row.Col(1)

		if col1 == "data" && col2 == "lançamento" {
			isEntry = true
			continue
		}

		if isEntry {
			if col1 == "" || col2 == "dólar de conversão" {
				continue
			}

			col4 := row.Col(3)

			value, err := strconv.ParseFloat(col4, 32)
			if err != nil {
				return entries, err
			}

			entries = append(entries, Entry{
				Date:        col1,
				Description: col2,
				Value:       money.NewFromFloat(-value, money.BRL),
			})
		}
	}

	return entries, nil
}

func generateXLXSSheet(entries []Entry) error {
	f := excelize.NewFile()

	defer func() error {
		if err := f.Close(); err != nil {
			return err
		}

		return nil
	}()

	f.SetCellValue("Sheet1", "A1", "Data")
	f.SetCellValue("Sheet1", "B1", "Descrição")
	f.SetCellValue("Sheet1", "C1", "Categoria")
	f.SetCellValue("Sheet1", "D1", "Valor")
	f.SetCellValue("Sheet1", "E1", "Situação")

	for i := 0; i < len(entries); i++ {
		row := i + 2

		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), entries[i].Date)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), entries[i].Description)
		f.SetCellFloat("Sheet1", fmt.Sprintf("D%d", row), entries[i].Value.AsMajorUnits(), 2, 32)

	}

	if err := f.SaveAs("fatura.xlsx"); err != nil {
		return err
	}

	return nil
}
