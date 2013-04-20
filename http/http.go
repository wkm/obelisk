package main

import (
	_ "circuit/load"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var (
	templates = template.Must(template.ParseGlob("templates/*"))
)

func main() {
	log.Printf("obelisk/")
	http.HandleFunc("/zk/", zkHandler)
	http.HandleFunc("/", indexHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("err: %s", err.Error())
	}
}

func error(rw http.ResponseWriter, msg string) {
	log.Printf("err: %s", msg)
	fmt.Fprintf(rw, "Error: %s", msg)
}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Obelisk")
}
