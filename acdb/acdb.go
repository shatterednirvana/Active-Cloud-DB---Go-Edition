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
	http.HandleFunc("/put", put)
	http.HandleFunc("/query", query)
	http.HandleFunc("/delete", delete)
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, mainPage)
}

const mainPage = `
<html>
  <head>
    <title>Active Cloud DB - Now with 100% more Go!</title>
    <link rel="stylesheet" href="/static/css/bootstrap-1.1.0.min.css">
    <script src="static/js/jquery.js" type="text/javascript"></script>
    <script src="static/js/jquery.tablesorter.min.js" type="text/javascript"></script>
    <script src="static/js/custom.js" type="text/javascript"></script>
  </head>
  <body>
    <div class="container-fluid">
      <div class="sidebar">
        <h4>Active Cloud DB</h4>
        <h4>Go Edition!</h4>
        <form class="clearfix">
          <legend>Add a new key/value:</legend>
          <label>Key:</label>
          <input class="medium" id="key" name="mInput" size="15" type="text" /> 
          <label>Value:</label>
          <input class="medium" id="val" name="mInput" size="15" type="text" /> 
          <br />
          <br />
          <a href="#" id="put" class="put btn primary">Save Changes</a>
          &nbsp;
          <button type="reset" class="btn">Clear</button> 
        </form>
      </div>
      <div class="content">
        <h4>All Items</h4>
        <form class="clearfix">
          <table id="queryData" class="common-table zebra-striped">
  	    <tr><th>Key</th><th>Value</th><th></th></tr>
          </table>
        </form>
	<div>
	<img src="/static/img/appengine-silver-120x30.gif" alt="Powered by Google App Engine" />
	</div>
      </div>
    </div>
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
		"error":"",
	}

	if err := datastore.Get(c, key, entity); err == nil {
		result[keyName] = entity.Value
	} else {
		result["error"] = fmt.Sprintf("%s", err)
	}

	fmt.Fprintf(w, "%s", mapToJson(result))
}

func put(w http.ResponseWriter, r *http.Request) {
	keyName := r.FormValue("key")
	value := r.FormValue("val")

	c := appengine.NewContext(r)

	key := datastore.NewKey("Entity", keyName, 0, nil)
	entity := new(Entity)
	entity.Value = value

	result := map[string] string {
		"error":"",
	}
	if _, err := datastore.Put(c, key, entity); err != nil {
		result["error"] = fmt.Sprintf("%s", err)
	}

	fmt.Fprintf(w, "%s", mapToJson(result))
}

func query(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Entity")

	result := map[string] string {}
	for t := q.Run(c); ; {
		var entity Entity
		key, err := t.Next(&entity)
		if err == datastore.Done {
			break
		}
		if err != nil {
			result["error"] = fmt.Sprintf("%s", err)
		}
		keyString := fmt.Sprintf("%s", key)
		result[keyString] = entity.Value
	}

	fmt.Fprintf(w, "%s", mapToJson(result))
}

func delete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	keyName := r.FormValue("key")
	key := datastore.NewKey("Entity", keyName, 0, nil)

	result := map[string] string {
		"error":"",
	}
	if err := datastore.Delete(c, key); err != nil {
		result["error"] = fmt.Sprintf("%s", err)
	}

	fmt.Fprintf(w, "%s", mapToJson(result))
}

func mapToJson(mapToConvert map[string] string) []byte {
	jsonResult, err := json.Marshal(mapToConvert)
        if err != nil {
                fmt.Printf("json marshalling saw error: %s\n", err)
        }

	return jsonResult
}
