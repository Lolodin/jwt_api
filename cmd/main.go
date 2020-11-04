package main

import (
	"REST_API/controller"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/getToken", controller.GetTokens)


	http.ListenAndServe(":8080", mux)
}
