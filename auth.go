package main

import (
	"context"
	"errors"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

type Authenticator struct {
	config *oauth2.Config
}

type AuthenticatorOption func(a *Authenticator)

func WithScopes(scopes ...string) AuthenticatorOption {
	return func(a *Authenticator) {
		a.config.Scopes = scopes
	}

}

func NewAuthenticator(opts ...AuthenticatorOption) *Authenticator {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Endpoint:     spotify.Endpoint,
	}

	a := &Authenticator{
		config: conf,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a Authenticator) AuthURL(state string, opts ...oauth2.AuthCodeOption) string {
	return a.config.AuthCodeURL(state, opts...)
}

func (a Authenticator) Token(ctx context.Context, state string, r *http.Request, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	values := r.URL.Query()
	if e := values.Get("error"); e != "" {
		return nil, errors.New("spotify: authentication failed: " + e)
	}

	code := values.Get("code")
	if code == "" {
		return nil, errors.New("spotify: no code")
	}

	returedState := values.Get("state")
	if returedState != state {
		return nil, errors.New("spotify: state did not match")
	}

	return a.config.Exchange(ctx, code)
}

func (a Authenticator) Client(ctx context.Context, token *oauth2.Token) *http.Client {
	return a.config.Client(ctx, token)
}
