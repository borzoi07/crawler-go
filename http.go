package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "FriendlyCrawler/1.0")

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("statusErr: error status code %d", res.StatusCode)
	}
	if !strings.Contains(res.Header.Get("Content-Type"), "text/html") {
		return "", fmt.Errorf("contentErr: invalid content-type: %s", res.Header.Get("Content-Type"))
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	// make sure it's crawling the same domain
	// TODO: detecting a substring may be unpredictable, consider using 'net/url' or 'path'
	if !strings.Contains(rawCurrentURL, rawBaseURL) {
		fmt.Printf("--DEBUG: detected different domain: %s\n", rawCurrentURL)
		return
	}

	normalCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("> returned err: %v\n", err)
		return
	}

	// if it already crawled this return
	if _, ok := pages[normalCurrentURL]; ok {
		pages[normalCurrentURL]++
		return
	} else {
		pages[normalCurrentURL] = 1
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("> returned err: %v\n", err)
		return
	}
	fmt.Printf("--DEBUG: got html from %q\n", rawCurrentURL)

	pageData := extractPageData(html, rawCurrentURL)

	for _, url := range pageData.OutgoingLinks {
		crawlPage(rawBaseURL, url, pages)
	}
}
