package main

import (
        "fmt"
        "net/http"
)

func main() {
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                fmt.Fprintf(w, "hello!")
        })

        http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("content-type", "application/json")
                fmt.Fprintf(w, `{"status":"UP"}`)
        })

        err := http.ListenAndServe(":8020", nil)
        if err == nil {
                panic(err)
        }
}
