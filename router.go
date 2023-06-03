package main

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type NewRouter struct {
	mux *http.ServeMux
}

var (
	auth       = NewAuthenticator(WithScopes("user-read-private", "user-read-email"))
	state      = GenerateRandomString(16)
	ch         = make(chan *Client, 10)
	login bool = false
)

func (r *NewRouter) AddHandler(path string, handler http.HandlerFunc) {
	r.mux.HandleFunc(path, handler)
}

func startServer() {
	r := &NewRouter{
		mux: http.NewServeMux(),
	}

	r.AddHandler("/", indexHandler)
	r.AddHandler("/login", loginHandler)
	r.AddHandler("/callback", callbackHandler)

	http.ListenAndServe(":3000", r.mux)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	if !login {
		fmt.Println("error: no client")
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	client := <-ch

	fmt.Println("Here")

	playlist, err := client.GetPlaylist(r.Context(), "2PQGFlUBOGUCq4N6JxNcda")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(playlist.Name)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := auth.AuthURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("show_dialog", "true"))
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := NewClient(auth.Client(r.Context(), token))
	select {
	case ch <- client:
		fmt.Println("client added")
	default:
		fmt.Println("client not added")
		return
	}
	login = true

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
