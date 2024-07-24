package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNew(t *testing.T) {
	req := httptest.NewRequest("POST", "/some-url", nil)
	postedData := url.Values{}
	postedData.Add("test", "value")

	form := New(req.PostForm)
	if !form.Valid() {
		t.Error("No new form was created.")
	}
}

func TestForm_Valid(t *testing.T) {
	req := httptest.NewRequest("POST", "/some-url", nil)
	form := New(req.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("Got an invalid form, while it should have been valid.")
	}
}

func TestForm_Required(t *testing.T) {
	req := httptest.NewRequest("POST", "/some-url", nil)
	form := New(req.PostForm)

	form.Required("a", "b")
	if form.Valid() {
		t.Error("Form deemed valid while required fields are missing.")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")

	req, _ = http.NewRequest("POST", "/some-url", nil)
	req.PostForm = postedData

	form = New(req.PostForm)
	form.Required("a", "b")
	if !form.Valid() {
		t.Error("Form was invalid with no required field missing.")
	}
}

func TestForm_Has(t *testing.T) {
	req := httptest.NewRequest("POST", "/some-url", nil)
	postedData := url.Values{}
	postedData.Add("valid", "test-value")
	form := New(req.PostForm)

	if form.Has("invalid", req) {
		t.Error("Form has a field while it's supposed not to have.")
	}
	if !form.Has("valid", req) {
		t.Error("Form doesn't have a field it's supposed to have.")
	}
}
