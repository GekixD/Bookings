package handlers

import (
	"encoding/json"
	"fmt"
	"log"
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

	render.RenderTemplate(res, req, "home.page.tmpl", &models.TemplateData{})
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

	render.RenderTemplate(res, req, "about.page.tmpl", tmplData)
}

// Reservation renders the make a reservation page and displays form
func (r *Repository) Reservation(res http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(res, req, "make-reservations.page.tmpl", &models.TemplateData{})
}

// Generals renders the General's room page
func (r *Repository) Generals(res http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(res, req, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the Majors's room page
func (r *Repository) Majors(res http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(res, req, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the Majors's room page
func (r *Repository) Availability(res http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(res, req, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability renders the Majors's room page
func (r *Repository) PostAvailability(res http.ResponseWriter, req *http.Request) {
	start := req.Form.Get("start")
	end := req.Form.Get("end")

	res.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", start, end)))
}

// Contact renders the page that contains contact info
func (r *Repository) Contact(res http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(res, req, "contact.page.tmpl", &models.TemplateData{})
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles a request for availability and send a JSON response
func (r *Repository) AvailabilityJSON(res http.ResponseWriter, req *http.Request) {
	response := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		log.Println("JSON parsing error: ", err)
	}

	// standard Content-Type Header is "text/html" so we need to define it to be JSON
	res.Header().Set("Content-Type", "application/json")
	res.Write(out)
}
