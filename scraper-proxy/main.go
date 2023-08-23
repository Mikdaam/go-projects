package main

import (
	"bytes"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"log"
)

func main() {
	c := colly.NewCollector(colly.AllowURLRevisit())

	// rotate two socks5 proxies
	rp, err := proxy.RoundRobinProxySwitcher("socks5://162.223.94.163:80", "socks5://103.49.202.252:80")
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)

	c.OnResponse(func(r *colly.Response) {
		log.Printf("%s\n", bytes.Replace(r.Body, []byte("\n"), nil, -1))
	})

	for i := 0; i < 8; i++ {
		err := c.Visit("https://httpbin.org/ip")
		if err != nil {
			log.Println(err)
		}
	}
}
