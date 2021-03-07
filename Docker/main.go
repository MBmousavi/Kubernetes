package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello!")
	})

	err := http.ListenAndServe(":8020", nil)
	if err == nil {
		panic(err)
	}
}
