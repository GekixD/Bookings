package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/GekixD/Bookings/internal/config"
	"github.com/GekixD/Bookings/internal/forms"
	"github.com/GekixD/Bookings/internal/models"
	"github.com/GekixD/Bookings/internal/render"
)

// Repository type initialization
type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
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
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderTemplate(res, req, "make-reservations.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostReservation handles the posting of a reservation form
func (r *Repository) PostReservation(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Println("Error parsing the Form data: ", err)
		return
	}

	reservation := models.Reservation{
		FirstName: req.Form.Get("first_name"),
		LastName:  req.Form.Get("last_name"),
		Email:     req.Form.Get("email"),
		Phone:     req.Form.Get("phone"),
	}

	form := forms.New(req.PostForm)
	// Form validation logic, the order the messages will be displayed is the order the errors appears (due to the forms.Get function)
	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3, req) // The first name needs to be at least 3 characters long
	form.MinLength("last_name", 3, req)  // The last name needs to be at least 3 characters long
	form.MinLength("phone", 10, req)     // The phone number needs to be at least 3 characters long
	form.IsEmail("email")
	// If you find any validation errors
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		// render the make reservation page with the form and data provided but don't continue (return)
		render.RenderTemplate(res, req, "make-reservations.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	r.App.Session.Put(req.Context(), "reservation", reservation)         // Redirect to reservation summary with all the required data
	http.Redirect(res, req, "/reservation-summary", http.StatusSeeOther) // redirect after submitting the form to the reservation summary page
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

// ReservationSummary renders the page that displays the reservation's details to the user
func (r *Repository) ReservationSummary(res http.ResponseWriter, req *http.Request) {
	reservation, ok := r.App.Session.Get(req.Context(), "reservation").(models.Reservation)
	if !ok {
		// In case the user tried to access this page without submitting a form
		log.Println("Can not get item from session")
		r.App.Session.Put(req.Context(), "error", "Can't get reservation from current session")
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}

	// remove reservation data from the session after the post is complete
	r.App.Session.Remove(req.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplate(res, req, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
