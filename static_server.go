package main

import (
	"log"
	"net/http"
)

func staticServer(path string) {
	fs := http.FileServer(http.Dir(path))
	http.Handle("/", fs)

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":80", nil))
}
