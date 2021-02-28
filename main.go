package main

import (
	"log"
	"net/http"

	"github.com/Kudesnjk/http_proxy/proxy"
)

func main() {
	port := 8080
	proxy := proxy.NewProxy(port)

	server := http.Server{
		Addr: proxy.Port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				proxy.HandleHttps(w, r)
			} else {
				proxy.HandleHttp(w, r)
			}
		}),
	}

	log.Fatalln(server.ListenAndServe())
}
