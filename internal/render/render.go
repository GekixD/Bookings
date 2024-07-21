package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/GekixD/Bookings/internal/config"
	"github.com/GekixD/Bookings/internal/models"
	"github.com/justinas/nosurf"
)

// METHOD 3 - AUTOMATIC CACHE

var app *config.AppConfig
var tmplPath = "./tempates"        // the path to the template files
var functions = template.FuncMap{} // a map of functions that will render with the templates

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds the default data packed in a models.TemplateData (eg CSRF token)
func AddDefaultData(tmplData *models.TemplateData, req *http.Request) *models.TemplateData {
	tmplData.FlashMsg = app.Session.PopString(req.Context(), "flash")  // Puts something in the session until the next time a
	tmplData.Warning = app.Session.PopString(req.Context(), "warning") // until the next time a page is displayed
	tmplData.Error = app.Session.PopString(req.Context(), "error")     // and then it's automatically taken out
	tmplData.CSRFToken = nosurf.Token(req)                             // CSRF token insterted into template data
	return tmplData
}

// RenderTemplate handles the client-side rendering of a page handler
func RenderTemplate(res http.ResponseWriter, req *http.Request, tmpl string, tmplData *models.TemplateData) {

	var templCache map[string]*template.Template

	if app.UseCache {
		// if UseCache is true, get the template cache from the app config
		templCache = app.TemplateCache
	} else {
		// else rebuild it from the start
		templCache, _ = CreateTemplateCache()
	}
	// get requested template from cahce
	t, ok := templCache[tmpl]
	if !ok {
		log.Fatal("could not get template from cache")
	}

	// we can create a buffer for a finer grain error checking (will be touched on later)
	buffer := new(bytes.Buffer)

	tmplData = AddDefaultData(tmplData, req)

	err := t.Execute(buffer, tmplData)
	if err != nil {
		log.Println("Error executing template: ", err)
	}

	// render the template
	_, err = buffer.WriteTo(res)
	if err != nil {
		fmt.Println("Error printing the template to browser: ", err)
	}
}

// CreateTemplateCache creates the template cache for all templates, without the need to call it everytime
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all *.page.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", tmplPath))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		tmplSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// check if there are any .layout.tmpl files
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", tmplPath))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			tmplSet, err = tmplSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", tmplPath))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = tmplSet
	}

	return myCache, nil
}

/*
	METHOD 1 - NAIVE
	This would be the naive approach for rendering templates

	func naiveRenderTemplate(res http.ResponseWriter, tmpl string) {
		parseTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")
		err := parseTemplate.Execute(res, nil)
		if err != nil {
			log.Println("Error parsing template: ", err)
			return
		}
	}
*/

/*
METHOD 2 - MANUAL CACHE
This method uses a manual cache to store all pages visited already
// This variable allows us to store the already parsed template files in cache
var tmplCahce = make(map[string]*template.Template)

func RenderTemplate(res http.ResponseWriter, t string) {
	var templ *template.Template
	var err error
	// We check if the template t is in cahce already
	_, inMap := tmplCahce[t]

	if !inMap {
		log.Println("creating tempalte and storing in in cahce")
		err = createTemplateCahce(t)
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Println("using cached template")
	}

	templ = tmplCahce[t]

	err = templ.Execute(res, nil)
	if err != nil {
		log.Println(err)
	}
}

// We add the template t in cache, as a key
func createTemplateCahce(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	tmplCahce[t] = tmpl

	return nil
}
*/
