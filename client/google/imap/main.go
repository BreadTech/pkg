package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/coreos/go-oidc"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/google/uuid"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	goauth "google.golang.org/api/oauth2/v2"

	"github.com/breadtech/pkg/crypto/sasl/xoauth2"
)

var scopes = []string{
	goauth.UserinfoEmailScope,
	goauth.UserinfoProfileScope,
	goauth.OpenIDScope,
	gmail.MailGoogleComScope,
}

type OIDCClaims struct {
	Email    string `json:"email"`
	Verified bool   `json:"email_verified"`
}

func main() {
	// load creds
	credFile := "creds.json"
	if len(os.Args) > 1 {
		credFile = os.Args[1]
	}

	credBytes, err := os.ReadFile(credFile)
	if err != nil {
		panic(err)
	}

	cfg, err := google.ConfigFromJSON(credBytes, scopes...)
	if err != nil {
		panic(err)
	}

	// open browser with oauth request
	state := uuid.New().String()
	authURL := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	browser.OpenURL(authURL)

	// setup oidc provider and verifier
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		panic(err)
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})

	// handle oauth response
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		vals := r.URL.Query()
		if vals.Get("state") != state {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "bad")
			return
		}

		// Do code exchange.
		oauth2Token, err := cfg.Exchange(ctx, vals.Get("code"))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "oauth code exchange failed: %v", err)
			return
		}

		// Extract the ID Token from OAuth2 token.
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "oauth exchange extra id token missing")
			return
		}

		// Parse and verify ID Token payload.
		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "oauth exchange extra id token missing")
			return
		}

		// Extract custom claims
		claims := new(OIDCClaims)
		if err := idToken.Claims(claims); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "failed to extract claims from oidc verification token")
			return
		}

		// Authenticate to Gmail IMAP using SASL-XOAuth2
		c, err := client.DialTLS("imap.gmail.com:993", nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "failed to load email client")
			return
		}
		c.SetDebug(os.Stdout)
		defer c.Logout()

		sc := xoauth2.New(claims.Email, oauth2Token.AccessToken)
		if err := c.Authenticate(sc); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "failed to login email %v", err)
			return
		}

		// List mailboxes
		mailboxes := make(chan *imap.MailboxInfo, 10)
		done := make(chan error, 1)
		go func() {
			done <- c.List("", "*", mailboxes)
		}()

		resp := "Mailboxes:"
		for m := range mailboxes {
			resp += "* " + m.Name
		}

		if err := <-done; err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, resp)
	})
	socket, err := net.Listen("tcp", "127.0.0.1:2442")
	if err != nil {
		panic(err)
	}
	http.Serve(socket, mux)
}
