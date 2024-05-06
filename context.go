package gonvoy

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"sync"

	"github.com/ardikabs/gonvoy/pkg/types"
	"github.com/ardikabs/gonvoy/pkg/util"
	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
	"github.com/go-logr/logr"
)

type Context interface {
	// SetErrorHandler sets a custom error handler for an Http Filter
	//
	SetErrorHandler(ErrorHandler)

	// RegisterHandler adds an Http Filter Handler to the chain,
	// which should be run during filter startup (HttpFilter.OnStart).
	//
	RegisterHandler(HttpFilterHandler)

	// RequestHeader provides an interface to access and modify HTTP Request header, including
	// add, overwrite, or delete existing header.
	// RequestHeader will panic when it used without initializing the request header map first.
	//
	RequestHeader() Header

	// ResponseHeader provides an interface to access and modify HTTP Response header, including
	// add, overwrite, or delete existing header.
	// ResponseHeader will panic when it used without initializing the response header map first.
	//
	ResponseHeader() Header

	// RequestBodyBuffer provides an interface to access and manipulate an HTTP Request body.
	//
	RequestBody() Body

	// ResponseBody provides an interface access and manipulate an HTTP Response body.
	//
	ResponseBody() Body

	// Request returns an http.Request struct, which is a read-only data.
	// Attempting to modify this value will have no effect.
	// To make modifications to the request header, please use the RequestHeader() method instead.
	//
	Request() *http.Request

	// Response returns an http.Response struct, which is a read-only data.
	// It means, update anything to this value will result nothing.
	// To make modifications to the response header, please use the ResponseHeader() method instead.
	// To make modifications to the response body, please use the BufferWriter() method instead.
	//
	Response() *http.Response

	// SetRequestHeader is a low-level API, it set request header from RequestHeaderMap interface during DecodeHeaders phase
	//
	SetRequestHeader(api.RequestHeaderMap)

	// SetResponseHeader is a low-level API, it set response header from ResponseHeaderMap interface during EncodeHeaders phase
	//
	SetResponseHeader(api.ResponseHeaderMap)

	// SetRequestBody is a low-level API, it set request body from BufferInstance interface during DecodeData phase
	//
	SetRequestBody(api.BufferInstance)

	// SetResponseBody is a low-level API, it set response body from BufferInstance interface during EncodeData phase
	//
	SetResponseBody(api.BufferInstance)

	// Store allows you to save a value of any type under a key of any type.
	// It is designed for sharing data within a Context.
	// If you wish to share data throughout the lifetime of Envoy,
	// please refer to the Configuration.Cache.
	//
	// Please be cautious! The Store function overwrites any existing data.
	Store(key, value any)

	// Load retrieves a value associated with a specific key and assigns it to the receiver.
	// It is designed for sharing data within a Context.
	// If you wish to share data throughout the lifetime of Envoy,
	// please refer to the Configuration.Cache.
	//
	// It returns true if a compatible value is successfully loaded,
	// and false if no value is found or an error occurs during the process.
	Load(key any, receiver interface{}) (ok bool, err error)

	// Log provides a logger from the plugin to the Envoy Log. It accessible under Envoy `http` component.
	// e.g., Envoy flag `--component-log-level http:{debug,info,warn,error,critical}`
	//
	Log() logr.Logger

	// JSON sends a JSON response with a status code.
	//
	// This action halts the handler chaining and immediately returns back to Envoy.
	JSON(code int, b []byte, headers map[string][]string, opts ...ReplyOption) error

	// String sends a plain text response with a status code.
	//
	// This action halts the handler chaining and immediately returns back to Envoy.
	String(code int, s string, headers map[string][]string, opts ...ReplyOption) error

	// StatusType is a low-level API used to specify the type of status to be communicated to Envoy.
	//
	StatusType() api.StatusType

	// Committed indicates whether the current context has already completed its processing
	// within the plugin and forwarded the result to Envoy.
	//
	Committed() bool

	// StreamInfo offers an interface for retrieving comprehensive details about the incoming HTTP traffic, including
	// information such as the route name, filter chain name, dynamic metadata, and more.
	// It provides direct access to low-level Envoy information, so it's important to use it with a clear understanding of your intent.
	//
	StreamInfo() api.StreamInfo

	// Metrics sets gauge stats that could to record both increase and decrease metric. E.g., current active requests.
	//
	Metrics() Metrics

	// GetProperty is a helper function to fetch Envoy attributes based on https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes.
	// Currently, it only supports value that has a string format, work in progress for List/Map format.
	//
	GetProperty(name, defaultVal string) (string, error)

	// Configuration provides access to the filter configuration,
	// while also enabling users to persist and share data throughout Envoy's lifespan.
	//
	Configuration() Configuration

	// CanModifyRequestBody indicates whether an HTTP Request is eligible to modify the body/payload.
	//
	CanModifyRequestBody() bool

	// CanModifyResponseBody indicates whether an HTTP Response is eligible to modify the body/payload.
	//
	CanModifyResponseBody() bool

	//
	// --- Internal ---
	//

	// serveFilter serves the Http Filter for the specified phase.
	// This method is designed for internal use as it is directly invoked within each filter instance's phase.
	//
	serveFilter(HttpFilterPhase) api.StatusType
}

