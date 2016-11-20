package models

import (
	"html/template"
)

// Page struct to pass data to the html template
type Page struct {
	Body template.HTML
}
