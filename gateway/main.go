package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/mux"
)

func reverseProxyHandler(targetURL string) http.Handler {
	target, _ := url.Parse(targetURL)
	return httputil.NewSingleHostReverseProxy(target)
}
func main() {
	router := mux.NewRouter()
	backendService1URL := fmt.Sprintf("http://%s:8081", os.Getenv("service1"))
	backendService2URL := fmt.Sprintf("http://%s:8082", os.Getenv("service2"))

	// Create routes for the API gateway
	router.PathPrefix("/team").Handler(reverseProxyHandler(backendService1URL))
	router.PathPrefix("/match").Handler(reverseProxyHandler(backendService2URL))
	http.Handle("/", router)
	http.ListenAndServe(":8080", router)
}
