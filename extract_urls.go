package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL            string
	Heading        string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func extractPageData(html, pageURL string) PageData {
	heading := getHeadingFromHTML(html)
	fparagraph := getFirstParagraphFromHTML(html)

	baseURL, err := url.Parse(pageURL)
	if err != nil {
		return PageData{
			URL:            pageURL,
			Heading:        heading,
			FirstParagraph: fparagraph,
			OutgoingLinks:  nil,
			ImageURLs:      nil,
		}
	}

	links, err := getURLsFromHTML(html, baseURL)
	if err != nil {
		links = nil
	}
	imageURLs, err := getImagesFromHTML(html, baseURL)
	if err != nil {
		imageURLs = nil
	}

	return PageData{
		URL:            pageURL,
		Heading:        heading,
		FirstParagraph: fparagraph,
		OutgoingLinks:  links,
		ImageURLs:      imageURLs,
	}
}

// return <h1> tag if present, <h2> as fallback,
// if not found returns empty string
func getHeadingFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Printf("err: Couldn't parse raw html string: %s ", err)
		return ""
	}

	h := doc.Find("h1, h2").First().Text()
	return strings.TrimSpace(h)
}

func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Printf("err: Couldn't parse raw html string: %s ", err)
		return ""
	}

	var p string
	main := doc.Find("main")
	if main.Length() > 0 {
		p = main.Find("p").First().Text()
	} else {
		p = doc.Find("p").First().Text()
	}

	return strings.TrimSpace(p)
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var res []string

	// for each <a href> run the following function
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if !exists || strings.TrimSpace(link) == "" {
			return
		}
		lURL, err := url.Parse(link)
		if err != nil {
			fmt.Printf("err: returned error while parsing %q, %v\n", link, err)
			return
		}

		abs := baseURL.ResolveReference(lURL)
		// below if block may not be required
		if !abs.IsAbs() {
			abs.Scheme = baseURL.Scheme
			if abs.Hostname() == "" {
				abs.Host = baseURL.Host
			}
		}

		res = append(res, abs.String())
	})

	return res, nil
}

// TODO REFACTOR: create a helper function to remove duplicate code between getImagesFromHTML and getURLsFromHTML

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var res []string

	// for each <img src> run the following function
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		link, exists := s.Attr("src")
		if !exists || strings.TrimSpace(link) == "" {
			return
		}
		lURL, err := url.Parse(link)
		if err != nil {
			fmt.Printf("err: returned error while parsing %q, %v\n", link, err)
			return
		}

		abs := baseURL.ResolveReference(lURL)
		// below if block may not be required
		if !abs.IsAbs() {
			abs.Scheme = baseURL.Scheme
			if abs.Hostname() == "" {
				abs.Host = baseURL.Host
			}
		}

		res = append(res, abs.String())
	})

	return res, nil
}