type context struct {
	config Configuration

	reqHeaderMap       api.RequestHeaderMap
	respHeaderMap      api.ResponseHeaderMap
	reqBufferInstance  api.BufferInstance
	respBufferInstance api.BufferInstance

	callback   api.FilterCallbacks
	statusType api.StatusType
	metrics    Metrics

	httpReq  *http.Request
	httpResp *http.Response

	storage   sync.Map
	logger    logr.Logger
	committed bool

	httpFilterManager HttpFilterHandlerManager
}

func newContext(cb api.FilterCallbacks, cfg Configuration) (Context, error) {
	if cb == nil {
		return nil, errors.New("filter callback can not be nil")
	}

	c := &context{
		callback: cb,
		logger:   NewLogger(cb),
		httpFilterManager: &DefaultHttpFilterHandlerManager{
			errorHandler: DefaultHttpFilterErrorHandler,
		},
	}

	type validator interface {
		Validate() error
	}

	if validate, ok := cfg.GetFilterConfig().(validator); ok {
		if err := validate.Validate(); err != nil {
			return nil, fmt.Errorf("invalid filter config; %w", err)
		}
	}

	cc := cfg.GetConfigCallbacks()
	if cc == nil {
		return nil, errors.New("config callbacks can not be nil")
	}

	c.config = cfg
	c.metrics = NewMetrics(cfg)
	return c, nil
}

func (c *context) SetErrorHandler(e ErrorHandler) {
	c.httpFilterManager.SetErrorHandler(e)
}

func (c *context) RegisterHandler(handler HttpFilterHandler) {
	c.httpFilterManager.RegisterHandler(handler)
}

func (c *context) serveFilter(phase HttpFilterPhase) api.StatusType {
	return c.httpFilterManager.serve(c, phase)
}

func (c *context) Configuration() Configuration {
	return c.config
}

func (c *context) Metrics() Metrics {
	return c.metrics
}

func (c *context) StreamInfo() api.StreamInfo {
	return c.callback.StreamInfo()
}

func (c *context) Log() logr.Logger {
	return c.logger
}

func (c *context) JSON(code int, body []byte, headers map[string][]string, opts ...ReplyOption) error {
	options := NewDefaultReplyOptions()
	for _, opt := range opts {
		opt(options)
	}

	if headers == nil {
		headers = make(map[string][]string)
	}

	if body == nil {
		body = []byte("{}")
	}

	headers["content-type"] = []string{"application/json"}
	c.callback.SendLocalReply(code, string(body), headers, options.grpcStatusCode, options.responseCodeDetails)
	c.committed = true
	c.statusType = options.statusType

	runtime.GC()
	return nil
}

func (c *context) String(code int, s string, headers map[string][]string, opts ...ReplyOption) error {
	options := NewDefaultReplyOptions()
	for _, opt := range opts {
		opt(options)
	}

	c.callback.SendLocalReply(code, s, headers, options.grpcStatusCode, options.responseCodeDetails)
	c.committed = true
	c.statusType = options.statusType

	return nil
}

func (c *context) RequestHeader() Header {
	if c.reqHeaderMap == nil {
		panic("Request Header has not been initialized yet")
	}

	return &header{c.reqHeaderMap}
}

