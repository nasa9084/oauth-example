package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	cbHTML = `<html>
<body>
%s<br />
get Token: %s
</body>
</html>
`
	indexHTML = `<html>
<body>
<a href="/authz">Start AuthZ</button>
</body>
</html>`
)

func main() { os.Exit(exec()) }

// main logic
func exec() int {
	// binding handlers
	http.HandleFunc(`/`, indexHandler)
	http.HandleFunc(`/authz`, authzHandler)
	http.HandleFunc(`/callback`, cbHandler)

	// run http server
	addr := ":8000"
	log.Printf("client server is listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Printf("%s", err)
		return 1
	}
	return 0
}

// top page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(indexHTML))
}

func authzHandler(w http.ResponseWriter, r *http.Request) {
	// Authorization Request
	// OAuth 2.0 4.1.1, Appendix B
	query := url.Values{}
	query.Add(`response_type`, `code`) // REQUIRED, must be `code`
	query.Add(`client_id`, `client application`) // REQUIRED
	query.Add(`redirect_uri`, `http://localhost:8000/callback`) // OPTIONAL
	// in production, use TLS
	w.Header().Set(`Location`, `http://localhost:8080/authorize?` + query.Encode())
	w.WriteHeader(http.StatusFound)
}

// callback endpoint
func cbHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var body string
	var token string
	// Validate authorization code
	code := r.FormValue(`code`)
	if code == "" {
		err := r.FormValue(`error`)
		body = fmt.Sprintf(`Error: %s`, err)
	} else {
		// OAuth 2.0 4.1.3
		// access token request
		body = fmt.Sprintf(`Code: %s`, code)
		buf := bytes.Buffer{}
		buf.WriteString(`grant_type=authorization_code&code=`)
		buf.WriteString(code)
		res, _ := http.Post(`http://localhost:8080/token`, `application/x-www-form-urlencoded`, &buf)
		tokenb, _ := ioutil.ReadAll(res.Body)
		token = string(tokenb)
	}

	w.Write([]byte(fmt.Sprintf(cbHTML, body, token)))
}
