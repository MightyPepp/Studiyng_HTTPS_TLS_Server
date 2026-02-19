package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var usersPasswords = map[string][]byte {
	"joe": 	[]byte("$2a$10$8kfCgrw46NKF79hfwmQTRuMi/EG5VYORx2fOTuqHlTUSfp4irKMNm"), //Actual password - "112233"
	"mary": []byte("$2a$10$O.CMUNt0vN5GoWijbFmqBe7I/lMtnoU8Gvs5eLr10ynvZ37f6ZaK2"), //Actual password - "SuperSecretPassword0987"
}

func verifyUserPass(username, password string) bool {
	wantPass, hasUser := usersPasswords[username]
	if !hasUser {
		return false
	} 
	if cmperr := bcrypt.CompareHashAndPassword(wantPass, []byte(password)); cmperr == nil {
		return true
	}
	return false
}

func main() {
	addr := flag.String("addr", ":8443", "HTTPS network addres")
	certFile := flag.String("certfile", "cert.pem", "certificate PEM file")
	keyFile := flag.String("keyfile", "key.pem", "key PEM file")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "С гордостью представляю свой TLS сервер")
	})

	mux.HandleFunc("/secret/", func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if ok && verifyUserPass(user, pass) {
			fmt.Fprintf(w, "U get to see the secret!")
		} else {
			w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
			http.Error(w, "Unautgirized", http.StatusUnauthorized)
		}
	})

	srv := &http.Server {
		Addr: 		*addr,
		Handler: 	mux,
		TLSConfig: 	&tls.Config{
					MinVersion: tls.VersionTLS13,
					PreferServerCipherSuites: true,
		},
	}

	log.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServeTLS(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}
}