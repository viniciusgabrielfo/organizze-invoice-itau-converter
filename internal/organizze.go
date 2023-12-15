package internal

import (
	"fmt"

	"github.com/viniciusgabrielfo/organizze-invoice-itau-converter/pkg/model"
	"github.com/xuri/excelize/v2"
)

func GenerateOrganizzeXLXSSheet(entries []model.Entry) error {
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
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), entries[i].Category)
		f.SetCellFloat("Sheet1", fmt.Sprintf("D%d", row), entries[i].Value.AsMajorUnits(), 2, 32)

	}

	if err := f.SaveAs("fatura.xlsx"); err != nil {
		return err
	}

	return nil
}
