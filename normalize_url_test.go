package main

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		expected  string
		keepQuery bool
	}{
		{
			name:     "remove scheme",
			inputURL: "https://www.youtube.com/some/path",
			expected: "www.youtube.com/some/path",
		},
		// more for edge cases
		{
			name:     "remove extra slash",
			inputURL: "https://www.youtube.com/some/path/",
			expected: "www.youtube.com/some/path",
		},
		{
			name:      "keep query",
			inputURL:  "https://www.youtube.com/some/path?foo=bar&foo1=bar1",
			expected:  "www.youtube.com/some/path?foo=bar&foo1=bar1",
			keepQuery: true,
		},
		{
			name:     "remove query",
			inputURL: "https://www.youtube.com/some/path?foo=bar&foo1=bar1",
			expected: "www.youtube.com/some/path",
		},
		{
			name:     "remove http and userinfo and port, make lowercase",
			inputURL: "http://foo:1234@wWW.yOUtUBE.COM:1111/some/path",
			expected: "www.youtube.com/some/path",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var actual string
			var err error
			if tc.keepQuery {
				actual, err = normalizeURL(tc.inputURL, true)
			} else {
				actual, err = normalizeURL(tc.inputURL)
			}
			if err != nil {
				t.Errorf("Test %d - '%s' FAIL: unexpected error: %v", i, tc.name, err)
			}
			if actual != tc.expected {
				t.Errorf("Test %d - %s FAIL: expected URL: %s, actual: %s", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func _toLowerASCIIsingleLoop(str string) string {
	res := []byte(str)
	for i := 0; i < len(str); i++ { // A: 65 , Z: 90 ; a: 97, z: 122
		c := res[i]
		if c >= 0x41 && c <= 0x5A {
			res[i] = c + 0x20
		}
	}

	return string(res)
}

func _toLowerUTF8(str string) string {
	var res []rune = make([]rune, 0, len(str))

	for _, r := range str {
		if r >= 0x41 && r <= 0x5A {
			r += 0x20
		}
		res = append(res, r)
	}

	return string(res)
}

const (
	lowercase = "www.testhostname.com"
	uppercase = "www.TestHostName.COM"
)

func BenchmarkOneLoop_Lowercase(b *testing.B) {
	for b.Loop() {
		_toLowerASCIIsingleLoop(lowercase)
	}
}

func BenchmarkTwoLoops_Lowercase(b *testing.B) {
	for b.Loop() {
		toLowerASCII(lowercase)
	}
}

func BenchmarkOneLoop_Uppercase(b *testing.B) {
	for b.Loop() {
		_toLowerASCIIsingleLoop(uppercase)
	}
}

func BenchmarkTwoLoops_Uppercase(b *testing.B) {
	for b.Loop() {
		toLowerASCII(uppercase)
	}
}

func BenchmarkTwoLoops_UTF_Lowercase(b *testing.B) {
	for b.Loop() {
		_toLowerUTF8(lowercase)
	}
}

func BenchmarkTwoLoops_UTF_Uppercase(b *testing.B) {
	for b.Loop() {
		_toLowerUTF8(uppercase)
	}
}

/*
func toLowerTEST(str string) string {
	var res []rune = make([]rune, 0, len(str))

	var p int
	for i := 0; i < len(str); i++ {
		// NOTE: breaks unicode characters
		r := rune(str[i])
		if r >= 0x41 && r <= 0x5A {
			r += 0x20

			if p < i {
				// NOTE: breaks unicode characters because this slices by byte indices but URLs shouldn't contain non-ASCII chars so it should be ok
				res = append(res, []rune(str[p:i])...)
			}
			res = append(res, r)

			p = i + 1
		}
	}
	if p < len(str) {
		res = append(res, []rune(str[p:])...)
	}

	return string(res)
}
*/
