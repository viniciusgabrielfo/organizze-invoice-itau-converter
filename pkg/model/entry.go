package model

import (
	"github.com/Rhymond/go-money"
	"github.com/viniciusgabrielfo/organizze-invoice-itau-converter/pkg/category_definer"
)

type Entry struct {
	Date        string
	Description string
	Category    category_definer.Category
	Value       *money.Money
}
