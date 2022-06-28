package data_parser

const ItemsOnPage = 99
const url = "https://prodservices.waters.com/api/waters/search/category_facet$%s:%s?isocode=en_US&page=%d&rows=99"

// TODO: это URL, который извлекает все товары

type cntItemsFromOnePage struct {
	NumFound int `json:"num_found"`
}

//TODO: добавить поле с типом товара
type page struct {
	Items []struct {
		Title string `json:"title"`
		Id    string `json:"skucode"`
	} `json:"documents"`
}

type price struct {
	Val struct {
		Currency string  `json:"currencyCode"`
		Value    float64 `json:"value"`
	} `json:"basePrice"`
}

// TODO: доступность по странам (иногда в зависимости от страны меняется доступность)
type availability struct {
	Status string `json:"productStatus"`
}

type product struct {
	Title        string
	Id           string
	Availability string
	Price        struct {
		Currency string
		Value    float64
	}
}
