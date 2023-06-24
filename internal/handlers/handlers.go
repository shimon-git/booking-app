package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
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
	dateFormat := "D-M-Y"

	startDate, err := helpers.StringToDateTime(dateFormat, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := helpers.StringToDateTime(dateFormat, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	rooms, err := m.DB.CheckAvialibilityByDatesForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No Avilability")
		http.Redirect(w, r, "/check-avilability", http.StatusSeeOther)
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)
	render.Template(w, r, "choose-room.page.html", &models.TemplateData{
		Data: data,
	})
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
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cennot get reservation from the session"))
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	dateLayout := "D-M-Y"
	startDateS, err := helpers.DateTimeToString(res.StartDate, dateLayout)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDateS, err := helpers.DateTimeToString(res.EndDate, dateLayout)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	stringMap := make(map[string]string)
	data := make(map[string]interface{})
	data["reservation"] = res
	stringMap["start_date"] = startDateS
	stringMap["end_date"] = endDateS

	render.Template(w, r, "make-reservation.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handle  the posting of reservation form
func (m *Reposetory) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot retrive reservation from the session"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

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

	m.App.Session.Put(r.Context(), "reservation", reservation)

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
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

	dateLayout := "D-M-Y"
	startDate, err := helpers.DateTimeToString(reservation.StartDate, dateLayout)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := helpers.DateTimeToString(reservation.EndDate, dateLayout)
	if err != nil {
		helpers.ServerError(w, err)
	}

	stringMap := make(map[string]string)
	stringMap["start_date"] = startDate
	stringMap["end_date"] = endDate

	render.Template(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Reposetory) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cennot get the resevation from the session"))
	}

	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
