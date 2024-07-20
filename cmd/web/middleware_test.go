package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler

	h := NoSurf(&myH)
	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("type mismatch, expected 'http.Handler' but instead got: %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler

	h := NoSurf(&myH)
	switch v := h.(type) {
	case http.Handler:
		// do nothing, test passed
	default:
		t.Errorf("type mismatch, expected 'http.Handler' but instead got: %T", v)
	}
}
