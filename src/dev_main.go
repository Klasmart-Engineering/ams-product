// +build !lambda

package main

import (
	"net/http"

	"bitbucket.org/calmisland/product-lambda-funcs/src/server"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	err := handleHTTPRequest(w, r)
	if err != nil {
		panic(err)
	}
}

func main() {
	server.Setup()
	initLambdaFunctions()
	initLambdaDevFunctions()

	http.HandleFunc("/", handleRequest)

	httpServer := &http.Server{
		Addr:    ":8044",
		Handler: nil,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
