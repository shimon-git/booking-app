package handlers

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/shimon-git/booking-app/internal/config"
	"github.com/shimon-git/booking-app/internal/models"
	"github.com/shimon-git/booking-app/internal/render"
)

var pathToTemplate = "./../../templates"
var app config.AppConfig
var session *scs.SessionManager
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	// what am i going to put in the session
	gob.Register(models.Reservation{})

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	templateCache, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Cennot create template cache.")
	}

	app.TemplateCache = templateCache
	app.UseCache = true
	repo := NewRepo(&app)
	NewHandlers(repo)
	render.NewTemplates(&app)

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	//mux.Use(Nosurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Get("/check-avilability", Repo.Avilability)
	mux.Post("/check-avilability", Repo.PostAvilability)
	mux.Post("/check-avilability-json", Repo.JsonAvilability)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

// Nosurf adds CSRF protection to all post requests
func Nosurf(next http.Handler) http.Handler {
	csfrHandler := nosurf.New(next)

	csfrHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csfrHandler
}

// SessionLoad loads and saves the sessions on every request.
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache create a chace as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {

	// creating our map to store the tempates cache
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.html from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplate))
	if err != nil {
		return myCache, err
	}

	// range trough all files ending with .page.html
	for _, page := range pages {
		name := filepath.Base(page)
		tmpl, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		layouts, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplate))
		if err != nil {
			return myCache, err
		}

		if len(layouts) > 0 {
			tmpl, err = tmpl.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplate))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = tmpl
	}
	return myCache, nil
}
