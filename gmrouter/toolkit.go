package gmrouter

func (r *Router) RequestHeaderSet(key string, value string) {
	r.ApiContext.Request.Header.Set(key, value)
}

func (r *Router) RequestHeaderGet(key string) string {
	return r.ApiContext.Request.Header.Get(key)
}
