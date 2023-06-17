package render

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/justinas/nosurf"
	"github.com/shimon-git/booking-app/internal/config"
	"github.com/shimon-git/booking-app/internal/models"
)

var functions = template.FuncMap{}
var pathToTemplate = "./templates"
var app *config.AppConfig

// NewTemplates sets the config for the remplate package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(templateData *models.TemplateData, r *http.Request) *models.TemplateData {
	templateData.Flash = app.Session.PopString(r.Context(), "flash")
	templateData.Warning = app.Session.PopString(r.Context(), "warning")
	templateData.Error = app.Session.PopString(r.Context(), "error")
	templateData.CSRFToken = nosurf.Token(r)
	return templateData
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, templateName string, templateData *models.TemplateData) error {
	var templateCache map[string]*template.Template
	if app.UseCache {
		// create template cache
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	// get requested template from cache
	tmpl, ok := templateCache[templateName]
	if !ok {
		fmt.Println("Could not get template from template cache.")
		return errors.New("cannot get template from chace")
	}

	buf := new(bytes.Buffer)
	templateData = AddDefaultData(templateData, r)
	err := tmpl.Execute(buf, templateData)
	if err != nil {
		log.Println("Error writing template to browser", err)
	}
	//render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func CreateTemplateCache() (map[string]*template.Template, error) {

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
