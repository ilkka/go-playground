package hello

import (
	"net/http"
	"appengine"
	"appengine/user"
	"appengine/datastore"
	"html/template"
	"time"
)

type Greeting struct {
	Author string
	Content string
	Timestamp time.Time
}

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/sign", sign)
}

func root(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
	q := datastore.NewQuery("Greeting").Order("-Timestamp").Limit(10)
	greetings := make([]Greeting, 0, 10)
	if _, err := q.GetAll(c, &greetings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := guestBookTemplate.Execute(w, greetings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)		
	}
}

var guestBookTemplate = template.Must(template.New("book").Parse(guestBookTemplateHtml))

const guestBookTemplateHtml = `
<html>
<body>
	{{range .}}
		{{with .Author}}
		<p><b>{{.}}</b> wrote:</p>
		{{else}}
		<p>Anonymous wrote:</p>
		{{end}}
		<p>{{.Content}}</p>
	{{end}}
	<form action="/sign" method="post">
		<label for="content">Content:</label>
		<textarea name="content" rows="3" cols="60"></textarea><br/>
		<button type="submit">Sign</button>
	</form>
</body>
</html>
`

func sign(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	g := Greeting {
		Content: r.FormValue("content"),
		Timestamp: time.Now(),
	}
	if u := user.Current(c); u != nil {
		g.Author = u.String()
	}
	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Greeting", nil), &g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

