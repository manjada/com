package web

import (
	"bytes"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/manjada/com/config"
	_ "github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Fiber struct {
	App         *fiber.App
	GroupRouter fiber.Router
}

func (f *Fiber) USE(handler func(web Context) error) Web {
	if group, ok := f.GroupRouter.(*fiber.Group); ok {
		group.Use(func(c *fiber.Ctx) error {
			return handler(&FiberCtx{c})
		})
	} else {
		f.App.Use(func(c *fiber.Ctx) error {
			return handler(&FiberCtx{c})
		})
	}
	return f
}

func (f *Fiber) PUT(path string, handler func(c Context) error) {
	if group, ok := f.GroupRouter.(*fiber.Group); ok {
		group.Put(path, func(c *fiber.Ctx) error {
			return handler(&FiberCtx{c})
		})
	} else {
		f.App.Put(path, func(c *fiber.Ctx) error {
			return handler(&FiberCtx{c})
		})
	}

}

func (f *Fiber) DELETE(path string, handler func(c Context) error) {
	if group, ok := f.GroupRouter.(*fiber.Group); ok {
		group.Delete(path, func(c *fiber.Ctx) error {
			return handler(&FiberCtx{c})
		})
	} else {
		f.App.Delete(path, func(c *fiber.Ctx) error {
			return handler(&FiberCtx{c})
		})
	}
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
	if group, ok := f.GroupRouter.(*fiber.Group); ok {
		group.Get(path, func(c *fiber.Ctx) error {
			return handler(&FiberCtx{c})
		})
	} else {
		f.App.Get(path, func(c *fiber.Ctx) error {
			return handler(&FiberCtx{c})
		})
	}
}

func (f *Fiber) POST(path string, handler func(c Context) error) {
	if group, ok := f.GroupRouter.(*fiber.Group); ok {
		group.Post(path, func(c *fiber.Ctx) error {
			return handler(&FiberCtx{c})
		})
	} else {
		f.App.Post(path, func(c *fiber.Ctx) error {
			return handler(&FiberCtx{c})
		})
	}
}

func (f *Fiber) Group(path string, handler ...func(web Context) error) Web {
	group := f.App.Group(path)
	for _, h := range handler {
		group.Use(func(c *fiber.Ctx) error {
			return h(&FiberCtx{c})
		})
	}
	f.GroupRouter = group
	return f
}

func NewFiber() Web {
	f := fiber.New()
	return &Fiber{App: f}
}

func (f *Fiber) Start(addr string) error {
	return f.App.Listen(addr)
}

// JwtConfig returns a configuration struct for JWT middleware
func JwtConfigFiber() jwtware.Config {
	secretKey := config.GetConfig().AppJwt.AccessSecret
	return jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte(secretKey)},
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		},
	}
}

// FiberJwtMiddleware returns JWT middleware for Fiber framework
func FiberJwtMiddleware() fiber.Handler {
	configJwt := JwtConfigFiber()
	return jwtware.New(configJwt)
}
