package gmrouter

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func (r *Router) LambdaResponse(code int, data interface{}) events.APIGatewayProxyResponse {
	response := r.NewResponseMessage(code, data)
	rsp, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		Body:       string(rsp),
		StatusCode: code,
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Cache-Control": "public, max-age=30, stale-while-revalidate=10",
		},
	}
}

func (r *Router) LambdaResponseOk(data interface{}) events.APIGatewayProxyResponse {
	return r.LambdaResponse(http.StatusOK, data)
}

func (r *Router) LambdaResponseBadRequest() events.APIGatewayProxyResponse {
	return r.LambdaResponse(http.StatusBadRequest, nil)
}

func (r *Router) LambdaResponseUnauthorized() events.APIGatewayProxyResponse {
	return r.LambdaResponse(http.StatusUnauthorized, nil)
}

func (r *Router) LambdaResponseInternalServerError() events.APIGatewayProxyResponse {
	return r.LambdaResponse(http.StatusInternalServerError, nil)
}
