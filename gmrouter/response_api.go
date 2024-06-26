package gmrouter

import "net/http"

func (r *Router) ApiResponse(code int, data interface{}) {
	r.ApiContext.Header("Access-Control-Allow-Origin", "*")
	r.ApiContext.Header("Access-Control-Allow-Methods", "*")
	r.ApiContext.Header("Access-Control-Allow-Headers", "*")

	r.ApiContext.JSON(code, r.NewResponseMessage(code, data))
}

func (r *Router) ApiResponseOk(data interface{}) {
	r.ApiResponse(http.StatusOK, data)
}

func (r *Router) ApiResponseBadRequest() {
	r.ApiResponse(http.StatusBadRequest, nil)
}

func (r *Router) ApiResponseUnauthorized() {
	r.ApiResponse(http.StatusUnauthorized, nil)
}

func (r *Router) ApiResponseInternalServerError() {
	r.ApiResponse(http.StatusInternalServerError, nil)
}
