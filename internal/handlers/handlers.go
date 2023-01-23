package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/kapi1023/bookingWebsite/internal/config"
	"github.com/kapi1023/bookingWebsite/internal/driver"
	"github.com/kapi1023/bookingWebsite/internal/forms"
	"github.com/kapi1023/bookingWebsite/internal/helpers"
	"github.com/kapi1023/bookingWebsite/internal/models"
	render "github.com/kapi1023/bookingWebsite/internal/renders"
	"github.com/kapi1023/bookingWebsite/internal/repository"
	"github.com/kapi1023/bookingWebsite/internal/repository/dbrepo"
)

//Repo is repository used by the handlers

var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DataBaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostGresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "home.page.html", &models.TemplateData{}, r)
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	// send data to the template
	render.Template(w, "about.page.html", &models.TemplateData{}, r)
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, "make-reservation.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	}, r)
}

// Reservation renders the reservation
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	stDate := r.Form.Get("start_date")
	edDate := r.Form.Get("end_date")
	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, stDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, edDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomId:    roomId,
	}

	log.Println(reservation.FirstName, reservation.LastName, reservation.Email, reservation.Phone, reservation.StartDate, reservation.EndDate, reservation.RoomId)
	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		}, r)
		return
	}

	newReservationId, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restricion := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomId:        roomId,
		ReservationId: newReservationId,
		RestricionId:  1,
	}
	err = m.DB.InsertRoomRestriction(restricion)
	if err != nil {
		helpers.ServerError(w, err)
	}
	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// PostReservation renders the reservation form
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "generals.page.html", &models.TemplateData{}, r)
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "majors.page.html", &models.TemplateData{}, r)
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "search-availability.page.html", &models.TemplateData{}, r)
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("Start: %s End: %s", start, end)))

}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// JsonAvailability renders the search availability page
func (m *Repository) JsonAvailability(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		Ok:      true,
		Message: "Available",
	}
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "contact.page.html", &models.TemplateData{}, r)
}

// ReservationSummary renders reservation-summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cant get error from session")
		log.Println("cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Cant get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, "reservation-summary.page.html", &models.TemplateData{
		Data: data,
	}, r)
}
