package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const request = "https://prodservices.waters.com/api/waters/search/category_facet$shop:Shop&contenttype_facet$columns:Columns?isocode=en_US&page=1"

type cntItemsFromOnePage struct {
	NumFound int `json:"num_found"`
}

type item struct {
	title string
	id    int
}

type page struct {
	items []struct {
		it item
	}
}

func main() {
	fmt.Print(GetPagesCnt1(request))
}

func MakeRequest(request string) map[string]interface{} {
	resp, err := http.Get(request)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalln(resp.Status)
	}

	var data map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatalln(err)
	}
	return data
}

func MakeRequest1(request string) []byte {
	resp, err := http.Get(request)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalln(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

//Для того, чтобу узнать общее количество страниц товаров данного вида
func GetPagesCnt(request string) (int, error) {
	data := MakeRequest(request)
	cnt, ok := data["num_found"].(float64)
	if !ok {
		strErr := fmt.Sprintf("Cannon't convert %v to int", cnt)
		return -1, errors.New(strErr)
	}
	//Todo найти метод округления вверх
	cntInt := int(cnt/12 + 1)
	return cntInt, nil
}

func GetPagesCnt1(request string) (int, error) {
	body := MakeRequest1(request)
	items := cntItemsFromOnePage{}
	err := json.Unmarshal(body, &items)
	if err != nil {
		log.Fatal(err)
	}

	//Todo найти метод округления вверх
	cntInt := int(items.NumFound/12 + 1)
	return cntInt, nil
}

/*
func GetAllItemsFromOnePage(request string) int {
	data := MakeRequest(request)
	arr := data["documents"]
	pg := page{items: []struct{ it item }{}}
	for _, item := range arr {

	}
}*/
