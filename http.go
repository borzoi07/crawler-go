package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{} // will block if buffer limit is reached
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	// immediately return when maxPage limit is reached
	if cfg.getPageCount() >= cfg.maxPages {
		return
	}

	parsedCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("> returned err: %v\n", err)
		return
	}

	// make sure it's crawling the same domain
	if cfg.baseURL.Hostname() != parsedCurrentURL.Hostname() {
		// TODO: count external links and add it to the PageData struct
		fmt.Printf("--DEBUG: detected different domain: %s\n", rawCurrentURL)
		return
	}

	normalCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("> returned err: %v\n", err)
		return
	}

	// if it already crawled this increment and return
	if !cfg.addPageVisit(normalCurrentURL) {
		return
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("> returned err: %v\n", err)
		return
	}
	fmt.Printf("--DEBUG: got html from %q\n", rawCurrentURL)

	pageData := extractPageData(html, rawCurrentURL)

	// add page data
	cfg.mu.Lock()
	pageData.TimesVisited = 1
	cfg.pages[normalCurrentURL] = pageData
	cfg.mu.Unlock()

	for _, url := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(url)
	}
}

func (cfg *config) addPageVisit(normaliedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if pageData, ok := cfg.pages[normaliedURL]; ok {
		pageData.TimesVisited++
		cfg.pages[normaliedURL] = pageData
		return false
	}

	return true
}

func (cfg *config) getPageCount() int {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	return len(cfg.pages)
}
