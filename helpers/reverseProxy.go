package helpers

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ReverseProxy proxies requests
func ReverseProxy(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("Request Info: ", r.Host, r.Method, r.Proto)
	//fmt.Println("14 ", r.Header)
	fmt.Println("15 in reverse proxy")

	rawUrl := "http://localhost:3000/"

	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		fmt.Println(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedUrl)
	proxy.ServeHTTP(w, r)

}
