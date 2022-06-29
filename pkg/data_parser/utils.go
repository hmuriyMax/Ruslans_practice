package data_parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// getItemsCnt для того, чтобы узнать общее количество товаров данного вида
func getItemsCnt(url string) (int, error) {
	body, err := MakeRequest(url, nil, 0)
	if err != nil {
		return 0, err
	}
	cnt := cntItemsFromOnePage{}
	err = json.Unmarshal(body, &cnt)
	if err != nil {
		return 0, err
	}
	return cnt.NumFound, nil
}

func getAllItemsFromOnePage(request string) (page, error) {
	body, err := MakeRequest(request, nil, 0)
	if err != nil {
		return page{}, err
	}
	var pg page
	if err := json.Unmarshal(body, &pg); err != nil {
		return page{}, err
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
	body, err := MakeRequest(siteUrl, headers, 0)
	if err != nil {
		return price{}, err
	}
	var pr []price
	if err := json.Unmarshal(body, &pr); err != nil {
		return price{}, err
	}
	return pr[0], nil
}

func getAvailabilityOfOneItem(itemId string, countryCode string) (availability, error) {
	//Todo узнать про страну запроса (от этого зависит результат)
	siteUrl := fmt.Sprintf("https://prodservices.waters.com/api/waters/product/v1/availability/%s/%s", itemId, countryCode)
	body, err := MakeRequest(siteUrl, nil, 0)
	if err != nil {
		return availability{}, err
	}
	var av availability
	if err := json.Unmarshal(body, &av); err != nil {
		return availability{}, err
	}
	return av, nil
}

func MakeRequest(siteURL string, headers map[string]string, timeout int) ([]byte, error) {
	body := io.Reader(nil)
	req, err := http.NewRequest(http.MethodGet, siteURL, body)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Response is %d", resp.StatusCode))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func Print(logger *log.Logger, prod []product) {
	for _, el := range prod {
		logger.Printf(StructFormat+"\n", el.Id, el.Title, el.Price.Value, el.Price.Currency, el.Url)
		for countryCode, avlbty := range el.Availability {
			logger.Printf("\tAVAILABLE IN %s:  %v\n", countryCode, avlbty)
		}
		logger.Println()
	}
}
