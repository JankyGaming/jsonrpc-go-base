package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var ctx = context.Background()

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/public/api", publicHandler)
	r.HandleFunc("/private/api", privateHandler)

	err := http.ListenAndServe(":80", r)
	if err != nil {
		fmt.Println(err)
	}
}
