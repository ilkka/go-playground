package hello

import (
	"fmt"
	"net/http"
	"appengine"
	"appengine/user"
	"html/template"
)

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
	fmt.Fprintf(w, guestBookForm)
}

const guestBookForm = `
<html>
<body>
	<form action="/sign" method="post">
		<label for="content">Content:</label>
		<textarea name="content" rows="3" cols="60"></textarea><br/>
		<button type="submit">Sign</button>
	</form>
</body>
</html>
`

func sign(w http.ResponseWriter, r *http.Request) {
	err := signTemplate.Execute(w, r.FormValue("content"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var signTemplate = template.Must(template.New("sign").Parse(signTemplateHtml))

const signTemplateHtml = `
<html>
<body>
	<p>You wrote: {{.}}</p>
</body>
</html>
`