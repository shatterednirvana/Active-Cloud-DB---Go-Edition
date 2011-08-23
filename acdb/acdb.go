package hello

import (
	"fmt"
	"http"
	"template"
)

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/get", sign)
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, guestbookForm)
}

const guestbookForm = `
<html>
  <head>
    <title>Active Cloud DB - Now with 100% more Go!</title>
    <link rel="stylesheet" href="http://twitter.github.com/bootstrap/assets/css/bootstrap-1.1.0.min.css">
    <link rel="javascript" href="http://ajax.googleapis.com/ajax/libs/jquery/1.6.2/jquery.min.js">
    <link rel="javascript" href="/static/js/custom.js">
  </head>
  <body>
  </body>
</html>
`

func sign(w http.ResponseWriter, r *http.Request) {
	err := signTemplate.Execute(w, r.FormValue("content"))
	if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
	}
}

var signTemplate = template.MustParse(signTemplateHTML, nil)

const signTemplateHTML = `
<html>
  <body>
    <p>You wrote:</p>
    <pre>{@|html}</pre>
  </body>
</html>
`

