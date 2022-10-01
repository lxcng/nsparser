package parser

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

const (
	nyaaRootUrl = "https://nyaa.si"
	lookupUrl   = nyaaRootUrl + "/user/%s?f=0&c=0_0&q=%s+%d&p=%d"
)

func getPage(user, title string, resoluiton int, p int) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	url := fmt.Sprintf(
		lookupUrl,
		user,
		strings.ReplaceAll(title, " ", "+"),
		resoluiton,
		p,
	)
	req.SetRequestURI(url)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := fasthttp.Do(req, resp)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}
