package internal

import (
	"errors"
	"strconv"

	"github.com/Rhymond/go-money"
	"github.com/viniciusgabrielfo/organizze-invoice-itau-converter/pkg/category_definer"
	"github.com/viniciusgabrielfo/organizze-invoice-itau-converter/pkg/model"
	"github.com/viniciusgabrielfo/xls"
)

func GetEntriesFromItauInvoice(filePath string) ([]model.Entry, error) {
	f, err := xls.Open(filePath, "utf-8")
	if err != nil {
		return nil, err
	}

	sheet := f.GetSheet(0)

	if sheet == nil {
		return nil, errors.New("invalid sheet")
	}

	var isEntry bool

	entries := make([]model.Entry, 0)

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

			entries = append(entries, model.Entry{
				Date:        col1,
				Description: col2,
				Category:    category_definer.GetCategoryFromDescription(col2),
				Value:       money.NewFromFloat(-value, money.BRL),
			})
		}
	}

	return entries, nil
}
