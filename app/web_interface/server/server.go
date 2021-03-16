package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/Kudesnjk/http_proxy/app/cacher"
	"github.com/Kudesnjk/http_proxy/app/proxy/proxy"
	"github.com/Kudesnjk/http_proxy/app/scanner"
	"github.com/labstack/echo/v4"
)

const (
	ERR_INTERNAL        = "Internal server error"
	ERR_BAD_REQUEST     = "Incorrect request received"
	ERR_REQUEST         = "Error occured while executing request"
	ERR_RESPONSE        = "Incorrect response received"
	ERR_QUERY_PARAM     = "Error occured while scanning request"
	ERR_NO_QUERY_PARAMS = "Request contains no query params"
	ERR_NO_REQUESTS     = "No requests sent with proxy"
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
	server.GET("/repeat/:id", wi.repeatRequestHandler())
	server.GET("/scan/:id", wi.scanRequestHandler())
	log.Fatalln(server.Start(address))
}

func (wi *WebInterface) getRequestsHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requests, err := wi.cacher.GetRequests()
		if err != nil {
			log.Println(err)
			return ctx.String(http.StatusInternalServerError, ERR_INTERNAL)
		}

		if len(requests) == 0 {
			return ctx.String(http.StatusOK, ERR_NO_REQUESTS)
		}

		var buffer bytes.Buffer
		for idx, req := range requests {
			strReq, _ := httputil.DumpRequest(&req, true)
			buffer.WriteString("ID: " + strconv.Itoa(idx) + "\n")
			buffer.Write(strReq)
		}

		return ctx.String(http.StatusOK, buffer.String())
	}
}

func (wi *WebInterface) getSingleRequestHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.String(http.StatusInternalServerError, ERR_BAD_REQUEST)
		}

		req, err := wi.cacher.GetSingleRequest(id)
		if err != nil {
			return ctx.String(http.StatusBadRequest, err.Error())
		}

		strReq, err := httputil.DumpRequest(req, true)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, ERR_INTERNAL)
		}

		return ctx.String(http.StatusOK, string(strReq))
	}
}

func (wi *WebInterface) repeatRequestHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.String(http.StatusInternalServerError, ERR_BAD_REQUEST)
		}

		req, err := wi.cacher.GetSingleRequest(id)
		if err != nil {
			return ctx.String(http.StatusBadRequest, err.Error())
		}

		client := http.Client{
			CheckRedirect: proxy.HandleRedirect(),
			Timeout:       time.Second * 10,
		}

		res, err := client.Do(req)
		if err != nil {
			return ctx.String(http.StatusBadRequest, ERR_REQUEST)
		}

		resStr, err := httputil.DumpResponse(res, true)
		if err != nil {
			return ctx.String(http.StatusBadRequest, ERR_RESPONSE)
		}

		return ctx.String(http.StatusOK, string(resStr))
	}
}

func (wi *WebInterface) scanRequestHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.String(http.StatusInternalServerError, ERR_BAD_REQUEST)
		}

		req, err := wi.cacher.GetSingleRequest(id)
		if err != nil {
			return ctx.String(http.StatusBadRequest, err.Error())
		}

		isVulMap, err := scanner.ScanRequest(req)
		if err != nil {
			return ctx.String(http.StatusBadRequest, ERR_QUERY_PARAM)
		}

		if len(isVulMap) == 0 {
			return ctx.String(http.StatusOK, ERR_NO_QUERY_PARAMS)
		}

		var buffer bytes.Buffer
		for queryParam, isVul := range isVulMap {
			if isVul {
				buffer.WriteString(fmt.Sprintf("Query param: '%s' is vulnerable!!!", queryParam))
			} else {
				buffer.WriteString(fmt.Sprintf("Query param: '%s' is not vulnerable.", queryParam))
			}
		}

		return ctx.String(http.StatusOK, buffer.String())
	}
}
