package hello

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"

	"fmt"
	"http"
	"json"
	"strconv"
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
        <form class="form-stacked">
          <legend>Add a new key/value:</legend>
          <br />
          <label>Key:</label>
          <input class="medium" id="key" name="mInput" size="15" type="text" /> 
          <br />
          <br />
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
        <div id="success-flash-div" class="alert-message success"><p id="success-flash" /></div>
        <div id="error-flash-div" class="alert-message error"><p id="error-flash" /></div>
        <h4>All Items</h4>
        <form class="clearfix">
          <table id="queryData" class="common-table zebra-striped">
  	    <tr><th>Key</th><th>Value</th><th></th></tr>
          </table>
        </form>
	<div>
	<a href="http://code.google.com/appengine/">
          <img src="/static/img/appengine-silver-120x30.gif" alt="Powered by Google App Engine" />
        </a>
	<a href="http://golang.org">
          <img src="/static/img/Golang.png" alt="Powered by Go" />
        </a>
	</div>
      </div>
    </div>
  </body>
</html>
`

func get(w http.ResponseWriter, r *http.Request) {
	keyName := r.FormValue("key")

	c := appengine.NewContext(r)

	result := map[string] string {
		keyName:"",
		"error":"",
	}

	if item, err := memcache.Get(c, keyName); err == nil {
		result[keyName] = fmt.Sprintf("%q", item.Value)
		fmt.Fprintf(w, "%s", mapToJson(result))
		return
	}

	key := datastore.NewKey("Entity", keyName, 0, nil)
	entity := new(Entity)

	if err := datastore.Get(c, key, entity); err == nil {
		result[keyName] = entity.Value

		// Set the value to speed up future reads - errors here aren't
		// that bad, so don't worry about them
		item := &memcache.Item{
			Key: keyName,
			Value: []byte(entity.Value),
		}
		memcache.Set(c, item)
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

	// Set the value to speed up future reads - errors here aren't
	// that bad, so don't worry about them
	item := &memcache.Item{
		Key: keyName,
		Value: []byte(value),
	}
	memcache.Set(c, item)
	bumpGeneration(c)

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

	if err := memcache.Delete(c, keyName); err != nil {
		result["error"] += fmt.Sprintf("%s", err)
	}
	bumpGeneration(c)

	fmt.Fprintf(w, "%s", mapToJson(result))
}

func mapToJson(mapToConvert map[string] string) []byte {
	jsonResult, err := json.Marshal(mapToConvert)
        if err != nil {
                fmt.Printf("json marshalling saw error: %s\n", err)
        }

	return jsonResult
}

const generationKey = "GENERATION_NUMBER"

func bumpGeneration(c appengine.Context) {
	if item, err := memcache.Get(c, generationKey); err == memcache.ErrCacheMiss {
		newItem := &memcache.Item{
			Key: generationKey,
			Value: []byte("0"),
		}
		memcache.Set(c, newItem)
        } else {
		oldValue, _ := strconv.Atoi(fmt.Sprintf("%d", item.Value))
		newValue := int(oldValue) + 1
		item.Value = []byte(string(newValue))
		memcache.CompareAndSwap(c, item)
	}
}
