package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/shimon-git/booking-app/pkg/config"
	"github.com/shimon-git/booking-app/pkg/models"
)

func AddDefaultData(templateData *models.TemplateData) *models.TemplateData {
	return templateData
}

var app *config.AppConfig

// NewTemplates sets the config for the remplate package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func RenderTemplate(w http.ResponseWriter, templateName string, templateData *models.TemplateData) {
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
		log.Fatal("Could not get template from template cache.")
	}

	buf := new(bytes.Buffer)

	err := tmpl.Execute(buf, templateData)
	if err != nil {
		log.Println("Error writing template to browser", err)
	}
	//render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {

	// creating our map to store the tempates cache
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.html from ./templates
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	// range trough all files ending with .page.html
	for _, page := range pages {
		name := filepath.Base(page)
		tmpl, err := template.New(name).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		layouts, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(layouts) > 0 {
			tmpl, err = tmpl.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = tmpl
	}
	return myCache, nil
}
