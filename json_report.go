package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

func writeJSONReport(pages map[string]PageData, filename string) error {
	if len(pages) == 0 {
		fmt.Println("> there is no page data")
		return nil
	}

	keys := make([]string, 0, len(pages))
	for k := range pages {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pageData := make([]PageData, 0, len(pages))
	for _, key := range keys {
		pageData = append(pageData, pages[key])
	}

	jsonPageData, err := json.MarshalIndent(pageData, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, jsonPageData, 0644)
	if err != nil {
		return err
	}

	fmt.Println("\n> successfully created a json report file:", filename)
	return nil
}
