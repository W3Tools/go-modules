package gmrouter

import "net/http"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var DefaultMessage = map[int]string{
	http.StatusOK:                  "OK",
	http.StatusBadRequest:          "BadRequest",
	http.StatusInternalServerError: "Internal Server Error",
	http.StatusUnauthorized:        "Unauthorized",
}

func (*Router) NewResponseMessage(code int, data interface{}) Response {
	return Response{
		Code:    code,
		Message: DefaultMessage[code],
		Data:    data,
	}
}

func (r *Router) ResponseMessageOk(data interface{}) Response {
	return r.NewResponseMessage(http.StatusOK, data)
}

func (r *Router) ResponseMessageBadRequest(data interface{}) Response {
	return r.NewResponseMessage(http.StatusBadRequest, data)
}

func (r *Router) ResponseMessageUnauthorized(data interface{}) Response {
	return r.NewResponseMessage(http.StatusUnauthorized, data)
}

func (r *Router) ResponseMessageInternalServerError(data interface{}) Response {
	return r.NewResponseMessage(http.StatusInternalServerError, data)
}
