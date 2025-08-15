package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

var (
	clientID     = "6cdf8826-8ae3-4467-b477-9d3647ba6740"
	clientSecret = "arz8Q~DN6uNYV.vnSzA2sfwwrCh5JxNxoLNS3azu"
	tenantID     = "2d535af8-0d98-4df0-aa2f-8ff43c65dea7"
	redirectURL  = "http://localhost:8080/callback"

	oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", tenantID),
			TokenURL: fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID),
		},
	}
)

const style = `
	<style>
		body {
			font-family: Arial, sans-serif;
			background-color: #f5f7fa;
			text-align: center;
			padding: 40px;
		}
		h1 {
			color: #333;
		}
		a.button {
			display: inline-block;
			padding: 12px 24px;
			font-size: 16px;
			background-color: #0078D4;
			color: white;
			border-radius: 6px;
			text-decoration: none;
			transition: background-color 0.3s ease;
		}
		a.button:hover {
			background-color: #005A9E;
		}
		.token-box {
			margin: 20px auto;
			padding: 15px;
			width: 90%;
			max-width: 800px;
			background-color: #fff;
			border: 1px solid #ddd;
			border-radius: 6px;
			font-family: monospace;
			white-space: pre-wrap;
			word-break: break-all;
			text-align: left;
			box-shadow: 0px 2px 6px rgba(0,0,0,0.1);
		}
	</style>
`

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", callbackHandler)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<html><head>%s</head><body>
		<h1>Welcome to My Go App</h1>
		<a class="button" href="/login">Login with Azure AD</a>
	</body></html>`, style)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != "state" {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Token exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	idToken := token.Extra("id_token")
	if idToken != nil {
		fmt.Fprintf(w, `<html><head>%s</head><body>
			<h1>Authentication Successful</h1>
			<h2>Raw ID Token</h2>
			<div class="token-box">%s</div>
			<a class="button" href="/">Back to Home</a>
		</body></html>`, style, idToken)
	} else {
		fmt.Fprintf(w, `<html><head>%s</head><body>
			<h1>No ID token found</h1>
			<a class="button" href="/">Back to Home</a>
		</body></html>`, style)
	}
}
