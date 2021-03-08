package cacher

import "net/http"

type Cacher interface {
	InsertRequest(r *http.Request) error
}
