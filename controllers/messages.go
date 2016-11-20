package controllers

import (
	"github.com/perthgophers/puddle/models"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

// loadFile locates the log and returns it
func loadFile() (*models.Page, error) {
	filename := path.Join("logs", "log.txt")
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &models.Page{Body: body}, nil
}

// renderTemplate locates and renders the html file
func renderTemplate(w http.ResponseWriter, p *models.Page) {
	lp := path.Join("views", "log.html")
	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		log.Println(err)
	}

	tmpl.ExecuteTemplate(w, "log.html", p)
}

// ServeTastic executes the loadFile and renderTemplate functions
func ServeTastic(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "405: Method Not Allowed", http.StatusMethodNotAllowed)

	} else {

		p, err := loadFile()
		if err != nil {
			log.Println(err)
		}

		renderTemplate(w, p)
	}

}
