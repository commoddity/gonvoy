package main

import (
	"net/http"

	"github.com/commoddity/gonvoy"
)

func init() {
	gonvoy.RunHttpFilter(
		customErrResponseFilterName,
		func() gonvoy.HttpFilter {
			return new(Filter)
		},
		gonvoy.ConfigOptions{},
	)
}

func main() {}

const customErrResponseFilterName = "custom_err_response"

<<<<<<< HEAD
=======
type Filter struct{}

>>>>>>> filter-factory-func-instead-of-reflect
func (Filter) OnBegin(c gonvoy.RuntimeContext, ctrl gonvoy.HttpFilterController) error {
	ctrl.AddHandler(Handler{})
	return nil
}

func (Filter) OnComplete(c gonvoy.Context) error {
	return nil
}

type Handler struct {
	gonvoy.PassthroughHttpFilterHandler
}

func (h Handler) OnResponseHeader(c gonvoy.Context) error {
	switch code := c.Response().StatusCode; code {
	case http.StatusUnauthorized:
		headers := gonvoy.NewGatewayHeadersWithEnvoyHeader(c.ResponseHeader())
		rcdetails := gonvoy.MustGetProperty(c, "response.code_details", gonvoy.DefaultResponseCodeDetails)

		return c.JSON(code, gonvoy.NewMinimalJSONResponse("UNAUTHORIZED", "UNAUTHORIZED"),
			gonvoy.LocalReplyWithHTTPHeaders(headers),
			gonvoy.LocalReplyWithRCDetails(rcdetails))

	case http.StatusTooManyRequests:
		headers := gonvoy.NewGatewayHeadersWithEnvoyHeader(c.ResponseHeader())
		rcdetails := gonvoy.MustGetProperty(c, "response.code_details", gonvoy.DefaultResponseCodeDetails)

		return c.JSON(code, gonvoy.NewMinimalJSONResponse("TOO_MANY_REQUESTS", "TOO_MANY_REQUESTS"),
			gonvoy.LocalReplyWithHTTPHeaders(headers),
			gonvoy.LocalReplyWithRCDetails(rcdetails))

	case http.StatusServiceUnavailable:
		headers := gonvoy.NewGatewayHeadersWithEnvoyHeader(c.ResponseHeader())
		rcdetails := gonvoy.MustGetProperty(c, "response.code_details", gonvoy.DefaultResponseCodeDetails)

		return c.JSON(code,
			gonvoy.NewMinimalJSONResponse("SERVICE_UNAVAILABLE", "SERVICE_UNAVAILABLE"),
			gonvoy.LocalReplyWithHTTPHeaders(headers),
			gonvoy.LocalReplyWithRCDetails(rcdetails))
	}

	return nil
}
