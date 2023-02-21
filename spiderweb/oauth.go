package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  os.Getenv("OAUTH_CALLBACK_URL"),
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint: google.Endpoint,
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func generateStateOauthCookie(res http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration, Path: "/"}
	http.SetCookie(res, &cookie)
	return state
}

func oauthGoogleLogin(res http.ResponseWriter, req *http.Request) {
	oauthState := generateStateOauthCookie(res)
	u := googleOauthConfig.AuthCodeURL(oauthState)
	log.Println(req.RemoteAddr, "Login request")
	http.Redirect(res, req, u, http.StatusTemporaryRedirect)
}

func oauthGoogleCallback(res http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, cookieName)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	oauthState, _ := req.Cookie("oauthstate")

	if req.FormValue("state") != oauthState.Value {
		log.Println(req.RemoteAddr, "Invalid oauth google state")
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := getUserDataFromGoogle(req.FormValue("code"))
	if err != nil {
		log.Println(req.RemoteAddr, err.Error())
		http.Redirect(res, req, "/fail/", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(data, &googUser); err != nil {
		log.Println(err)
	}

	user := &User{
		Email:         strings.ToLower(googUser.Email),
		Username:      strings.Replace(strings.ToLower(googUser.Name), " ", "", -1),
		Forename:      strings.ToLower(googUser.GivenName),
		Surname:       strings.ToLower(googUser.FamilyName),
		Authenticated: true,
	}
	log.Println("Creating session for ", googUser.Email)
	session.Values["user"] = user
	session.Save(req, res)

	http.Redirect(res, req, "/", http.StatusFound)
	return
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
