package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/calculate", calcHandler)
	http.ListenAndServe(":8080", nil)
}
