package proxy

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/Kudesnjk/http_proxy/cacher"
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
	r.Header.Del(PROXY_HEADER)
	client := http.Client{
		CheckRedirect: pr.handleRedirect(),
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

func (pr *Proxy) handleRedirect() RedirectHandler {
	return func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
}

func (pr *Proxy) HandleHttps(w http.ResponseWriter, r *http.Request) {
	r.Header.Del(PROXY_HEADER)
	dest_conn, err := net.DialTimeout("tcp", r.Host, pr.Timeout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	err = pr.Cacher.InsertRequest(r)
	if err != nil {
		fmt.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	go pr.transfer(dest_conn, client_conn)
	go pr.transfer(client_conn, dest_conn)
}

func (pr *Proxy) transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
