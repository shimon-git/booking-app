package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/shimon-git/booking-app/internal/config"
	"github.com/shimon-git/booking-app/internal/driver"
	"github.com/shimon-git/booking-app/internal/handlers"
	"github.com/shimon-git/booking-app/internal/helpers"
	"github.com/shimon-git/booking-app/internal/models"
	"github.com/shimon-git/booking-app/internal/render"
)

const portNumber = ":4444"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Printf("Starting Application on port number: %s...\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}

func run() (*driver.DB, error) {
	// registering our models into the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restrition{})

	app.InProduction = false

	app.InfoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=ShimonTest123!")
	if err != nil {
		log.Fatalf("Cennot connect to th DB dying...\n%v", err)
		return nil, err
	}
	log.Println("Connected to the DB...")

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cennot create template cache.")
		return nil, err
	}

	app.TemplateCache = templateCache
	app.UseCache = false
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	return db, nil
}
