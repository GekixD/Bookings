package handlers

import (
	"net/http"

	"github.com/GekixD/Bookings/pkg/config"
	"github.com/GekixD/Bookings/pkg/models"
	"github.com/GekixD/Bookings/pkg/render"
)

// Repository type initialization
type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

// This creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// This sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (r *Repository) Home(res http.ResponseWriter, req *http.Request) {
	remoteIP := req.RemoteAddr
	r.App.Session.Put(req.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(res, "home.page.tmpl", &models.TemplateData{})
}

func (r *Repository) About(res http.ResponseWriter, req *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "OK"

	remoteIP := r.App.Session.GetString(req.Context(), "remote_ip")

	stringMap["remote_ip"] = remoteIP

	// Template Data to be passed on to RenderTemplate()
	tmplData := &models.TemplateData{
		StringMap: stringMap,
	}

	render.RenderTemplate(res, "about.page.tmpl", tmplData)
}

// Reservation renders the make a reservation page and displays form
func (r *Repository) Reservation(res http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(res, "make-reservation.page.tmpl", &models.TemplateData{})
}

// Generals renders the General's room page
func (r *Repository) Generals(res http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(res, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the Majors's room page
func (r *Repository) Majors(res http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(res, "generals.page.tmpl", &models.TemplateData{})
}

// Availability renders the Majors's room page
func (r *Repository) Availability(res http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(res, "search-availability.page.tmpl", &models.TemplateData{})
}
