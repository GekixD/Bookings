package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

// myHandler implements the same methods as a "real" http.Handler, for testing purposes
type myHandler struct{}

func (mh *myHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {}
