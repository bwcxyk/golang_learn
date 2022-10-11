package cros

import (
	"strings"
	"testing"
)

func Test_getReferer(t *testing.T) {
	urls := []string{
		"https://itsos.ltd/aaa/bbb",
		"https://itsos.ltd/",
		"https://itsos.ltd",
	}

	expected := "https://itsos.ltd"

	for _, url := range urls {
		start := strings.Index(url, ":") + 3
		end := strings.Index(url[start:], "/")
		if end > -1 {
			url = url[:start+end]
		}
		if url != expected {
			t.Error("不符合预期.")
		}
	}
}
