package web

import (
	"net/http"
)

type Web interface {
	Start(addr string) error
	GET(path string, handler func(c Context) error, middleware ...Use)
	POST(path string, handler func(c Context) error, middleware ...Use)
	PUT(path string, handler func(c Context) error, middleware ...Use)
	DELETE(path string, handler func(c Context) error, middleware ...Use)
	Group(path string, handler ...func(web Context) error) Web
	USE(handler ...Use) Web
}

type Context interface {
	Bind(i interface{}) error
	JSON(code int, i interface{}) error
	Request() *http.Request
	Param(key string) string
	Query(key string, typeData string) any
	Queries() map[string]string
	AllParams() map[string]string
}

type Use interface {
	Handle(c Context) error
}
