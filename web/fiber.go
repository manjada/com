package web

import "github.com/gofiber/fiber/v2"

type Fiber struct {
	*fiber.App
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
