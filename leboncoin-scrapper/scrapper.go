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
	Name  string `json:"nom"`
	Phone string `json:"telephone"`
	IsPro bool   `json:"pro"`
}

type carAd struct {
	Name        string   `json:"nom"`
	Price       int      `json:"prix"`
	Url         string   `json:"url"`
	Year        int      `json:"annee"`
	Mileage     int      `json:"kilometrage"`
	Fuel        string   `json:"carburant"`
	Gearbox     string   `json:"boite_de_vitesse"`
	Description string   `json:"description"`
	Pictures    []string `json:"photos"`
	Location    string   `json:"localisation"`
	Seller      *seller  `json:"vendeur"`
}

func main() {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"),
		//colly.Async(true),
	)

	// Define headers
	headers := map[string]string{
		"accept":          "application/json",
		"accept-language": "en-FR,en;q=0.9,fr-FR;q=0.8,fr;q=0.7,en-GB;q=0.6,en-US;q=0.5",
		"authorization":   "Bearer ", // Replace with your actual access token
		"content-type":    "application/x-www-form-urlencoded",
	}

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

	detailCollector.OnHTML("aside > div[class^=\"_1o09v\"] > section > div:nth-child(1)", func(e *colly.HTMLElement) {
		seller := new(seller)
		seller.Name = e.ChildText("div:nth-child(1) > div:nth-child(2) > a")

		//#aside > div._1o09v > section > div.sc-11fe9401-0.kiGxvX > div.src__Box-sc-10d053g-0.eWhWEg > a
		if seller.Name == "" {
			link := e.ChildAttr("div:nth-child(2) > a", "href")
			seller.Name = e.ChildText("div:nth-child(1) > div:nth-child(2) > a")

			fmt.Println("link: ", link)
		}

		// Check if the seller is a pro
		proText := e.ChildText("div.sc-11fe9401-0.kiGxvX > div.src__Box-sc-10d053g-0.bfCalR > span")
		if proText != "" {
			seller.IsPro = true
		}

		// Get the phone number
		/*detailCollector.Post(
			"https://api.leboncoin.fr/api/utils/phonenumber.json",
			map[string]string{"app_id": "leboncoin_web_utils", "list_id": "2400112289", "text": "1"},
		)*/

		ads[len(ads)-1].Seller = seller
	})

	detailCollector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting Details page ", r.URL)
		if r.Method == "POST" {
			// Set the headers
			for key, value := range headers {
				r.Headers.Set(key, value)
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting ", r.URL)
	})

	detailCollector.OnResponse(func(r *colly.Response) {
		log.Println("Visited (Probably Captcha 😭", r.Request.URL)
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

	// Wait until threads are finished
	/*detailCollector.Wait()

	c.Wait()*/
}
