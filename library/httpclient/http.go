package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/ernesto-jimenez/httplogger"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

type Helper interface {
	Do(httpReq *http.Request, resp any) (*http.Response, error)
	ErrorParseHTTPRequest(req *http.Request, err error) error
	ErrorParseHTTPResponse(res *http.Response, err error) error
	JoinPath(base string, elem ...string) string
	AnyToBuffer(v any) *bytes.Buffer
}

func New(log *otelzap.Logger) Helper {
	transport := otelhttp.NewTransport(
		&http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	)

	return &helper{
		log: log,
		client: &http.Client{
			Timeout: time.Duration(5) * time.Minute,
			Transport: httplogger.NewLoggedTransport(
				transport,
				&httpLogger{
					log,
				},
			),
		},
	}
}

type helper struct {
	log    *otelzap.Logger
	client *http.Client
}

func (h *helper) Do(httpReq *http.Request, resp any) (*http.Response, error) {
	httpReq.Header.Set("Content-type", "application/json")
	// TO DO: parameterize request headers
	httpResp, err := h.client.Do(httpReq)
	if err != nil {
		return httpResp, err
	}
	return httpResp, json.NewDecoder(httpResp.Body).Decode(resp)
}

func (h *helper) ErrorParseHTTPRequest(req *http.Request, err error) error {
	if req == nil {
		return err
	}
	return fmt.Errorf("[HTTP] Failed to do HTTP Request. Method: %s, url: %s, err: %v", req.Method, req.URL, err)
}

func (h *helper) ErrorParseHTTPResponse(res *http.Response, err error) error {
	if res == nil {
		return err
	}
	return fmt.Errorf("[HTTP] Failed to get Success HTTP Status Code. Method: %s, url: %s, status_code: %d, err: %v", res.Request.Method, res.Request.URL, res.StatusCode, err)
}

func (h *helper) JoinPath(base string, elem ...string) string {
	fullPath, err := url.JoinPath(base, elem...)
	if err != nil {
		h.log.Error("[HTTP] Error JoinPath", zap.Error(err))
	}
	return fullPath
}

func (h *helper) AnyToBuffer(v any) *bytes.Buffer {
	bData, err := json.Marshal(v)
	if err != nil {
		h.log.Error("[HTTP] Error Convert to Buffer", zap.Error(err))
	}
	return bytes.NewBuffer(bData)
}
