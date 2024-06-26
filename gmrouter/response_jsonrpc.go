package gmrouter

import (
	"encoding/json"
	"net/http"

	gm "github.com/W3Tools/go-modules"
)

const DefaultJsonRPCVersion = "2.0"

// JSON-RPC request message structure
type JsonRPCRequest struct {
	ID      json.RawMessage `json:"id"`
	Jsonrpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

// JSON-RPC response message structure
type JsonRPCResponse struct {
	ID      json.RawMessage `json:"id,omitempty"`
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
	OK             int = 200
	InvalidRequest int = -32600
	MethodNotFound int = -32601
	InvalidParams  int = -32602
	InternalError  int = -32603
	ParseError     int = -32700
)

var DefaultJsonRPCMessage = map[int]string{
	InvalidRequest: "Invalid request",
	MethodNotFound: "Method not found",
	InvalidParams:  "Invalid params",
	InternalError:  "Internal server error",
	ParseError:     "Parse error",
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

func (*Router) NewJsonRPCResponseMessage(id json.RawMessage, code int, data interface{}) JsonRPCResponse {
	response := JsonRPCResponse{Jsonrpc: DefaultJsonRPCVersion, ID: id}

	if code == OK {
		response.Result = data
	} else {
		response.Error = &JsonRPCError{
			Code:    code,
			Message: DefaultJsonRPCMessage[code],
			Data:    data,
		}
	}
	return response
}

func (r *Router) JsonRPCResponse(id json.RawMessage, code int, data interface{}) {
	r.ApiContext.Header("Access-Control-Allow-Origin", "*")
	r.ApiContext.Header("Access-Control-Allow-Methods", "*")
	r.ApiContext.Header("Access-Control-Allow-Headers", "*")

	r.ApiContext.JSON(http.StatusOK, r.NewJsonRPCResponseMessage(id, code, data))
}

func (r *Router) JsonRPCResponseOk(id json.RawMessage, data interface{}) {
	r.JsonRPCResponse(id, http.StatusOK, data)
}

func (r *Router) JsonRPCResponseInvalidRequest(id json.RawMessage, data interface{}) {
	r.JsonRPCResponse(id, InvalidRequest, data)
}

func (r *Router) JsonRPCResponseMethodNotFound(id json.RawMessage, data interface{}) {
	r.JsonRPCResponse(id, MethodNotFound, data)
}

func (r *Router) JsonRPCResponseInvalidParams(id json.RawMessage, data interface{}) {
	r.JsonRPCResponse(id, InvalidParams, data)
}

func (r *Router) JsonRPCResponseInternalError(id json.RawMessage, data interface{}) {
	r.JsonRPCResponse(id, InternalError, data)
}

func (r *Router) JsonRPCResponseParseError(id json.RawMessage, data interface{}) {
	r.JsonRPCResponse(id, ParseError, data)
}
