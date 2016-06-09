package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World")
}

func main() {
	http.HandleFunc("/", handler)

	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
