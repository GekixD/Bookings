package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GekixD/Bookings/internal/config"
	"github.com/GekixD/Bookings/internal/forms"
	"github.com/GekixD/Bookings/internal/helpers"
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
	err := render.RenderTemplate(res, req, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		helpers.ServerError(res, err)
	}
}

func (r *Repository) About(res http.ResponseWriter, req *http.Request) {
	err := render.RenderTemplate(res, req, "about.page.tmpl", &models.TemplateData{})
	if err != nil {
		helpers.ServerError(res, err)
	}
}

// Reservation renders the make a reservation page and displays form
func (r *Repository) Reservation(res http.ResponseWriter, req *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	err := render.RenderTemplate(res, req, "make-reservations.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
	if err != nil {
		helpers.ServerError(res, err)
	}
}

// PostReservation handles the posting of a reservation form
func (r *Repository) PostReservation(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		helpers.ServerError(res, err)
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
	form.MinLength("first_name", 3) // The first name needs to be at least 3 characters long
	form.MinLength("last_name", 3)  // The last name needs to be at least 3 characters long
	form.MinLength("phone", 10)     // The phone number needs to be at least 3 characters long
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
	err := render.RenderTemplate(res, req, "generals.page.tmpl", &models.TemplateData{})
	if err != nil {
		helpers.ServerError(res, err)
	}
}

// Majors renders the Majors's room page
func (r *Repository) Majors(res http.ResponseWriter, req *http.Request) {
	err := render.RenderTemplate(res, req, "majors.page.tmpl", &models.TemplateData{})
	if err != nil {
		helpers.ServerError(res, err)
	}
}

// Availability renders the Majors's room page
func (r *Repository) Availability(res http.ResponseWriter, req *http.Request) {
	err := render.RenderTemplate(res, req, "search-availability.page.tmpl", &models.TemplateData{})
	if err != nil {
		helpers.ServerError(res, err)
	}
}

// PostAvailability renders the Majors's room page
func (r *Repository) PostAvailability(res http.ResponseWriter, req *http.Request) {
	start := req.Form.Get("start")
	end := req.Form.Get("end")

	res.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", start, end)))
}

// Contact renders the page that contains contact info
func (r *Repository) Contact(res http.ResponseWriter, req *http.Request) {
	err := render.RenderTemplate(res, req, "contact.page.tmpl", &models.TemplateData{})
	if err != nil {
		helpers.ServerError(res, err)
	}
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
		helpers.ServerError(res, err)
		return
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
		r.App.ErrorLog.Println("Can not get item from session")
		r.App.Session.Put(req.Context(), "error", "Can't get reservation from current session")
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}

	// remove reservation data from the session after the post is complete
	r.App.Session.Remove(req.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	err := render.RenderTemplate(res, req, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
	if err != nil {
		helpers.ServerError(res, err)
	}
}
