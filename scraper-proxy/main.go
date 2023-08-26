package main

import (
	"bytes"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"log"
)

func main() {
	c := colly.NewCollector(colly.AllowURLRevisit())

	// rotate two proxies
	rp, err := proxy.RoundRobinProxySwitcher(
		"https://181.129.74.58:40667",
		"https://144.49.99.190:8080",
	)
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)

	c.OnResponse(func(r *colly.Response) {
		log.Printf("%s\n", bytes.Replace(r.Body, []byte("\n"), nil, -1))
	})

	for i := 0; i < 5; i++ {
		err := c.Visit("https://httpbin.org/ip")
		if err != nil {
			log.Println(err)
		}
	}
}
