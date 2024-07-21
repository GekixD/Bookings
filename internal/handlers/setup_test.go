package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/GekixD/Bookings/internal/config"
	"github.com/GekixD/Bookings/internal/models"
	"github.com/GekixD/Bookings/internal/render"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var tmplPath = "./../../templates"

func getRoutes() http.Handler {
	// What we want our session to contain:
	gob.Register(models.Reservation{}) //What do I want to store in the session
	app.Prod = false                   // whether the web app is in producton or development

	session = scs.New()
	session.Lifetime = 24 * time.Hour              // set the lifetime for the session to 24 hours
	session.Cookie.Persist = true                  // whether the session will persis if they close the window
	session.Cookie.SameSite = http.SameSiteLaxMode // how strict is the cookie enforcement in the site
	session.Cookie.Secure = app.Prod               // whether the cookies are encrupted (http vs https)

	app.Session = session

	tmplCache, err := CreateTemplateCacheT()
	if err != nil {
		log.Fatal("Can not create template cache.")
	}
	app.TemplateCache = tmplCache
	app.UseCache = true

	repo := NewRepo(&app)
	NewHandlers(repo)
	render.NewTemplates(&app)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurfT) , Currently we don't need to test for the CSRF Token
	mux.Use(SessionLoadT)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/generals", Repo.Generals)
	mux.Get("/majors", Repo.Majors)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// NoSurfT is used to setup tests using NoSurf middleware
func NoSurfT(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoadT is used to setup test using SessionLoad middleware
func SessionLoadT(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTemplateCacheT is used to setup test using SessionLoad middleware
func CreateTemplateCacheT() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all *.page.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", tmplPath))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		tmplSet, err := template.New(name).ParseFiles(page)
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
