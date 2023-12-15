package category_definer

import "strings"

type Category string

var (
	Car                = Category("Carro")
	BarsAndRestaurants = Category("Bares e restaurantes")
	Food               = Category("Alimentação")
	Pharmacy           = Category("Farmácia")
	IFood              = Category("IFood")
	Market             = Category("Mercado")
	Pet                = Category("Pet")
	Uber               = Category("Uber")
)

var categoryKeyWords = map[Category][]string{
	Car:                {"posto", "autopost", "conectcar", "estacion", "meuestar"},
	BarsAndRestaurants: {"boteco", "bar", "restaurante", "churrascaria"},
	Food:               {"saintger", "saint ger", "coffee", "café"},
	Pharmacy:           {"panvel", "raia"},
	IFood:              {"ifood"},
	Market:             {"festval", "super beal", "mercado", "supermercado", "market4u"},
	Pet:                {"cobasi"},
	Uber:               {"uber"},
}

func GetCategoryFromDescription(description string) Category {
	for category, keys := range categoryKeyWords {
		for i := 0; i < len(keys); i++ {
			if strings.Contains(strings.ToLower(description), strings.ToLower(keys[i])) {
				return category
			}
		}
	}

	return ""
}
