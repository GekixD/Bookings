package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the configurations for the application
type AppConfig struct {
	UseCache      bool                          // This allows you to toggle the cached version of the app, used from memory
	TemplateCache map[string]*template.Template // This stores the cached templates
	InfoLog       *log.Logger                   // This is a logger
	Prod          bool                          // This signifies if the web app is in a production environment or not
	Session       *scs.SessionManager           // This manages the current session
}
