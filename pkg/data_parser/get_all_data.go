package data_parser

import (
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
)

func GetAllData(lgr *log.Logger, contentType string, countriesList []string) ([]product, error) {
	itemsCnt, err := getItemsCnt(fmt.Sprintf(url, strings.ToLower(contentType), contentType, 1))
	if err != nil {
		return nil, err
	}
	products := make([]product, itemsCnt)

	// TODO: выяснить, как работает арифметика при разных типах
	pageCnt := int(math.Ceil(float64(itemsCnt) / ItemsOnPage))
	lgr.Printf("%d items on %d pages found. Lets parse!", itemsCnt, pageCnt)

	wg := sync.WaitGroup{}
	wg.Add(itemsCnt)
	mu := sync.Mutex{}
	for i := 0; i < pageCnt; i++ {
		lgr.Printf("Page %3.d parsing...", i+1)
		pg, err := getAllItemsFromOnePage(fmt.Sprintf(url, strings.ToLower(contentType), contentType, i+1))
		if err != nil {
			log.Printf("Error parsing page %d: %v", i, err)
			continue
		}
		for j := 0; j < len(pg.Items); j++ {
			// Создам горутину для отдельного товара
			go func(i, j int, pg *page) {
				defer wg.Done()
				itemId := pg.Items[j].Id

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
					Id:           pg.Items[j].Id,
					Title:        pg.Items[j].Title,
					Url:          pg.Items[j].Url,
					Availability: avlb,
					Price:        price.Val,
				}
				mu.Unlock()
			}(i, j, &pg)
		}
	}
	wg.Wait()
	return products, nil
}
