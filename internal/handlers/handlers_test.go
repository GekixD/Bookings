package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// postData holds the data stored in post requests
type postData struct {
	key   string
	value string
}

var tests = []struct {
	name           string     // the test name
	url            string     // the url matched by the roots
	method         string     // specify Get or Post
	params         []postData // hte parameters for the test
	expectedStatus int        // the expected status code for the test (2xx, 3xx, ...)
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},                            // Get request to home page
	{"about", "/about", "GET", []postData{}, http.StatusOK},                      // Get request to about page
	{"generals", "/generals", "GET", []postData{}, http.StatusOK},                // Get request to general's page
	{"majors", "/majors", "GET", []postData{}, http.StatusOK},                    // Get request to major's page
	{"search-avail", "/search-availability", "GET", []postData{}, http.StatusOK}, // Get request to search availability page
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},                  // Get request to contact page
	{"make-res", "/make-reservation", "GET", []postData{}, http.StatusOK},        // Get request to search availability page
	{"post-search-avail", "/search-availability", "POST", []postData{
		{key: "start", value: "2024-01-01"},
		{key: "end", value: "2024-01-02"},
	}, http.StatusOK},
	{"post-search-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2024-01-01"},
		{key: "end", value: "2024-01-01"},
	}, http.StatusOK},
	{"post-make-res", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "Jon"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "me@here.com"},
		{key: "phone", value: "0123456789"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	testServer := httptest.NewTLSServer(routes) // Go allows us to create a testing server
	defer testServer.Close()                    // remember to close the server after done testing

	for _, test := range tests {
		if test.method == "GET" {
			response, err := testServer.Client().Get(testServer.URL + test.url)
			if err != nil {
				t.Logf("Error encountered when performing test: %s. Error: %s", test.name, err)
				t.Fatal(err)
			}

			if response.StatusCode != test.expectedStatus {
				t.Errorf("For %s, we expected status code %d but instead got %d", test.name, test.expectedStatus, response.StatusCode)
			}

		} else {
			values := url.Values{}
			for _, val := range test.params {
				values.Add(val.key, val.value)
			}
			response, err := testServer.Client().PostForm(testServer.URL+test.url, values)
			if err != nil {
				t.Logf("Error encountered when performing test: %s. Error: %s", test.name, err)
				t.Fatal(err)
			}

			if response.StatusCode != test.expectedStatus {
				t.Errorf("For %s, we expected status code %d but instead got %d", test.name, test.expectedStatus, response.StatusCode)
			}
		}
	}
}
