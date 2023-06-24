package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/shimon-git/booking-app/internal/config"
	"github.com/shimon-git/booking-app/internal/driver"
	"github.com/shimon-git/booking-app/internal/forms"
	"github.com/shimon-git/booking-app/internal/helpers"
	"github.com/shimon-git/booking-app/internal/models"
	"github.com/shimon-git/booking-app/internal/render"
	"github.com/shimon-git/booking-app/internal/reposetory"
	"github.com/shimon-git/booking-app/internal/reposetory/dbrepo"
)

// Repo the reposetory used by the handlers
var Repo *Reposetory

// Reposetory is the reposetory type
type Reposetory struct {
	App *config.AppConfig
	DB  reposetory.DatabaseRepo
}

// NewRepo creates a new reposetory
func NewRepo(a *config.AppConfig, db *driver.DB) *Reposetory {
	return &Reposetory{
		App: a,
		DB:  dbrepo.NewPostgresRepo(a, db.SQL),
	}
}

// NewHandlers sets the reposetory for the handlers
func NewHandlers(r *Reposetory) {
	Repo = r
}

// Home is the home page handler
func (m *Reposetory) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.html", &models.TemplateData{})
}

// About is the about page handler
func (m *Reposetory) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.html", &models.TemplateData{})
}

// Generals renders the generals room page
func (m *Reposetory) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.html", &models.TemplateData{})
}

// Majors renders the majors room page
func (m *Reposetory) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.html", &models.TemplateData{})
}

// Avilability renders the check-avilability page
func (m *Reposetory) Avilability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "check-avilability.page.html", &models.TemplateData{})
}

// Contact renders the contact page
func (m *Reposetory) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.html", &models.TemplateData{})
}

// PostAvilability renders the check-avilability page
func (m *Reposetory) PostAvilability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")
	w.Write([]byte(fmt.Sprintf("start date is: %s and end date is: %s", start, end)))
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// JsonAvilability handles requests for avilability and send json response
func (m *Reposetory) JsonAvilability(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		Ok:      true,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the contact page
func (m *Reposetory) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyResevation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyResevation

	render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostReservation handle  the posting of reservation form
func (m *Reposetory) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	startDateS := r.Form.Get("start_date")
	endDateS := r.Form.Get("end_date")
	dateFormat := "Y-M-D"

	startDateT, err := helpers.DateConvertor(dateFormat, startDateS)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDateT, err := helpers.DateConvertor(dateFormat, endDateS)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDateT,
		EndDate:   endDateT,
		RoomID:    roomID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3)
	form.IsEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	reservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	restriction := models.RoomRestriction{
		StartDate:     startDateT,
		EndDate:       endDateT,
		RoomID:        roomID,
		ReservationID: reservationID,
		RestrictionID: 1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err = m.DB.InsertRoomRestriction(restriction); err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Reposetory) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cennot get item from session")
		m.App.Session.Put(r.Context(), "error", "Cen't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data: data,
	})
}