func (c *context) ResponseHeader() Header {
	if c.respHeaderMap == nil {
		panic("Response Header has not been initialized yet")
	}

	return &header{c.respHeaderMap}
}

func (c *context) RequestBody() Body {
	if c.reqBufferInstance == nil {
		panic("Request Body has not been initialized yet.")
	}

	return &bodyWriter{
		header: c.reqHeaderMap,
		buffer: c.reqBufferInstance,
	}
}

func (c *context) ResponseBody() Body {
	if c.respBufferInstance == nil {
		panic("Response Body has not been initialized yet.")
	}

	return &bodyWriter{
		header: c.respHeaderMap,
		buffer: c.respBufferInstance,
	}
}

func (c *context) SetRequestHeader(header api.RequestHeaderMap) {
	c.reset()

	req, err := types.NewRequest(
		header.Method(),
		header.Host(),
		types.WithRequestURI(header.Path()),
		types.WithRequestHeaderRangeSetter(header),
	)
	if err != nil {
		c.Log().Error(err, "while initialize Http Request")
		return
	}

	c.httpReq = req
	c.reqHeaderMap = header
}

func (c *context) SetResponseHeader(header api.ResponseHeaderMap) {
	c.reset()

	code, ok := header.Status()
	if !ok {
		return
	}

	resp, err := types.NewResponse(code, types.WithResponseHeaderRangeSetter(header))
	if err != nil {
		c.Log().Error(err, "while initialize Http Response")
		return
	}

	c.httpResp = resp
	c.respHeaderMap = header
}

func (c *context) SetRequestBody(buffer api.BufferInstance) {
	c.reset()

	if c.reqBufferInstance == nil {
		c.reqBufferInstance = buffer
	}

	if buffer.Len() == 0 {
		return
	}

	buf := bytes.NewBuffer(buffer.Bytes())
	c.httpReq.Body = io.NopCloser(buf)
}

func (c *context) SetResponseBody(buffer api.BufferInstance) {
	c.reset()

	if c.respBufferInstance == nil {
		c.respBufferInstance = buffer
	}

	if buffer.Len() == 0 {
		return
	}

	buf := bytes.NewBuffer(buffer.Bytes())
	c.httpResp.Body = io.NopCloser(buf)
}

func (c *context) Request() *http.Request {
	if c.httpReq == nil {
		panic("Http Request is not yet initialized, see SetRequestHeader.")
	}

	return c.httpReq
}

func (c *context) Response() *http.Response {
	if c.httpResp == nil {
		panic("Http Response is not yet initialized, see SetResponse.")
	}

	return c.httpResp
}

func (c *context) StatusType() api.StatusType {
	return c.statusType
}

func (c *context) Committed() bool {
	return c.committed
}

func (c *context) reset() {
	c.statusType = api.Continue
	c.committed = false
}

func (c *context) CanModifyRequestBody() bool {
	contentLength, ok := c.reqHeaderMap.Get(HeaderContentLength)
	if !ok {
		return false
	}

	isEmpty := contentLength == "" || contentLength == "0"
	return !isEmpty
}

func (c *context) CanModifyResponseBody() bool {
	contentLength, ok := c.respHeaderMap.Get(HeaderContentLength)
	if !ok {
		return false
	}

	isEmpty := contentLength == "" || contentLength == "0"
	return !isEmpty
}

func (c *context) Store(key any, value any) {
	c.storage.Store(key, value)
}

func (c *context) Load(key any, receiver interface{}) (bool, error) {
	if receiver == nil {
		return false, errors.New("context: receiver should not be nil")
	}

	v, ok := c.storage.Load(key)
	if !ok {
		return false, nil
	}

	if !util.CastTo(receiver, v) {
		return false, errors.New("context: receiver and value has an incompatible type")
	}

	return true, nil
}

func (c *context) GetProperty(name, defaultVal string) (string, error) {
	value, err := c.callback.GetProperty(name)
	if err != nil {
		if errors.Is(err, api.ErrValueNotFound) {
			return defaultVal, nil
		}

		return value, err
	}

	if value == "" {
		return defaultVal, nil
	}

	return value, nil
}
