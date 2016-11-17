package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var htmlString = `
<!doctype html>

<html lang="en">
<head>
  <meta charset="utf-8">
  <title>{{ .Title }}</title>
</head>

<body>
  {{ .Body }}
</body>
</html>
`

func main() {
	// parse template
	tpl, err := template.New("log_template").Parse(htmlString)

	// read log file into variable
	logBytes, err := ioutil.ReadFile("./logs/current.log")
	if err != nil {
		log.Fatalln("Did not parse log file!")
	}

	// Create inline struct with required data
	data := struct {
		Body  string
		Title string
	}{
		Body:  strings.Replace(string(logBytes), "\n", "<br/>", -1), // convert bytes to string and replace newlines with html breaks
		Title: "Current Logs",
	}

	// execute template, send data struct through as pointer
	err = tpl.Execute(os.Stdout, &data)
	if err != nil {
		log.Fatalln("Didn't work!")
	}
}
