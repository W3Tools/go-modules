package gmrouter

import "github.com/gin-gonic/gin/binding"

func (r *Router) BindJSON(obj any) error {
	return r.Api.ShouldBindWith(obj, binding.JSON)
}

func (r *Router) BindQuery(obj any) error {
	return r.Api.ShouldBindWith(obj, binding.Query)
}

func (r *Router) BindXML(obj any) error {
	return r.Api.ShouldBindWith(obj, binding.XML)
}

func (r *Router) BindYAML(obj any) error {
	return r.Api.ShouldBindWith(obj, binding.YAML)
}

func (r *Router) BindTOML(obj any) error {
	return r.Api.ShouldBindWith(obj, binding.TOML)
}
