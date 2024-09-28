package web

type Web interface {
	GET(path string, handler func(c Context) error)
	POST(path string, handler func(c Context) error)
}

type Context interface {
	Bind(i interface{}) error
	JSON(code int, i interface{}) error
}
