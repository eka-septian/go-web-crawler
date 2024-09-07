package report

import (
	"reflect"
	"testing"
)

func TestSortPages(t *testing.T) {
	tests := []struct {
		input    map[string]int
		expected []Page
	}{
		{
			input: map[string]int{
				"url1": 5,
				"url2": 1,
				"url3": 3,
				"url4": 10,
				"url5": 7,
			},
			expected: []Page{
				{URL: "url4", Count: 10},
				{URL: "url5", Count: 7},
				{URL: "url1", Count: 5},
				{URL: "url3", Count: 3},
				{URL: "url2", Count: 1},
			},
		},
		{
			input: map[string]int{
				"d": 1,
				"a": 1,
				"e": 1,
				"b": 1,
				"c": 1,
			},
			expected: []Page{
				{URL: "a", Count: 1},
				{URL: "b", Count: 1},
				{URL: "c", Count: 1},
				{URL: "d", Count: 1},
				{URL: "e", Count: 1},
			},
		},
		{
			input: map[string]int{
				"d": 2,
				"a": 1,
				"e": 3,
				"b": 1,
				"c": 2,
			},
			expected: []Page{
				{URL: "e", Count: 3},
				{URL: "c", Count: 2},
				{URL: "d", Count: 2},
				{URL: "a", Count: 1},
				{URL: "b", Count: 1},
			},
		},
		{
			input:    map[string]int{},
			expected: []Page{},
		},
		{
			input:    nil,
			expected: []Page{},
		},
		{
			input: map[string]int{
				"url1": 1,
			},
			expected: []Page{
				{URL: "url1", Count: 1},
			},
		},
	}

	for i, tc := range tests {
		pages := sortPages(tc.input)
		if !reflect.DeepEqual(pages, tc.expected) {
			t.Errorf("Test %d - expected: %v, got: %v", i, tc.expected, pages)
		}
	}
}
