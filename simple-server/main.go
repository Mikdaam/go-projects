package main

import (
	"fmt"
	"log"
	"net/http"
)

func aboutHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/about" {
		http.Error(w, "Error 404. Page not found\n", http.StatusNotFound)
		return
	}

	if req.Method != "GET" {
		http.Error(w, "Method not supported\n", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Welcome to the about page\n")
}

func formHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		fmt.Fprintf(w, "Can't parse form: %v\n", err)
	}
	firstname, lastname, email := req.FormValue("firstname"), req.FormValue("lastname"), req.FormValue("email")
	fmt.Fprintf(w, "Welcome %s %s, this is your email (%s).\n", firstname, lastname, email)
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))

	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/about", aboutHandler)

	fmt.Printf("Startting server at port 5050")
	if err := http.ListenAndServe(":5050", nil); err != nil {
		log.Fatal(err)
	}
}
