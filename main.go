package main

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {
	u, err := url.Parse(os.Getenv("PROXY_URL"))

	if err != nil {
		panic(err)
	}

	users := os.Getenv("USERS")

	userToPassword := make(map[string]string)

	for _, userAndPass := range strings.Split(users, ",") {
		parts := strings.Split(userAndPass, ":")

		userToPassword[parts[0]] = parts[1]
	}

	proxy := httputil.NewSingleHostReverseProxy(u)

	if os.Getenv("INSECURE") == "true" {
		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	http.ListenAndServe("0.0.0.0:8766", &BasicAuthProxy{
		users:   userToPassword,
		handler: proxy,
	})
}

type BasicAuthProxy struct {
	users   map[string]string
	handler http.Handler
}

func (ap *BasicAuthProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

	username, password, authOK := r.BasicAuth()
	if authOK == false {
		http.Error(w, "Not authorized", 401)
		return
	}

	usersPass, ok := ap.users[username]

	if !ok || usersPass != password {
		http.Error(w, "Not authorized", 401)
		return
	}

	ap.handler.ServeHTTP(w, r)
}
