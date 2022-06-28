package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"webScraper/pkg/utils"
)

const url = "https://prodservices.waters.com/api/waters/search/category_facet$shop:Shop&contenttype_facet$columns:Columns?isocode=en_US&page=1"

type cntItemsFromOnePage struct {
	NumFound int `json:"num_found"`
}

type page struct {
	Items []struct {
		Title string `json:"title"`
		Id    string `json:"skucode"`
	} `json:"documents"`
}

//Todo Возможно понадобиться изменить структуру для сбора цены (см. ответ запроса цены https://api.waters.com/waters-product-exp-api-v1/api/products/prices?customerNumber=anonymous&productNumber=186010094)

type price struct {
	Val struct {
		Currency string  `json:"currencyCode"`
		Value    float64 `json:"value"`
	} `json:"basePrice"`
}

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

func main() {
	//fmt.Print(GetPriceOfOneItem("405000981"))
	start := time.Now()
	result, _ := GetAllDataFromOnePage(url)
	duration := time.Since(start)
	fmt.Println(duration)
	fmt.Println(duration.Nanoseconds())
	fmt.Print(result)
}

//Для того, чтобу узнать общее количество страниц товаров данного вида
func GetItemsCnt(url string) (int, error) {
	body := utils.MakeRequest(url, nil, 0)
	cnt := cntItemsFromOnePage{}
	err := json.Unmarshal(body, &cnt)
	if err != nil {
		log.Fatal(err)
	}
	return cnt.NumFound, nil
}

func GetPagesCnt(url string) (int, error) {
	//Todo найти метод округления вверх
	//Todo обработка ошибок
	cnt, _ := GetItemsCnt(url)
	cntInt := int(cnt/12 + 1)
	return cntInt, nil
}

func GetAllItemsFromOnePage(request string) (page, error) {
	body := utils.MakeRequest(request, nil, 0)
	var pg page
	//Todo доработать метод (https://stackoverflow.com/questions/42377989/unmarshal-json-array-of-arrays-in-go)
	if err := json.Unmarshal(body, &pg); err != nil {
		log.Fatal(err)
	}
	return pg, nil
}

func GetPriceOfOneItem(itemId string) (price, error) {
	const baseUrl = "https://api.waters.com/waters-product-exp-api-v1/api/products/prices?customerNumber=anonymous&productNumber="
	siteUrl := baseUrl + itemId
	headers := map[string]string{
		"countryCode": "us",
		"channel":     "ECOMM",
		"language":    "en",
	}
	body := utils.MakeRequest(siteUrl, headers, 0)
	var pr []price
	if err := json.Unmarshal(body, &pr); err != nil {
		log.Fatal(err)
	}
	return pr[0], nil
}

func GetAvailabilityOfOneItem(itemId string) (availability, error) {
	//Todo узнать про страну запроса (от этого зависит результат)
	siteUrl := fmt.Sprintf("https://prodservices.waters.com/api/waters/product/v1/availability/%s/US", itemId)
	body := utils.MakeRequest(siteUrl, nil, 0)
	var av availability
	if err := json.Unmarshal(body, &av); err != nil {
		log.Fatal(err)
	}
	return av, nil
}

func GetAllData(request string) ([]product, error) {
	products := []product{}
	page, _ := GetAllItemsFromOnePage(request)
	for i := 0; i < len(page.Items); i++ {
		itemId := page.Items[i].Id
		//Todo распараллелить
		price, _ := GetPriceOfOneItem(itemId)
		availability, _ := GetAvailabilityOfOneItem(itemId)
		products = append(products,
			product{
				Id:           page.Items[i].Id,
				Title:        page.Items[i].Title,
				Availability: availability.Status,
				Price: struct {
					Currency string
					Value    float64
				}(price.Val),
			})
	}
	return products, nil
}

func GetAllData(request string) ([]product, error) {

}
