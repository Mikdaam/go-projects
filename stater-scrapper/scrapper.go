package main

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
)

type Book struct {
	Link         string `json:"link"`
	Name         string `json:"name"`
	Price        string `json:"price"`
	Availability string `json:"availability"`
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector(colly.Async(true))

	// Array of books
	var books []Book

	c.OnHTML("div.side_categories li ul li a", func(h *colly.HTMLElement) {
		link := h.Attr("href")
		err := c.Visit(h.Request.AbsoluteURL(link))
		if err != nil {
			return
		}
	})

	c.OnHTML("li.next a", func(h *colly.HTMLElement) {
		err := c.Visit(h.Request.AbsoluteURL(h.Attr("href")))
		if err != nil {
			return
		}
	})

	c.OnHTML("article.product_pod", func(h *colly.HTMLElement) {
		book := Book{
			Link:         h.ChildAttr("h3 a", "href"),
			Name:         h.ChildAttr("h3 a", "title"),
			Price:        h.ChildText("div.product_price p.price_color"),
			Availability: h.ChildText("div.product_price p.instock.availability"),
		}
		books = append(books, book)
	})

	c.OnRequest(func(r *colly.Request) {
		println("Visiting", r.URL.String())
	})

	/*err := c.Visit("https://books.toscrape.com/catalogue/page-1.html")
	if err != nil {
		return
	}*/

	err := c.Visit("https://books.toscrape.com/catalogue/category/books/travel_2/index.html")
	if err != nil {
		return
	}

	c.Wait()

	data, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// write the json data to file
	err = ioutil.WriteFile("books.json", data, 0644)
}
