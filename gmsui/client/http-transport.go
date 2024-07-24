package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

func (client *SuiClient) request(input SuiTransportRequestOptions, output interface{}) error {
	reflectValue := reflect.ValueOf(output)
	if output != nil && reflectValue.Kind() != reflect.Pointer {
		return fmt.Errorf("output not a pointer or nil pointer")
	}

	message, err := client.newRequestMessage(input.Method, input.Params)
	if err != nil {
		return err
	}

	jsb, err := json.Marshal(message)
	if err != nil {
		return err
	}

	httpRequest, err := http.NewRequestWithContext(client.ctx, http.MethodPost, client.rpc, bytes.NewReader(jsb))
	if err != nil {
		return err
	}
	httpRequest.ContentLength = int64(len(jsb))
	httpRequest.Header.Set("Content-Type", "application/json")

	response, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %v, status: %v", response.StatusCode, response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var v jsonrpcResponse
	err = json.Unmarshal(body, &v)
	if err != nil {
		return err
	}

	if v.Error != nil {
		return fmt.Errorf("unexpected code: %d, message: %v", v.Error.Code, string(v.Error.Message))
	}

	return json.Unmarshal(v.Result, &output)
}

type jsonrpcRequest struct {
	Jsonrpc string          `json:"jsonrpc,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
}

func (client *SuiClient) newRequestMessage(method string, params []any) (*jsonrpcRequest, error) {
	id, err := json.Marshal(client.requestId)
	if err != nil {
		return nil, err
	}

	requestMessage := &jsonrpcRequest{Jsonrpc: "2.0", ID: id, Method: method}
	if !reflect.ValueOf(params).IsNil() {
		requestMessage.Params, err = json.Marshal(params)
		if err != nil {
			return nil, err
		}
	}

	return requestMessage, nil
}

type jsonrpcResponse struct {
	Jsonrpc string          `json:"jsonrpc,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  json.RawMessage `json:"result"`
	Error   *jsonrpcError   `json:"error,omitempty"`
}

type jsonrpcError struct {
	Code    int64           `json:"code,omitempty"`
	Message json.RawMessage `json:"message,omitempty"`
}
