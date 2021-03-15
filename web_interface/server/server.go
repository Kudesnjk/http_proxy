package server

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/Kudesnjk/http_proxy/cacher"
	"github.com/labstack/echo/v4"
)

type WebInterface struct {
	cacher cacher.Cacher
}

func NewWebInterface(cacher cacher.Cacher) *WebInterface {
	return &WebInterface{
		cacher: cacher,
	}
}

func (wi *WebInterface) RunWebInterface(address string) {
	server := echo.New()
	server.GET("/requests", wi.getRequestsHandler())
	server.GET("/requests/:id", wi.getSingleRequestHandler())
	server.POST("/repeat/:id", wi.repeatRequestHandler())
	server.GET("/scan/:id", wi.scanRequestHandler())
	log.Fatalln(server.Start(address))
}

func (wi *WebInterface) getRequestsHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requests, err := wi.cacher.GetRequests()
		if err != nil {
			return echo.ErrInternalServerError
		}

		var buffer bytes.Buffer
		for idx, req := range requests {
			strReq, _ := httputil.DumpRequest(&req, true)
			buffer.WriteString("ID: " + strconv.Itoa(idx) + "\n")
			buffer.Write(strReq)
		}

		if err != nil {
			return echo.ErrInternalServerError
		}

		return ctx.String(http.StatusOK, buffer.String())
	}
}

func (wi *WebInterface) getSingleRequestHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return nil
	}
}

func (wi *WebInterface) repeatRequestHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return nil
	}
}

func (wi *WebInterface) scanRequestHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return nil
	}
}
