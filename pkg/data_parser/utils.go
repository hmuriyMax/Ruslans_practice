package data_parser

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// getItemsCnt для того, чтобы узнать общее количество товаров данного вида
func getItemsCnt(url string) (int, error) {
	body := MakeRequest(url, nil, 0)
	cnt := cntItemsFromOnePage{}
	err := json.Unmarshal(body, &cnt)
	if err != nil {
		log.Fatal(err)
	}
	return cnt.NumFound, nil
}

func getAllItemsFromOnePage(request string) (page, error) {
	body := MakeRequest(request, nil, 0)
	var pg page
	//Todo доработать метод (https://stackoverflow.com/questions/42377989/unmarshal-json-array-of-arrays-in-go)
	if err := json.Unmarshal(body, &pg); err != nil {
		log.Fatal(err)
	}
	return pg, nil
}

func getPriceOfOneItem(itemId string) (price, error) {
	const baseUrl = "https://api.waters.com/waters-product-exp-api-v1/api/products/prices?customerNumber=anonymous&productNumber="
	siteUrl := baseUrl + itemId
	headers := map[string]string{
		"countryCode": "us",
		"channel":     "ECOMM",
		"language":    "en",
	}
	body := MakeRequest(siteUrl, headers, 0)
	var pr []price
	if err := json.Unmarshal(body, &pr); err != nil {
		log.Fatal(err)
	}
	return pr[0], nil
}

func getAvailabilityOfOneItem(itemId string) (availability, error) {
	//Todo узнать про страну запроса (от этого зависит результат)
	siteUrl := fmt.Sprintf("https://prodservices.waters.com/api/waters/product/v1/availability/%s/US", itemId)
	body := MakeRequest(siteUrl, nil, 0)
	var av availability
	if err := json.Unmarshal(body, &av); err != nil {
		log.Fatal(err)
	}
	return av, nil
}

func MakeRequest(siteURL string, headers map[string]string, timeout int) []byte {
	body := io.Reader(nil)
	req, err := http.NewRequest(http.MethodGet, siteURL, body)
	if err != nil {
		log.Fatalln(err)
	}
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	// Use the default timeout if the timeout parameter isn't configured.
	reqTimeout := 10 * time.Second
	if timeout != 0 {
		reqTimeout = time.Duration(timeout) * time.Second
	}

	// Use default http Client.
	httpClient := &http.Client{
		Transport:     http.DefaultTransport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       reqTimeout,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalln(resp.Status)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return respBody
}