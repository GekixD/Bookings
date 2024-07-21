package render

import (
	"net/http"
	"testing"

	"github.com/GekixD/Bookings/internal/models"
)

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	context := r.Context()
	context, _ = session.Load(context, r.Header.Get("X-Session"))

	r = r.WithContext(context)

	return r, nil
}

func TestAddDefaultData(t *testing.T) {
	var tmplData models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error("Error while creating an http request: ", err)
	}

	session.Put(r.Context(), "flash", "Test - 123")
	result := AddDefaultData(&tmplData, r)
	if result == nil {
		t.Error("Error while adding the default tempalte data")
	} else if result.FlashMsg != "Test - 123" {
		t.Error("Flash message value (of Test - 123) not found in session")
	}
}

// TestRenderTemplate
func TestRenderTemplate(t *testing.T) {
	tmplPath = "./../../templates"
	tmplCache, err := CreateTemplateCache()
	if err != nil {
		t.Error("Error creating the template cache: ", err)
	}

	app.TemplateCache = tmplCache

	req, err := getSession()
	if err != nil {
		t.Error("Error creating the session context for the request: ", err)
	}

	var resWriter myWriter

	// Error shouldn't exist
	err = RenderTemplate(&resWriter, req, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("Error writing template to browser: ", err)
	}

	// Error should exit
	err = RenderTemplate(&resWriter, req, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("Rendered a non existent template")
	}

	// Testing template cache
	app.UseCache = true // set UseCache to true
	err = RenderTemplate(&resWriter, req, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("Error while using app.TempalteCache: ", err)
	}
	app.UseCache = false // set UseCahce to false again
}

// New templates propably doesn't need testing but we did for test coverage
func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCahce(t *testing.T) {
	tmplPath = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error("Errro creating the template cahce: ", err)
	}
}
