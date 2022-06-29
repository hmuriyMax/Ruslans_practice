package data_parser

import (
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
)

func GetAllData(lo *log.Logger, contentType string, countriesList []string) ([]product, error) {
	itemsCnt, err := getItemsCnt(fmt.Sprintf(url, strings.ToLower(contentType), contentType, 1))
	if err != nil {
		return nil, err
	}
	products := make([]product, itemsCnt)

	// TODO: выяснить, как работает арифметика при разных типах
	pageCnt := int(math.Ceil(float64(itemsCnt) / ItemsOnPage))
	lo.Printf("%d items on %d pages found. Lets parse!", itemsCnt, pageCnt)

	wg := sync.WaitGroup{}
	wg.Add(itemsCnt)
	mu := sync.Mutex{}
	for i := 0; i < pageCnt; i++ {
		lo.Printf("Page %3.d parsing...", i+1)
		page, err := getAllItemsFromOnePage(fmt.Sprintf(url, strings.ToLower(contentType), contentType, i+1))
		if err != nil {
			log.Printf("Error parsing page %d: %v", i, err)
			continue
		}
		for j := 0; j < len(page.Items); j++ {
			// Создам горутину для отдельного товара
			go func(i, j int) {
				defer wg.Done()
				itemId := page.Items[j].Id

				price, err := getPriceOfOneItem(itemId)
				if err != nil {
					log.Printf("Error parsing price of item %s: %v", itemId, err)
					return
				}

				avlb := make(map[string]bool)
				for _, el := range countriesList {
					availability, err := getAvailabilityOfOneItem(itemId, el)
					if err != nil {
						log.Printf("Error parsing availability in %s of item %s: %v", el, itemId, err)
					} else {
						avlb[el] = availability.Status == "IN_STOCK"
					}
				}

				mu.Lock()
				products[i*ItemsOnPage+j] = product{
					Id:           page.Items[j].Id,
					Title:        page.Items[j].Title,
					Url:          page.Items[j].Url,
					Availability: avlb,
					Price:        price.Val,
				}
				mu.Unlock()
			}(i, j)
		}
	}
	wg.Wait()
	return products, nil
}
