package render

import (
	"net/http"
	"testing"

	"github.com/shimon-git/booking-app/internal/models"
)

func TestNewRenderer(t *testing.T) {
	NewRenderer(app)
}

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)
	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
}

func TestTemplate(t *testing.T) {
	pathToTemplate = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	var ww myWriter
	err = Template(&ww, r, "about.page.html", &models.TemplateData{})
	if err != nil {
		t.Errorf("error ocuured while writing template to browser - %v", err)
	}

}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplate = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}
	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	return r.WithContext(ctx), nil
}
