package main

import (
	"net/url"
	"path"
)

func normalizeURL(rawURL string, includeQuery ...bool) (string, error) {
	URL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// userinfo := URL.User.String()
	// if userinfo != "" {
	// 	userinfo += "@"
	// }

	// domains are case-insensitive
	hostname := toLowerASCII(URL.Hostname())

	var query string = ""
	if len(includeQuery) > 0 && includeQuery[0] == true {
		query = URL.Query().Encode()
		if query != "" {
			query = "?" + query
		}
	}

	URL.Path = path.Clean(URL.Path)

	return hostname + URL.EscapedPath() + query, nil
}

func toLowerASCII(str string) string {
	// most of the time domains are going to be lowercase
	var hasUpper bool
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 0x41 && c <= 0x5A {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		return str
	}

	res := []byte(str)
	for i := 0; i < len(str); i++ { // A: 65 , Z: 90 ; a: 97, z: 122
		c := res[i]
		if c >= 0x41 && c <= 0x5A {
			res[i] = c + 0x20
		}
	}

	return string(res)
}
