package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Kudesnjk/http_proxy/app/cacher"
)

const (
	PROXY_HEADER = "Proxy-Connection"
)

type RedirectHandler = func(*http.Request, []*http.Request) error

type Proxy struct {
	Port    string
	Timeout time.Duration
	Cacher  cacher.Cacher
}

func NewProxy(port int, timeout time.Duration, cacher cacher.Cacher) *Proxy {
	return &Proxy{
		Port:    ":" + strconv.Itoa(port),
		Timeout: timeout,
		Cacher:  cacher,
	}
}

func (pr *Proxy) HandleHttp(w http.ResponseWriter, r *http.Request) {
	if r.Proto != "HTTP/1.1" || r.Method == "CONNECT" {
		http.Error(w, "Proxy support HTTP only", http.StatusForbidden)
		return
	}

	url, _ := url.Parse(r.RequestURI)
	r.RequestURI = ""
	r.URL = url
	r.Header.Del(PROXY_HEADER)

	client := http.Client{
		CheckRedirect: HandleRedirect(),
		Timeout:       pr.Timeout,
	}

	res, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	err = pr.Cacher.InsertRequest(r)
	if err != nil {
		fmt.Println(err)
	}
}

func HandleRedirect() RedirectHandler {
	return func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
}
