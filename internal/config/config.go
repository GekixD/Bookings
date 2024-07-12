package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the configurations for the application
type AppConfig struct {
	UseCache      bool // This allows you to toggle the cached version of the app, used from memory
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	Prod          bool
	Session       *scs.SessionManager
}
