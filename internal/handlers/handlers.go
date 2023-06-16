package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/shimon-git/booking-app/internal/config"
	"github.com/shimon-git/booking-app/internal/forms"
	"github.com/shimon-git/booking-app/internal/models"
	"github.com/shimon-git/booking-app/internal/render"
)

// Repo the reposetory used by the handlers
var Repo *Reposetory

// Reposetory is the reposetory type
type Reposetory struct {
	App *config.AppConfig
}

// NewRepo creates a new reposetory
func NewRepo(a *config.AppConfig) *Reposetory {
	return &Reposetory{
		App: a,
	}
}

// NewHandlers sets the reposetory for the handlers
func NewHandlers(r *Reposetory) {
	Repo = r
}

// Home is the home page handler
func (m *Reposetory) Home(w http.ResponseWriter, r *http.Request) {
	remotIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remotIP)
	render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{})
}

// About is the about page handler
func (m *Reposetory) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["testAbout"] = "Hello again! --> about page"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Generals renders the generals room page
func (m *Reposetory) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.html", &models.TemplateData{})
}

// Majors renders the majors room page
func (m *Reposetory) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.html", &models.TemplateData{})
}

// Avilability renders the check-avilability page
func (m *Reposetory) Avilability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "check-avilability.page.html", &models.TemplateData{})
}

// Contact renders the contact page
func (m *Reposetory) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.html", &models.TemplateData{})
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
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the contact page
func (m *Reposetory) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyResevation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyResevation

	render.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostReservation handle  the posting of reservation form
func (m *Reposetory) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	//form.Has("first_name", r)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Reposetory) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "reservation-summary.page.html", &models.TemplateData{})
}
