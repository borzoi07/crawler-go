package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestExtractPageData(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://crawler-test.com",
		Heading:        "Test Title",
		FirstParagraph: "This is the first paragraph.",
		OutgoingLinks:  []string{"https://crawler-test.com/link1"},
		ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestExtractPageData2nd(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body>
		<h3>Test Title</h3>
        <a href="https://google.com/link1">Link 1</a>
		<h2>Some Header</h2>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://crawler-test.com",
		Heading:        "Some Header",
		FirstParagraph: "",
		OutgoingLinks:  []string{"https://google.com/link1"},
		ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestGetHeadingFromHTMLBasic(t *testing.T) {
	inputBody := "<html><body><h1>Test Title</h1></body></html>"
	actual := getHeadingFromHTML(inputBody)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	inputBody := `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "Main paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetIMGsAndURLs(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		inputBody string
		expected  []string
		fn        rune
	}{
		{
			name:      "get url absolute",
			input:     "https://crawler-test.com",
			inputBody: `<html><body><a href="https://crawler-test.com"><span>test absolute</span></a></body></html>`,
			expected:  []string{"https://crawler-test.com"},
			fn:        'U',
		},
		{
			name:      "get image relative",
			input:     "https://crawler-test.com",
			inputBody: `<html><body><img src="/logo.png" alt="Logo"></body></html>`,
			expected:  []string{"https://crawler-test.com/logo.png"},
			fn:        'I',
		},
		{
			name:      "get url relative",
			input:     "https://crawler-test.com",
			inputBody: `<html><body><a href="/homepage"><span>test relative</span></a></body></html>`,
			expected:  []string{"https://crawler-test.com/homepage"},
			fn:        'U',
		},
		{
			name:      "get image absolute",
			input:     "https://crawler-test.com",
			inputBody: `<html><body><img src="https://crawler-test.com/logo.png" alt="Logo"></body></html>`,
			expected:  []string{"https://crawler-test.com/logo.png"},
			fn:        'I',
		},
		{
			name:      "get image absolute foreign url",
			input:     "https://crawler-test.com",
			inputBody: `<html><body><img src="https://some-foreign-url/logo.png" alt="Logo"></body></html>`,
			expected:  []string{"https://some-foreign-url/logo.png"},
			fn:        'I',
		},
		{
			name:      "get url absolute foreign url",
			input:     "https://crawler-test.com",
			inputBody: `<html><body><a href="https://some-foreign-url.com"><span>test absolute</span></a></body></html>`,
			expected:  []string{"https://some-foreign-url.com"},
			fn:        'U',
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.input)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}

			var actual []string
			switch tc.fn {
			case 'U':
				actual, err = getURLsFromHTML(tc.inputBody, baseURL)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			case 'I':
				actual, err = getImagesFromHTML(tc.inputBody, baseURL)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %d - %s FAIL: expected: %s, actual: %s", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestExtractURLs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		fn       rune
	}{
		{
			fn:   'P',
			name: "extract paragraph",
			input: `<html><body>
		<p>First paragraph.</p>
		<p>Second paragraph.</p>
	</body></html>`,
			expected: "First paragraph.",
		},
		{
			fn:       'H',
			name:     "extract h2 even if theres h1",
			input:    "<html><body><h2>Second Title</h2><p>some paragraph</p><h1>Test Title</h1></body></html>",
			expected: "Second Title",
		},
		{
			fn:       'H',
			name:     "extract h2",
			input:    "<html><body><p>some paragraph</p><h2>Test Title</h2></body></html>",
			expected: "Test Title",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var actual string
			switch tc.fn {
			case 'P':
				actual = getFirstParagraphFromHTML(tc.input)
			case 'H':
				actual = getHeadingFromHTML(tc.input)
			}

			if actual != tc.expected {
				t.Errorf("Test %d - %s FAIL: expected HTML: %s, actual: %s", i, tc.name, tc.expected, actual)
			}
		})
	}
}
