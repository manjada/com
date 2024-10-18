package web

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	_ "github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Fiber struct {
	*fiber.App
}

func (f *Fiber) PUT(path string, handler func(c Context) error) {
	f.App.Put(path, func(c *fiber.Ctx) error {
		return handler(&FiberCtx{c})
	})
}

func (f *Fiber) DELETE(path string, handler func(c Context) error) {
	f.App.Delete(path, func(c *fiber.Ctx) error {
		return handler(&FiberCtx{c})
	})
}

type FiberCtx struct {
	*fiber.Ctx
}

func (fc *FiberCtx) Bind(i interface{}) error {
	return fc.Ctx.BodyParser(i)
}

func (fc *FiberCtx) JSON(code int, i interface{}) error {
	return fc.Ctx.Status(code).JSON(i)
}

func (fc *FiberCtx) Param(key string) string {
	return fc.Ctx.Params(key)
}

func (fc *FiberCtx) Request() *http.Request {
	fasthttpReq := fc.Ctx.Request()
	req := &http.Request{
		Method:     string(fasthttpReq.Header.Method()),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Body:       ioutil.NopCloser(bytes.NewReader(fasthttpReq.Body())),
		Host:       string(fasthttpReq.Host()),
	}

	fasthttpReq.Header.VisitAll(func(key, value []byte) {
		req.Header.Set(string(key), string(value))
	})

	req.URL = &url.URL{
		Scheme:   string(fasthttpReq.URI().Scheme()),
		Host:     string(fasthttpReq.URI().Host()),
		Path:     string(fasthttpReq.URI().Path()),
		RawQuery: string(fasthttpReq.URI().QueryString()),
	}

	return req
}

func (f *Fiber) GET(path string, handler func(c Context) error) {
	f.App.Get(path, func(c *fiber.Ctx) error {
		return handler(&FiberCtx{c})
	})
}

func (f *Fiber) POST(path string, handler func(c Context) error) {
	f.App.Post(path, func(c *fiber.Ctx) error {
		return handler(&FiberCtx{c})
	})
}

func NewFiber() Web {
	f := fiber.New()
	return &Fiber{f}
}

func (f *Fiber) Start(addr string) error {
	return f.App.Listen(addr)
}
