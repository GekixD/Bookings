package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GekixD/Bookings/internal/config"
	"github.com/GekixD/Bookings/internal/handlers"
	"github.com/GekixD/Bookings/internal/helpers"
	"github.com/GekixD/Bookings/internal/models"
	"github.com/GekixD/Bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	err := run()
	if err != nil {
		log.Fatal("Fatal error starting the web application: ", err)
	}

	srv := &http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}

	fmt.Printf("Starting web application on port %s \n", PORT)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Fatal Server Error: ", err)
	}
}

// run allows all app related logic to be outside the mail function
func run() error {
	// What we want our session to contain:
	gob.Register(models.Reservation{}) //What do I want to store in the session
	app.Prod = false                   // whether the web app is in producton or development

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                  // This is a logger for general info
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // This is a logger for error messages
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour              // set the lifetime for the session to 24 hours
	session.Cookie.Persist = true                  // whether the session will persis if they close the window
	session.Cookie.SameSite = http.SameSiteLaxMode // how strict is the cookie enforcement in the site
	session.Cookie.Secure = app.Prod               // whether the cookies are encrupted (http vs https)

	app.Session = session

	tmplCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Can not create template cache.")
		return err
	}
	app.TemplateCache = tmplCache
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	helpers.NewHelepers(&app)

	return nil
}
