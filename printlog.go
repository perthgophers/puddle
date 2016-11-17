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
	tpl, err := template.New("log_template").Parse(htmlString)
	logBytes, err := ioutil.ReadFile("./logs/current.log")
	if err != nil {
		log.Fatalln("Did not parse log file!")
	}

	data := struct {
		Body  string
		Title string
	}{
		Body:  strings.Replace(string(logBytes), "\n", "<br/>", -1),
		Title: "Current Logs",
	}

	err = tpl.Execute(os.Stdout, &data)
	if err != nil {
		log.Fatalln("Didn't work!")
	}
}
