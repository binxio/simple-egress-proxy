//   Copyright 2021 binx.io B.V.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//
package main

import (
	"crypto/tls"
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
	var insecure bool
	var targetURL string
	var listenPort string

	flag.BoolVar(&insecure, "insecure", true, "allow expired and unknown host certificates")
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

	if insecure {
		proxy.Transport =
			&http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
			}
	}
	http.Handle("/", &ProxyHandler{proxy: proxy, target: target})
	err = http.ListenAndServe(":"+listenPort, nil)
	if err != nil {
		log.Fatalf("server failed, %s", err)
	}
}
