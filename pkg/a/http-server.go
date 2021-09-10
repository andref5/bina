package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "A OK")
	})

	log.Println("HTTP server A startup 5011")
	http.ListenAndServe(":5011", nil)
}
