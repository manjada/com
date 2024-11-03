package web

import (
	"bytes"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

func (f *Fiber) USE(handlers ...Use) Web {
	if group, ok := f.GroupRouter.(*fiber.Group); ok {
		for _, handler := range handlers {
			group.Use(func(c *fiber.Ctx) error {
				return handler.Handle(&FiberCtx{c})
			})
		}
	} else {
		for _, handler := range handlers {
			f.App.Use(func(c *fiber.Ctx) error {
				return handler.Handle(&FiberCtx{c})
			})
		}
	}
	return f
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

func (f *Fiber) PUT(path string, handler func(c Context) error, middleware ...Use) {
	if group, ok := f.GroupRouter.(*fiber.Group); ok {
		group.Put(path, func(c *fiber.Ctx) error {
			return handle(c, middleware, handler)
		})
	} else {
		f.App.Put(path, func(c *fiber.Ctx) error {
			return handle(c, middleware, handler)
		})
	}

}

func (f *Fiber) DELETE(path string, handler func(c Context) error, middleware ...Use) {
	if group, ok := f.GroupRouter.(*fiber.Group); ok {
		group.Delete(path, func(c *fiber.Ctx) error {
			return handle(c, middleware, handler)
		})
	} else {
		f.App.Delete(path, func(c *fiber.Ctx) error {
			return handle(c, middleware, handler)
		})
	}
}

func (f *Fiber) GET(path string, handler func(c Context) error, middleware ...Use) {
	if group, ok := f.GroupRouter.(*fiber.Group); ok {
		group.Get(path, func(c *fiber.Ctx) error {
			return handle(c, middleware, handler)
		})
	} else {
		f.App.Get(path, func(c *fiber.Ctx) error {
			return handle(c, middleware, handler)
		})
	}
}

func (f *Fiber) POST(path string, handler func(c Context) error, middleware ...Use) {
	if group, ok := f.GroupRouter.(*fiber.Group); ok {
		group.Post(path, func(c *fiber.Ctx) error {
			return handle(c, middleware, handler)
		})
	} else {
		f.App.Post(path, func(c *fiber.Ctx) error {
			return handle(c, middleware, handler)
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

func handle(c *fiber.Ctx, middleware []Use, handler func(c Context) error) error {
	ctx := &FiberCtx{c}
	for _, a := range middleware {
		if err := a.Handle(ctx); err != nil {
			return err
		}
	}
	return handler(ctx)
}

func NewFiber() Web {
	f := fiber.New()
	return &Fiber{App: f}
}

func (f *Fiber) Start(addr string) error {
	return f.App.Listen(addr)
}

// JwtConfig returns a configuration struct for JWT middleware
func jwtConfigFiber() jwtware.Config {
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
func FiberJwtMiddleware() Use {
	configJwt := jwtConfigFiber()
	return NewJwtMiddleware(jwtware.New(configJwt))
}

type FiberMiddleware struct {
	handler fiber.Handler
}

func NewJwtMiddleware(handler fiber.Handler) *FiberMiddleware {
	return &FiberMiddleware{handler: handler}
}

func (m *FiberMiddleware) Handle(c Context) error {
	return m.handler(c.(*FiberCtx).Ctx)
}

func FiberCorsConfig() Use {
	corsConfig := initCorsConfig()
	return NewJwtMiddleware(cors.New(corsConfig))
}

func initCorsConfig() cors.Config {
	corsConfig := cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-CSRF-Token",
	}
	return corsConfig
}
