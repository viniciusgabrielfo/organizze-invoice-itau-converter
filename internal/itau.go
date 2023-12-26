package internal

import (
	"errors"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/viniciusgabrielfo/organizze-invoice-itau-converter/pkg/category_definer"
	"github.com/viniciusgabrielfo/organizze-invoice-itau-converter/pkg/model"
	"github.com/viniciusgabrielfo/xls"
)

type ItauImportConfigs struct {
	StartDate time.Time
	EndDate   time.Time
}

func GetEntriesFromItauInvoice(configs *ItauImportConfigs, filePath string) ([]model.Entry, error) {
	logger := slog.Default()

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

		date := row.Col(0)
		description := row.Col(1)

		if date == "data" && description == "lançamento" {
			isEntry = true
			continue
		}

		if isEntry {
			if date == "" || description == "dólar de conversão" {
				continue
			}

			entryDate, err := time.Parse("02/01/2006", date)
			if err != nil {
				logger.Error(err.Error())
				continue
			}

			if !IsBetweenConfigInternal(configs, entryDate) {
				continue
			}

			value, err := strconv.ParseFloat(row.Col(3), 32)
			if err != nil {
				return entries, err
			}

			if ok, installments := IsInstallmentPurchase(description); ok {
				value = value * float64(installments)
			}

			entries = append(entries, model.Entry{
				Date:        date,
				Description: description,
				Category:    category_definer.GetCategoryFromDescription(description),
				Value:       money.NewFromFloat(-value, money.BRL),
			})
		}
	}

	return entries, nil
}

func IsBetweenConfigInternal(configs *ItauImportConfigs, date time.Time) bool {
	if !configs.StartDate.IsZero() {
		if date.Before(configs.StartDate) {
			return false
		}
	}

	if !configs.EndDate.IsZero() {
		if date.After(configs.EndDate) {
			return false
		}
	}

	return true
}

func IsInstallmentPurchase(description string) (bool, int32) {
	logger := slog.Default()

	re, err := regexp.Compile("[0-9]+/[0-9]+")
	if err != nil {
		logger.Error(err.Error())
		return false, 0
	}

	installmentPattern := re.FindAllString(description, -1)
	if len(installmentPattern) == 0 {
		return false, 0
	}

	s := strings.Split(installmentPattern[0], "/")

	i, err := strconv.ParseInt(s[1], 10, 32)
	if err != nil {
		logger.Error(err.Error())
		return false, 0
	}

	return true, int32(i)
}
