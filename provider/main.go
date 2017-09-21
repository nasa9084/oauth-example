package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	authnPageHTML = `<html>
<body>
client_id: %s<br />
AuthN?<br />

<form action="/authorize" method="POST">
<label for="id">ID: </label>
<input name="id" placeholder="userid">
<label for="passwd">Password: </label>
<input name="passwd" type="password" placeholder="passwd">
<input type="submit">
</form>
</body>
</html>`
	authzPageHTML = `<html>
<body>
AuthZ?<br />
<a href="/authorize/yes">yes</a><br />
<a href="/authorize/no">no</a>
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

func authnHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(authnPageHTML, r.FormValue(`client_id`))))
}

func authzHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case `GET`:
		w.Header().Set(`Location`, `/authenticate?client_id=` + r.FormValue(`client_id`))
		w.WriteHeader(http.StatusFound)
	case `POST`:
		w.Write([]byte(authzPageHTML))
	}
}

func authzYesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(`Location`, `http://localhost:8000/callback?code=authorizedyes`)
	w.WriteHeader(http.StatusFound)
}

func authzNoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(`Location`, `http://localhost:8000/callback?error=access_denied`)
	w.WriteHeader(http.StatusFound)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(`Content-Type`, `applicaiton/json`)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenJSON))
}
