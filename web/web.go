package web

import (
	"net/http"
)

type Web interface {
	Start(addr string) error
	GET(path string, handler func(c Context) error)
	POST(path string, handler func(c Context) error)
	PUT(path string, handler func(c Context) error)
	DELETE(path string, handler func(c Context) error)
	Group(path string, handler ...func(web Context) error) Web
	USE(handler ...Use) Web
}

type Context interface {
	Bind(i interface{}) error
	JSON(code int, i interface{}) error
	Request() *http.Request
	Param(key string) string
}

type Use interface {
	Handle(c Context) error
}
