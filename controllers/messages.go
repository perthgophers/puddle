package controllers

import (
	"../models"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

// locates the log and passes it
func loadFile() (*models.Page, error) {
	filename := path.Join("logs", "log.txt")
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &models.Page{Body: body}, nil
}

// locats and renders the html file
func renderTemplate(w http.ResponseWriter, p *models.Page) {
	lp := path.Join("views", "log.html")
	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		log.Println(err)
	}

	tmpl.ExecuteTemplate(w, "log.html", p)
}

// puts everything together
func ServeTastic(w http.ResponseWriter, r *http.Request) {
	p, err := loadFile()
	if err != nil {
		log.Println(err)
	}

	renderTemplate(w, p)
}
