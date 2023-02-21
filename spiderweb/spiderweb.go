package main

import (
	"encoding/gob"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"spiderweb/kubernetes"
)

type User struct {
	Email         string
	Username      string
	Forename      string
	Surname       string
	Authenticated bool
}

type userData struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	Hd            string `json:"hd"`
}

var (
	store     *sessions.CookieStore
	googUser  *userData
	image     string
	namespace string
	appname   string
	port      string
	nfserver  string
	fqdn      string
)

const cookieName string = "spiderweb-app"

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	gob.Register(User{})
}

func main() {
	log.Print("SpIDErweb Startup")

	image = os.Getenv("SPIDER_IMAGE")
	namespace = os.Getenv("SPIDER_NAMESPACE")
	appname = os.Getenv("SPIDER_APPNAME")
	port = ":" + os.Getenv("SPIDERWEB_LISTEN_PORT")
	nfserver = os.Getenv("SPIDER_NFS_SERVER")
	fqdn = os.Getenv("SPIDER_FQDN")

	router := mux.NewRouter()
	router.HandleFunc("/auth/google/login", oauthGoogleLogin)
	router.HandleFunc("/auth/google/callback", oauthGoogleCallback)
	router.HandleFunc("/logout/", logout)
	router.HandleFunc("/favicon.ico", faviconHandler)
	router.HandleFunc("/{rest:.*}", index)
	log.Print("Listening on port ", port)
	http.ListenAndServe(port, router)
}

func getBackendURL(username string) string {
	name := cleanName(username)
	newBeUrl := "http://" + name + ":3000"
	log.Println(username, newBeUrl)
	return newBeUrl
}

func cleanName(givenName string) string {
	name := strings.ToLower(strings.Replace(givenName, " ", "", -1))
	return name
}

func cleanEmail(email string) string {
	var e string = email
	e = strings.ToLower(strings.Replace(e, "@", "-", -1))
	e = strings.ToLower(strings.Replace(e, ".", "-", -1))
	return e
}

func getUser(s *sessions.Session) User {
	val := s.Values["user"]
	var user = User{}
	user, ok := val.(User)
	if !ok {
		return User{Authenticated: false}
	}
	return user
}

func destroyEnvironment(user User) {
	kubernetes.DeleteDeployment(namespace, cleanEmail(user.Email))
	kubernetes.DeleteService(namespace, cleanEmail(user.Email))
}

func createEnvironment(user User) {
	// Create the kubernetes service to open a port to the IDE
	serviceSpec := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cleanEmail(user.Email),
			Namespace: namespace,
			Labels: map[string]string{
				"app":      appname,
				"forename": user.Forename,
				"surname":  user.Surname,
				"email":    cleanEmail(user.Email),
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{{
				Port:       3000,
				TargetPort: intstr.FromInt(3000),
			}},
			Selector: map[string]string{
				"app":      appname,
				"forename": user.Forename,
				"surname":  user.Surname,
				"email":    cleanEmail(user.Email),
			},
		},
	}

	// Create the deployment of the IDE
	deploymentSpec := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cleanEmail(user.Email),
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":      appname,
					"forename": user.Forename,
					"surname":  user.Surname,
					"email":    cleanEmail(user.Email),
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":      appname,
						"forename": user.Forename,
						"surname":  user.Surname,
						"email":    cleanEmail(user.Email),
					},
				},
				Spec: apiv1.PodSpec{
					SecurityContext: &apiv1.PodSecurityContext{
						RunAsUser:  int64Ptr(1001),
						RunAsGroup: int64Ptr(1001),
						FSGroup:    int64Ptr(2000),
					},
					Containers: []apiv1.Container{
						{
							Name:  user.Username,
							Image: image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 3000,
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "home",
									MountPath: "/home/coder",
									SubPath:   appname + "/" + cleanEmail(user.Email),
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "USER",
									Value: user.Forename,
								},
								{
									Name:  "EMAIL",
									Value: user.Email,
								},
								{
									Name:  "USER_FORENAME",
									Value: user.Forename,
								},
								{
									Name:  "USER_SURNAME",
									Value: user.Surname,
								},
								{
									Name:  "THEIA_HOSTS",
									Value: fqdn,
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "home",
							VolumeSource: apiv1.VolumeSource{
								NFS: &apiv1.NFSVolumeSource{
									Server:   nfserver,
									Path:     "/",
									ReadOnly: false,
								},
							},
						},
					},
				},
			},
		},
	}
	kubernetes.CreateService(namespace, serviceSpec)
	kubernetes.CreateDeployment(namespace, deploymentSpec)
}
