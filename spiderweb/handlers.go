package main

import (
	"fmt"
	"log"
	"net/http"
	"spiderweb/kubernetes"
	"strings"
)

func index(res http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, cookieName)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Println(req.RemoteAddr, err.Error())
		return
	}
	user := getUser(session)
	log.Println(req.RemoteAddr, user)

	flusher, ok := res.(http.Flusher)
	if !ok {
		http.Error(res, "ERROR: Server does not support Flusher!", http.StatusInternalServerError)
		log.Println(req.RemoteAddr, "ERROR: Server does not support Flusher!")
		return
	}

	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")

	// Check if the user has a valid session authenticated
	if auth := user.Authenticated; !auth {
		log.Println(req.RemoteAddr, "Session not authenticated, redirecting to sign on")
		err = session.Save(req, res)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			log.Println(req.RemoteAddr, "ERROR", err.Error())
			return
		}
		http.Redirect(res, req, "/auth/google/login", http.StatusFound)
		return
	}

	// Look for existing deployments and see if they match up to the current user.
	// Iterate over deployments in the namespace until we find a match
	// TODO: Serialize the email and use that as a deployment id? Makes OPS in cluster harder though, but could use labels.

	var i int = 0
	for _, d := range kubernetes.GetDeployments(namespace).Items {
		if cleanEmail(user.Email) == d.Name {
			i = 1
			fmt.Printf("Found existing deployment: * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
			break
		}
	}

	// if there's no deployment that matches, create one
	if i == 0 {
		log.Println(req.RemoteAddr, user.Email, "Creating deployment for user %s", user.Email)
		fmt.Fprintf(res, provisionWaitingRoom(user.Email, 0))
		flusher.Flush()
		createEnvironment(user)
	}

	// Backend for the given user.
	// TODO: Use Serialization on the email.
	backend := getBackendURL(cleanEmail(user.Email))
	beRes, err := http.Get(backend)
	if err != nil {
		log.Println(err)
	}

	if beRes.StatusCode != 200 {
		log.Println(req.RemoteAddr, beRes.StatusCode)
		flusher.Flush()
	}

	for beRes.StatusCode != 200 {
		log.Println(req.RemoteAddr, beRes.StatusCode)
		beRes, err = http.Get(backend)
		if err != nil {
			log.Println(err)
		}
	}

	// Ok we should have a valid environment, so punch a hole to it for the user.
	serveReverseProxy(backend, res, req)
}

func provisionWaitingRoom(email string, status int) string {
	// It aint pretty but it works!
	var body strings.Builder
	body.WriteString("<html>")
	body.WriteString("<head><title>spIDEr</title></head>")
	body.WriteString("<body style='background-color:#3c3c3c;'>")
	body.WriteString("<p style='font-family:Courier New; color:limegreen;'>")
	body.WriteString("<b>Please wait while your environment is privisioned</b>")
	body.WriteString("<br><br>")
	body.WriteString("Creating deployment for user:")
	body.WriteString(email)
	body.WriteString("</p>")
	return body.String()
}

func logout(res http.ResponseWriter, req *http.Request) {
	// Logout handler, this destroys the deployment.
	session, err := store.Get(req, cookieName)
	if err != nil {
		//http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Println(req.RemoteAddr, err.Error())
		return
	}

	user := getUser(session)
	destroyEnvironment(user)

	// Invalidate the session cookie
	session.Options.MaxAge = -1
	err = session.Save(req, res)
	if err != nil {
		//http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Println(req.RemoteAddr, err.Error())
		return
	}
	fmt.Fprintf(res, "<html><head><title>SpIDEr</title></head><body>")
	fmt.Fprintf(res, "Logged Out")
	fmt.Fprintf(res, "</body></html>")
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}

func int32Ptr(i int32) *int32 { return &i }
func int64Ptr(i int64) *int64 { return &i }
