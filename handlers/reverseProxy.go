package handlers

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// ReverseProxy proxies requests
func ReverseProxy(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Request Info: ", r.Host, r.Method, r.Proto)

	var rawUrl string

	subdomain := strings.Split(r.Host, ".")[1]
	fmt.Println("21 subdomain ", subdomain)

	if subdomain == "pid1" {
		rawUrl = "http://localhost:8100/"
	} else {
		rawUrl = fmt.Sprintf("http://localhost:%s", server2)
	}

	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		fmt.Println(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedUrl)
	proxy.ServeHTTP(w, r)

}
