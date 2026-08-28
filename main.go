package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("> no website provided")
		os.Exit(1)
	} else if len(args) > 4 {
		fmt.Println("> too many arguments provided")
		os.Exit(1)
	}
	BASE_URL := args[1]
	var maxConcurrency int
	var maxPages int

	if len(args) < 3 {
		maxConcurrency = 5
		fmt.Println("> setting default maxConcurrecy:", maxConcurrency)
	} else {
		var err error
		maxConcurrency, err = strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("> invalid maxConcurrency argument:", err)
			os.Exit(1)
		}
	}
	if len(args) < 4 {
		maxPages = 250
		fmt.Println("> setting default maxPages:", maxPages)
	} else {
		var err error
		maxPages, err = strconv.Atoi(args[3])
		if err != nil {
			fmt.Println("> invalid maxPages argument:", err)
			os.Exit(1)
		}
	}

	fmt.Printf("> starting crawl of: %s\n\n", BASE_URL)

	baseURL, err := url.Parse(BASE_URL)
	if err != nil {
		fmt.Println("> couldn't parse the base url:", err)
		os.Exit(1)
	}

	var cfg *config = &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency), // at most n goroutines at once
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	start := time.Now()

	cfg.wg.Add(1)
	go cfg.crawlPage(BASE_URL)
	cfg.wg.Wait() // waitgroup ensures that main doesn't exit until all goroutines are done

	elapsedTime := time.Since(start)
	fmt.Printf("\n- Finished in %q\n", elapsedTime)
	fmt.Printf("\n- Got '%d' pages\n", len(cfg.pages))
	for k, v := range cfg.pages {
		fmt.Printf(" • %d: %s\n", v.TimesVisited, k)
	}

	err = writeJSONReport(cfg.pages, "report.json")
	if err != nil {
		fmt.Println("> writeJSONReport returned err:", err)
		os.Exit(1)
	}
}
