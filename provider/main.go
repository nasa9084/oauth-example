package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	authnPageHTML = `<html>
<body>
client_id: %s<br />
%s<br />
AuthN?<br />

<form action="/authorize" method="POST">
<label for="id">ID: </label>
<input name="id" placeholder="userid">
<label for="passwd">Password: </label>
<input name="passwd" type="password" placeholder="passwd">
<input type="submit">

<input type="hidden" name="client_id" value="%s">
<input type="hidden" name="redirect_uri" value="%s">
</form>
</body>
</html>`
	authzPageHTML = `<html>
<body>
U R logged in.: %s<br />
AuthZ?<br />
<a href="/authorize/yes?%s">yes</a><br />
<a href="/authorize/no?%s">no</a>
</body>
</html>`
	tokenJSON = `{
  "access_token": "somethingToken",
  "token_type": "example"
}`
)

func main() { os.Exit(exec()) } // exec main logic and exit with status code 0/1

// main logic
func exec() int {
	// binding handlers
	http.HandleFunc(`/authenticate`, authnHandler)
	http.HandleFunc(`/authorize`, authzHandler)
	http.HandleFunc(`/authorize/yes`, authzYesHandler)
	http.HandleFunc(`/authorize/no`, authzNoHandler)
	http.HandleFunc(`/token`, tokenHandler)

	// run http server
	addr := ":8080"
	log.Printf("provider server is listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Printf("%s", err)
		return 1
	}
	return 0
}

// Authentication Endpoint
// This process is out of OAuth 2.0
// shows auth form
func authnHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.FormValue(`client_id`)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(authnPageHTML, clientID, r.URL.Query().Get(`error`), clientID, r.URL.Query().Get(`redirect_uri`))))
}

// Authorization Endpoint
func authzHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case `GET`: // GET method means the request is first Authorization request
		// redirect to authentication endpoint
		// this is out of OAuth 2.0
		query := url.Values{}
		query.Add(`client_id`, r.FormValue(`client_id`))
		query.Add(`redirect_uri`, r.FormValue(`redirect_uri`))
		w.Header().Set(`Location`, `/authenticate?`+query.Encode())
		w.WriteHeader(http.StatusFound)
	case `POST`: // POST method means the requrest is authenticate request
		// dummy authenticate process
		if authn(r.FormValue(`id`), r.FormValue(`passwd`)) {
			// shows Authorization form
			// resource owner chooses authoriza or not
			query := url.Values{}
			query.Add(`client_id`, r.FormValue(`client_id`))
			query.Add(`redirect_uri`, r.FormValue(`redirect_uri`))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(authzPageHTML, r.FormValue(`client_id`), query.Encode(), query.Encode())))
			return
		}
		// if authentication is failed, re-redirect to authentication form
		query := url.Values{}
		query.Add(`error`, `cantLogin`)
		query.Add(`client_id`, r.FormValue(`client_id`))
		query.Add(`redirect_uri`, r.FormValue(`redirect_uri`))
		w.Header().Set(`Location`, `/authenticate?`+query.Encode())
		w.WriteHeader(http.StatusFound)
	}
}

// dummy authentication
func authn(userid, password string) bool {
	return userid == `userid` && password == `passwd`
}

func authzYesHandler(w http.ResponseWriter, r *http.Request) {
	// OAuth 2.0 4.1.2
	// Successful Authorization Response
	query := url.Values{}
	// in production, this code must be invalidated in short time after(max 10 minutes is recommended)
	query.Add(`code`, `authorizedyes`) // Authorization code
	w.Header().Set(`Location`, r.FormValue(`redirect_uri`)+`?`+query.Encode())
	w.WriteHeader(http.StatusFound) // callback redirect
}

func authzNoHandler(w http.ResponseWriter, r *http.Request) {
	// OAuth 2.0 4.1.3.1
	// Error Response
	query := url.Values{}
	query.Add(`error`, `access_denied`) // Error code
	w.Header().Set(`Location`, r.FormValue(`redirect_uri`)+`?`+query.Encode())
	w.WriteHeader(http.StatusFound) // callback redirect
}

// Token Endpoint
func tokenHandler(w http.ResponseWriter, r *http.Request) {
	// in production, validate authorization code
	w.Header().Set(`Content-Type`, `applicaiton/json`)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenJSON))
}
