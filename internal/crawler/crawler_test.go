package crawler

import (
	"reflect"
	"strings"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input        string
		expected     string
		containError string
	}{
		{
			input:    "https://www.example.com/pat1/",
			expected: "www.example.com/pat1",
		},
		{
			input:    "https://www.example.com/pat1",
			expected: "www.example.com/pat1",
		},
		{
			input:    "https://WWW.example.COM/pat1/",
			expected: "www.example.com/pat1",
		},
		{
			input:    "http://WWW.example.COM/pat1/",
			expected: "www.example.com/pat1",
		},
		{
			input:        ":\\invalid",
			expected:     "",
			containError: "Couldn't normalize URL",
		},
	}

	for i, tt := range tests {
		url, err := normalizeURL(tt.input)
		if err != nil && !strings.Contains(err.Error(), tt.containError) {
			t.Fatalf("Tests %d - unexpected error: %s", i, err)
		} else if err != nil && tt.containError == "" {
			t.Fatalf("Tests %d - unexpected error: %s", i, err)
		} else if err == nil && tt.containError != "" {
			t.Fatalf("Tests %d - expected error containing %s got none.", i, tt.containError)
		}
		if url != tt.expected {
			t.Fatalf("Tests %d - expected url: %s, got: %s", i, tt.expected, url)
		}
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	cases := []struct {
		inputURL      string
		inputBody     string
		expected      []string
		errorContains string
	}{
		{
			inputURL: "https://blog.boot.dev",
			inputBody: `
            <html>
                <body>
                    <a href="https://blog.boot.dev">
                        <span>Boot.dev</span>
                    </a>
                </body>
            </html>
            `,
			expected: []string{"https://blog.boot.dev"},
		},
		{
			inputURL: "https://blog.boot.dev",
			inputBody: `
            <html>
                <body>
                    <a href="/path/one">
                        <span>Boot.dev</span>
                    </a>
                </body>
            </html>
            `,
			expected: []string{"https://blog.boot.dev/path/one"},
		},
		{
			inputURL: "https://blog.boot.dev",
			inputBody: `
            <html>
                <body>
                    <a href="/path/one">
                        <span>Boot.dev</span>
                    </a>
                    <a href="https://other.com/path/one">
                        <span>Boot.dev</span>
                    </a>
                </body>
            </html>
            `,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			inputURL: "https://blog.boot.dev",
			inputBody: `
            <html>
                <body>
                    <a>
                        <span>Boot.dev></span>
                    </a>
                </body>
            </html>
            `,
			expected: nil,
		},
		{
			inputURL: "https://blog.boot.dev",
			inputBody: `
            <html body>
                <a href="path/one">
                    <span>Boot.dev></span>
                </a>
            </html body>
            `,
			expected: []string{"https://blog.boot.dev/path/one"},
		},
		{
			inputURL: "https://blog.boot.dev",
			inputBody: `
            <html>
                <body>
                    <a href=":\\invalidURL">
                        <span>Boot.dev</span>
                    </a>
                </body>
            </html>
            `,
			expected: nil,
		},
		{
			inputURL: `:\\invalidBaseURL`,
			inputBody: `
            <html>
                <body>
                    <a href="/path">
                        <span>Boot.dev</span>
                    </a>
                </body>
            </html>
            `,
			expected:      nil,
			errorContains: "couldn't parse base URL",
		},
	}

	for i, tc := range cases {
		urls, err := getURLsFromHTMLBody(tc.inputBody, tc.inputURL)
		if err != nil && !strings.Contains(err.Error(), tc.errorContains) {
			t.Errorf("Tests %d - unexpected error: %v", i, err)
			return
		} else if err != nil && tc.errorContains == "" {
			t.Errorf("Tests %d - unexpected error: %v", i, err)
			return
		} else if err == nil && tc.errorContains != "" {
			t.Errorf("Tests %d - xpected error containing '%s', got none.", i, tc.errorContains)
			return
		}

		if !reflect.DeepEqual(urls, tc.expected) {
			t.Errorf("Test %d - expected URLs %v, got URLs %v", i, tc.expected, urls)
			return
		}
	}
}
