package scanner

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/Kudesnjk/http_proxy/app/proxy/proxy"
)

const ATTACK_VECTOR = `vulnerable'"><img src onerror=alert()>`

func ScanRequest(req *http.Request) (map[string]bool, error) {
	if req == nil {
		return nil, errors.New("Request is nil")
	}

	client := http.Client{
		CheckRedirect: proxy.HandleRedirect(),
		Timeout:       time.Second * 10,
	}

	isVul := make(map[string]bool, len(req.URL.Query()))

	for queryParam, val := range req.URL.Query() {
		req.URL.Query().Set(queryParam, ATTACK_VECTOR)
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		strRes, err := httputil.DumpResponse(res, true)
		if err != nil {
			return nil, err
		}

		isVul[queryParam] = strings.Contains(string(strRes), ATTACK_VECTOR)
		req.URL.Query().Set(queryParam, strings.Join(val, ","))
	}

	return isVul, nil
}
