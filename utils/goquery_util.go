package utils

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

// GetHrefFromUrl 从指定的uri中获取指定selector，指定属性的值
func GetHrefFromUrl(uri string, selector string, attr string) (attrv string, err error) {
	res, err := http.Get(uri)
	if err != nil {
		log.Printf("Failed to get Amazon seller center url: %v", err)
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("Response failed with status code: %d and\n", res.StatusCode)
		return "", errors.New(fmt.Sprintf("Response failed with status code: %d and", res.StatusCode))
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Printf("goquery load html failed: %v", err)
		return "", err
	}

	// Find the review items
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		title := s.Text()
		attrv, _ = s.Attr(attr)
		fmt.Printf("Review %d: %s, %s: %s\n", i, title, attr, attrv)
	})

	if attrv != "" {
		return attrv, nil
	}
	return "", errors.New(fmt.Sprintf("Can't find selector(%s)'s %s.'", selector, attr))
}
