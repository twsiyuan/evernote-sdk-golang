package utils

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/mrjones/oauth"
)

// this is a simple example to implement oauth flow
func TestOauth(t *testing.T) {
	host := GetHost(SANDBOX)
	client := oauth.NewConsumer(
		EvernoteKey, EvernoteSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   fmt.Sprintf("https://%s/oauth", host),
			AuthorizeTokenUrl: fmt.Sprintf("https://%s/OAuth.action", host),
			AccessTokenUrl:    fmt.Sprintf("https://%s/oauth", host),
		},
	)

	requestToken, url, err := client.GetRequestTokenAndUrl("http://127.0.0.1:8082/")
	if err != nil {
		t.Fatal(err)
	}
	println("Go to the url below to login")
	println(url)

	verifier := ""
	s := &http.Server{Addr: ":8082"}
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		values := req.URL.Query()
		if token, ok := values["oauth_token"]; ok && token[0] == requestToken.Token {
			if ver, ok := values["oauth_verifier"]; ok {
				w.Write(([]byte)("OK"))
				verifier = ver[0]
				s.Shutdown(context.Background())
				return
			}
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(([]byte)("Bad request"))
	})
	s.Handler = h
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		t.Fatal(err)
	}

	authorizedToken, err := client.AuthorizeToken(requestToken, verifier)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("AuthorizedToken: %s", authorizedToken.Token)
}
