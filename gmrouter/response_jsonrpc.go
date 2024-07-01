package gmrouter

import (
	"encoding/json"
	"net/http"

	gm "github.com/W3Tools/go-modules"
)

const DefaultJsonRPCVersion = "2.0"

// JSON-RPC request message structure
type JsonRPCRequest struct {
	ID      json.RawMessage `json:"id" validate:"required"`
	Jsonrpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

// JSON-RPC response message structure
type JsonRPCResponse struct {
	ID      json.RawMessage `json:"id"`
	Jsonrpc string          `json:"jsonrpc"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *JsonRPCError   `json:"error,omitempty"`
}

// JSON-RPC error message structure
type JsonRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	OK             string = "ok"
	InvalidRequest string = "Invalid request"
	MethodNotFound string = "Method not found"
	InvalidParams  string = "Invalid params"
	NoMoreParams   string = "No more params"
	InternalError  string = "Internal server error"
	ParseError     string = "Parse error"
)

var DefaultJsonRPCCode = map[string]int{
	OK:             200,
	InvalidRequest: -32600,
	MethodNotFound: -32601,
	InvalidParams:  -32602,
	NoMoreParams:   -32602,
	InternalError:  -32603,
	ParseError:     -32700,
}

func (r *Router) JsonRPCShouldBindJSON() (request *JsonRPCRequest, err error) {
	request = new(JsonRPCRequest)
	if err := r.ApiContext.ShouldBindJSON(&request); err != nil {
		return nil, err
	}

	if err := gm.ValidateStruct(request); err != nil {
		return nil, err
	}
	return request, nil
}

func (*Router) NewJsonRPCResponseMessage(id json.RawMessage, code int, msg string, data interface{}) JsonRPCResponse {
	response := JsonRPCResponse{Jsonrpc: DefaultJsonRPCVersion, ID: id}

	if msg == OK || code == DefaultJsonRPCCode[OK] {
		response.Result = data
	} else {
		response.Error = &JsonRPCError{
			Code:    code,
			Message: msg,
			Data:    data,
		}
	}
	return response
}

func (r *Router) JsonRPCResponse(id json.RawMessage, code int, msg string, data interface{}) {
	r.ApiContext.Header("Access-Control-Allow-Origin", "*")
	r.ApiContext.Header("Access-Control-Allow-Methods", "*")
	r.ApiContext.Header("Access-Control-Allow-Headers", "*")

	r.ApiContext.JSON(http.StatusOK, r.NewJsonRPCResponseMessage(id, code, msg, data))
}

func (r *Router) JsonRPCResponseOk(id json.RawMessage, data interface{}) {
	r.JsonRPCResponse(id, DefaultJsonRPCCode[OK], OK, data)
}

func (r *Router) JsonRPCResponseInvalidRequest(id json.RawMessage, data interface{}) {
	r.JsonRPCResponse(id, DefaultJsonRPCCode[InvalidRequest], InvalidRequest, data)
}

func (r *Router) JsonRPCResponseMethodNotFound(id json.RawMessage, data interface{}) {
	r.JsonRPCResponse(id, DefaultJsonRPCCode[MethodNotFound], MethodNotFound, data)
}

func (r *Router) JsonRPCResponseInvalidParams(id json.RawMessage, data interface{}) {
	r.JsonRPCResponse(id, DefaultJsonRPCCode[InvalidParams], InvalidParams, data)
}

func (r *Router) JsonRPCResponseInternalError(id json.RawMessage, data interface{}) {
	r.JsonRPCResponse(id, DefaultJsonRPCCode[InternalError], InternalError, data)
}

func (r *Router) JsonRPCResponseParseError(id json.RawMessage, data interface{}) {
	r.JsonRPCResponse(id, DefaultJsonRPCCode[ParseError], ParseError, data)
}

func (r *Router) JsonRPCResponseNoMoreParams(id json.RawMessage, data interface{}) {
	r.JsonRPCResponse(id, DefaultJsonRPCCode[NoMoreParams], NoMoreParams, data)
}

type JsonRPCHandlerFunc func(*Router, *JsonRPCRequest)

func WrapperJsonRPCHandler(router *Router, request *JsonRPCRequest, handlers ...JsonRPCHandlerFunc) {
	for _, handler := range handlers {
		handler(router, request)
		if router.ApiContext.IsAborted() {
			return
		}
	}
}
