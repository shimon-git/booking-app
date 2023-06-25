package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/shimon-git/booking-app/internal/config"
	"github.com/shimon-git/booking-app/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	// what am i going to put in the session
	gob.Register(models.Reservation{})

	testApp.InProduction = false
	testApp.InfoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	testApp.ErrorLog = log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session

	app := &testApp
	NewRenderer(app)

	os.Exit(m.Run())
}

type myWriter struct{}

func (w *myWriter) Header() http.Header {
	var header http.Header
	return header
}

func (w *myWriter) WriteHeader(statusCode int) {

}

func (w *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
