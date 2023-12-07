package gmrouter

func (r *Router) RequestHeaderSet(key string, value string) {
	r.Api.Request.Header.Set(key, value)
}

func (r *Router) RequestHeaderGet(key string) string {
	return r.Api.Request.Header.Get(key)
}
