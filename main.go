package main

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func getHtml(url string) io.Reader {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	//defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return charmap.ISO8859_1.NewDecoder().Reader(res.Body)
	//return res.Body
}

func getHtmlFake() io.Reader {
	file, err := os.Open("test/fixtures/19.99.html")
	if err != nil {
		panic(err)
	}
	r := charmap.ISO8859_1.NewDecoder().Reader(file)
	return r
}

func main() {

	fromPrice := os.Args[2]
	for {
		checkPrice(fromPrice)
		time.Sleep(10 * time.Minute)
	}
}

func checkPrice(fromPrice string) *goquery.Selection {
	//doc, err := goquery.NewDocumentFromReader(getHtmlFake())
	doc, err := goquery.NewDocumentFromReader(getHtml(os.Args[1]))
	if err != nil {
		log.Fatal(err)
	}
	return doc.Find(".product-summary").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Find(".product-title").Text()
		title2 := strings.ReplaceAll(title, "\n", "")

		r := regexp.MustCompile(`\s+`)
		title3 := r.ReplaceAllString(title2, " ")
		title4 := strings.TrimSpace(title3)

		price := s.Find(".product-price").Text()
		price2 := r.ReplaceAllString(price, " ")
		price3 := strings.TrimSpace(price2)
		price4 := strings.Split(price3, " ")[0]
		if fromPrice == price4 {
			log.Print("En attente d'un changement de prix...")
		} else {
			log.Printf("!! Le prix a changé !! ---> %s => %s\n", price4, title4)
			os.Exit(0)
		}
		return // only first item
	})
}
