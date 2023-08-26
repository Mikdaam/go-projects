package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type seller struct {
	Name  string
	Phone string
}

type carAd struct {
	Name        string
	Price       int
	Url         string
	Year        int
	Mileage     int
	Fuel        string
	Gearbox     string
	Description string
	Pictures    []string
	Location    string
	Seller      *seller
}

func main() {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"),
	)

	// another collector to scrape ad details
	detailCollector := c.Clone()

	// create a slice to hold the data
	ads := make([]*carAd, 0, 200)

	c.OnHTML("nav[aria-label=\"pagination\"] > ul > li > a[title=\"Page suivante\"]", func(h *colly.HTMLElement) {
		link := h.Attr("href")
		if strings.HasSuffix(link, "p-2") {
			log.Println("Reached the last page, exiting!")
			return
		}
		c.Visit(h.Request.AbsoluteURL(h.Attr("href")))
	})

	c.OnHTML("div.styles_adCard__HQRFN.styles_classified__rnsg4", func(e *colly.HTMLElement) {
		ad := new(carAd)

		ad.Name = e.ChildText("p[data-qa-id=\"aditem_title\"]")
		ad.Url = e.ChildAttr("a[data-qa-id=\"aditem_container\"]", "href")

		priceStr := e.ChildText("p[data-test-id=\"price\"] span span")
		priceStr = strings.ReplaceAll(priceStr, "\u00a0", "") // Remove non-breaking spaces
		priceStr = strings.ReplaceAll(priceStr, "€", "")
		ad.Price, _ = strconv.Atoi(priceStr)

		params := e.ChildText("div[data-test-id=\"ad-params-light\"] span")
		paramsArray := strings.Split(params, " • ")
		if len(paramsArray) == 4 {
			ad.Year, _ = strconv.Atoi(paramsArray[0])
			ad.Mileage, _ = strconv.Atoi(strings.ReplaceAll(paramsArray[1], " km", ""))
			ad.Fuel = paramsArray[2]
			ad.Gearbox = paramsArray[3]
		}

		ad.Location = e.ChildText("p[class^=\"_2k43C\"]")

		ads = append(ads, ad)

		detailCollector.Visit(e.Request.AbsoluteURL(ad.Url))
	})

	detailCollector.OnHTML("section[data-qa-id=\"adview_spotlight_container\"] > div > div > div", func(e *colly.HTMLElement) {
		ad := ads[len(ads)-1]
		pictures := e.ChildAttrs("img", "src")
		ad.Pictures = append(ad.Pictures, pictures...)
	})

	detailCollector.OnHTML("div[data-qa-id=\"adview_description_container\"]", func(e *colly.HTMLElement) {
		ad := ads[len(ads)-1]
		description := e.ChildText("div > p")
		ad.Description = strings.TrimSpace(description)
	})

	/*detailCollector.OnHTML("div[class^=\"_1o09v\"]", func(e *colly.HTMLElement) {
		seller := new(seller)
		seller.Name = e.ChildText("a[class^=\"text-headline-2\"]")
		sellerUrl := e.ChildAttr("a[class^=\"text-headline-2\"]", "href")

		// Visit the seller's profile page to get the phone number
		detailCollector.Visit(e.Request.AbsoluteURL(sellerUrl))

		ads[len(ads)-1].Seller = seller
	})*/

	detailCollector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting Details page ", r.URL)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting ", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL: ", r.Request.URL, " failed with response: ", string(r.Body), "\nError: ", err)
	})

	detailCollector.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL: ", r.Request.URL, " failed with response: ", string(r.Body), "\nError: ", err)
	})

	c.OnScraped(func(r *colly.Response) {
		log.Println("Finished scraping ", r.Request.URL)

		data, err := json.MarshalIndent(ads, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		// Write the data to a file
		fmt.Println("Writing data to file")
		err = os.WriteFile("ads.json", data, 0644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done!")
	})

	c.Visit("https://www.leboncoin.fr/f/voitures/u_car_brand--TOYOTA")
}
