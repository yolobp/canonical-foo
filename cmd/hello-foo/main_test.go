package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleFunc(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(handleFunc))
	defer s.Close()

	resp, err := http.DefaultClient.Get(s.URL)
	if err != nil {
		t.Fatalf("http request failed: %v", err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}
	defer resp.Body.Close()

	want := `Your request:
Method: GET
`
	if got := string(b); !strings.Contains(got, want) {
		t.Errorf("response got %q want substring %q", got, want)
	}
}
