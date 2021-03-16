package cacher

import "net/http"

type Cacher interface {
	InsertRequest(r *http.Request) error
	GetRequests() ([]http.Request, error)
	GetSingleRequest(id int) (*http.Request, error)
}
