package hello

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"http"
	"json"
)

type Entity struct {
	Value string
}

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/get", get)
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, mainPage)
}

const mainPage = `
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

func get(w http.ResponseWriter, r *http.Request) {
	keyName := r.FormValue("key")

	c := appengine.NewContext(r)

	key := datastore.NewKey("Entity", keyName, 0, nil)
	entity := new(Entity)

	result := map[string] string {
		keyName:"",
	}

	if err := datastore.Get(c, key, entity); err == nil {
		result[keyName] = entity.Value
	} else {
		result[keyName] = fmt.Sprintf("%s", err)
	}

	jsonResult, err := json.Marshal(result)

	if err != nil {
		fmt.Printf("json marshalling saw error: %s\n", err)
	}

	fmt.Fprintf(w, "%s", jsonResult)
}
