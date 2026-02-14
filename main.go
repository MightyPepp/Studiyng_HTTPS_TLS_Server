package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, HTTPS!")
	})

	log.Println("Сервер запущен на :8443")

	err := http.ListenAndServeTLS("localhost:8443", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера", err)
	}
}