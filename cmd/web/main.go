package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/GekixD/Bookings/internal/config"
	"github.com/GekixD/Bookings/internal/handlers"
	"github.com/GekixD/Bookings/internal/models"
	"github.com/GekixD/Bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	// What we want our session to contain:
	gob.Register(models.Reservation{}) //What do I want to store in the session
	app.Prod = false                   // whether the web app is in producton or development

	session = scs.New()
	session.Lifetime = 24 * time.Hour              // set the lifetime for the session to 24 hours
	session.Cookie.Persist = true                  // whether the session will persis if they close the window
	session.Cookie.SameSite = http.SameSiteLaxMode // how strict is the cookie enforcement in the site
	session.Cookie.Secure = app.Prod               // whether the cookies are encrupted (http vs https)

	app.Session = session

	tmplCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("can not create template cache, due to error: ", err)
	}
	app.TemplateCache = tmplCache
	app.UseCache = false

	repo := handlers.NewRepo(&app)

	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	srv := &http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}

	fmt.Println("Starting web application on port: ", PORT)

	err = srv.ListenAndServe()
	log.Fatal("Fatal Server Error: ", err)
}
