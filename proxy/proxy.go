package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	PROXY_HEADER = "Proxy-Connection"
)

type RedirectHandler = func(*http.Request, []*http.Request) error

type Proxy struct {
	Port string
}

func NewProxy(port int) *Proxy {
	return &Proxy{
		Port: ":" + strconv.Itoa(port),
	}
}

func (pr *Proxy) HandleHttp(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = ""
	r.URL = &url.URL{Host: r.Host, Path: r.URL.Path}
	r.URL.Scheme = "http"
	r.Header.Del(PROXY_HEADER)

	client := http.Client{
		CheckRedirect: pr.handleRedirect(),
	}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (pr *Proxy) handleRedirect() RedirectHandler {
	return func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
}

func (pr *Proxy) HandleHttps(w http.ResponseWriter, r *http.Request) {

}
