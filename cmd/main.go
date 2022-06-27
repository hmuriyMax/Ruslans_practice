package main

import (
	"encoding/json"
	"fmt"
	"log"
	"webScraper/pkg/utils"
)

const url = "https://prodservices.waters.com/api/waters/search/category_facet$shop:Shop&contenttype_facet$columns:Columns?isocode=en_US&page=1"

type cntItemsFromOnePage struct {
	NumFound int `json:"num_found"`
}

type Item struct {
	Title string
	Id    string
}

type page struct {
	Items []Item `json:"documents"`
}

//Todo Возможно понадобиться изменить структуру для сбора цены (см. ответ запроса цены https://api.waters.com/waters-product-exp-api-v1/api/products/prices?customerNumber=anonymous&productNumber=186010094)

type price struct {
	value struct {
		basePrice    float64
		currencyCode string
	}
}

func (it *Item) UnmarshalJSON(data []byte) error {
	var v []interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Printf("Error whilde decoding %v\n", err)
		return err
	}
	it.Title = v[2].(string)
	it.Id = v[6].(string)
	return nil
}

func main() {
	fmt.Print(GetPriceOfOneItem("186010094"))
}

//Для того, чтобу узнать общее количество страниц товаров данного вида
func GetPagesCnt(url string) (int, error) {
	body := utils.MakeRequest(url, nil, 0)
	cnt := cntItemsFromOnePage{}
	err := json.Unmarshal(body, &cnt)
	if err != nil {
		log.Fatal(err)
	}

	//Todo найти метод округления вверх
	//Todo обработка ошибок
	cntInt := int(cnt.NumFound/12 + 1)
	return cntInt, nil
}

/*
func GetAllItemsFromOnePage(request string) ([]Item, error) {
	body := MakeRequest(request)
	var pg page
	//Todo доработать метод (https://stackoverflow.com/questions/42377989/unmarshal-json-array-of-arrays-in-go)
	if err := json.Unmarshal(body, &pg); err != nil {
		log.Fatal(err)
	}
	return pg.Items, nil
}*/

func GetPriceOfOneItem(itemId string) (price, error) {
	const baseUrl = "https://api.waters.com/waters-product-exp-api-v1/api/products/prices?customerNumber=anonymous&productNumber="
	siteUrl := baseUrl + itemId
	headers := map[string]string{
		"countryCode": "us",
		"channel":     "ECOMM",
		"language":    "en",
	}
	body := utils.MakeRequest(siteUrl, headers, 0)
	var pr price
	if err := json.Unmarshal(body, &pr); err != nil {
		log.Fatal(err)
	}
	return pr, nil
}
