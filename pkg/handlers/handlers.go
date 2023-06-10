package handlers

import (
	"net/http"

	"github.com/shimon-git/booking-app/pkg/config"
	"github.com/shimon-git/booking-app/pkg/models"
	"github.com/shimon-git/booking-app/pkg/render"
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
	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

// About is the about page handler
func (m *Reposetory) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["testAbout"] = "Hello again! --> about page"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Generals renders the generals room page
func (m *Reposetory) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "generals.page.html", &models.TemplateData{})
}

// Majors renders the majors room page
func (m *Reposetory) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "majors.page.html", &models.TemplateData{})
}

// Avilability renders the check-avilability page
func (m *Reposetory) Avilability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "check-avilability.page.html", &models.TemplateData{})
}

// Contact renders the contact page
func (m *Reposetory) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact.page.html", &models.TemplateData{})
}

// Contact renders the contact page
func (m *Reposetory) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "make-reservation.page.html", &models.TemplateData{})
}
