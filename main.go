package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

type ProxyHandler struct {
	proxy  *httputil.ReverseProxy
	target *url.URL
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	r.Host = h.target.Host
	h.proxy.ServeHTTP(w, r)
}

func main() {
	var targetURL string
	var listenPort string

	flag.StringVar(&targetURL, "target-url", "", "to forward HTTP requests to")
	flag.Parse()
	if targetURL == "" {
		log.Fatal("option -target-url is missing")
	}

	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("failed to parse target URL %s, %s", targetURL, err)
	}
	if target.Scheme != "https" {
		log.Fatalf("invalid target url %s, only HTTPS target urls are supported", targetURL)
	}

	listenPort = os.Getenv("PORT")
	if listenPort == "" {
		listenPort = "8080"
	}

	if port, err := strconv.ParseUint(listenPort, 10, 64); err != nil || port > 65535 {
		log.Fatalf("the environment variable PORT is not a valid port number")
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	http.Handle("/", &ProxyHandler{proxy: proxy, target: target})
	err = http.ListenAndServe(":"+listenPort, nil)
	if err != nil {
		log.Fatalf("server failed, %s", err)
	}
}
