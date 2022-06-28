package data_parser

import (
	"fmt"
	"math"
	"strings"
)

//TODO: добавить адекватную обработку ошибок

func GetAllData(contentType string) ([]product, error) {
	num, _ := getItemsCnt(fmt.Sprintf(url, strings.ToLower(contentType), contentType, 1))
	products := make([]product, num)
	// TODO: выяснить, как работает арифметика при разных типах
	pageCnt := int(math.Ceil(float64(num) / 12))
	for i := 0; i < pageCnt; i++ {
		page, _ := getAllItemsFromOnePage(fmt.Sprintf(url, strings.ToLower(contentType), contentType, i+1))
		for j := 0; j < len(page.Items); j++ {
			itemId := page.Items[j].Id
			//Todo распараллелить
			price, _ := getPriceOfOneItem(itemId)
			availability, _ := getAvailabilityOfOneItem(itemId)
			products[i*ItemsOnPage+j] = product{
				Id:           page.Items[j].Id,
				Title:        page.Items[j].Title,
				Availability: availability.Status,
				Price: struct {
					Currency string
					Value    float64
				}(price.Val),
			}
		}
	}
	return products, nil
}
