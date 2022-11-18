package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/mrjones/oauth"
	"github.com/twsiyuan/evernote-sdk-golang/edamutil"
)

const (
	evernoteKey    = "YOUR_KEY"
	evernoteSecret = "YOUR_SECRET"

	evernoteEnvironment = edamutil.SANDBOX

	sessionSecret  = "SESSION_SECRET"
	sessionContext = "SESSION_CONTEXT"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func main() {
	host := edamutil.Host(evernoteEnvironment)
	client := oauth.NewConsumer(
		evernoteKey,
		evernoteSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   fmt.Sprintf("https://%s/oauth", host),
			AuthorizeTokenUrl: fmt.Sprintf("https://%s/OAuth.action", host),
			AccessTokenUrl:    fmt.Sprintf("https://%s/oauth", host),
		},
	)

	println("Go to the url below to login: http://127.0.0.1:8080/evernote/login")

	gob.Register(&oauth.RequestToken{})

	m := http.NewServeMux()
	s := &http.Server{
		Addr:    ":8080",
		Handler: gcontext.ClearHandler(m),
	}

	sessionMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			session, err := store.Get(req, sessionSecret)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			req = req.WithContext(context.WithValue(req.Context(), sessionContext, session))
			next(w, req)
		}
	}

	m.HandleFunc("/evernote/login", sessionMiddleware(func(w http.ResponseWriter, req *http.Request) {
		session := req.Context().Value(sessionContext).(*sessions.Session)
		requestToken, url, err := client.GetRequestTokenAndUrl("http://127.0.0.1:8080/evernote/callback")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Values["requestToken"] = requestToken
		session.Save(req, w)

		http.Redirect(w, req, url, http.StatusTemporaryRedirect)
	}))

	m.HandleFunc("/evernote/callback", sessionMiddleware(func(w http.ResponseWriter, req *http.Request) {
		session := req.Context().Value(sessionContext).(*sessions.Session)
		requestToken, ok := session.Values["requestToken"].(*oauth.RequestToken)
		if !ok {
			http.Error(w, "Bad request, NO SESSION", http.StatusBadRequest)
			return
		}

		values := req.URL.Query()
		if token, ok := values["oauth_token"]; !ok || token[0] != requestToken.Token {
			http.Error(w, "Bad request, NO OAUTH_TOKEN", http.StatusBadRequest)
			return
		}

		verifier, ok := values["oauth_verifier"]
		if !ok {
			http.Error(w, "User decline", http.StatusOK)
			return
		}

		authorizedToken, err := client.AuthorizeToken(requestToken, verifier[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Values["authorizedToken"] = authorizedToken.Token
		session.Save(req, w)

		http.Redirect(w, req, "/evernote/", http.StatusTemporaryRedirect)
	}))

	m.HandleFunc("/evernote/", sessionMiddleware(func(w http.ResponseWriter, req *http.Request) {
		session := req.Context().Value(sessionContext).(*sessions.Session)

		authorizedToken, ok := session.Values["authorizedToken"].(string)
		if !ok {
			http.Redirect(w, req, "/evernote/login", http.StatusTemporaryRedirect)
			return
		}

		us, err := edamutil.NewUserStore(evernoteEnvironment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		evernoteauthorizedToken := authorizedToken
		ns, err := edamutil.NewNoteStore(req.Context(), us, evernoteauthorizedToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		notebook, _ := ns.GetDefaultNotebook(req.Context(), evernoteauthorizedToken)
		fmt.Fprintf(w, "Default notebook: %s\r\n\r\n", notebook.GetName())

		notebooks, _ := ns.ListNotebooks(req.Context(), evernoteauthorizedToken)
		for idx, notebook := range notebooks {
			fmt.Fprintf(w, "Notebook[%d]: %s\r\n", idx, notebook.GetName())
		}
	}))

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
