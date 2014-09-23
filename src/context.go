package mantis

import (
	"io"
	"net/http"
)

type Context struct {
	Input  map[string]string
	Output map[string]string
	Bodys  string
	Tmpl   string
	Data   map[interface{}]interface{}
	Req    *http.Request
	Rw     http.ResponseWriter
}

func (ctx *Context) Run() {
	io.WriteString(ctx.Rw, ctx.Bodys)
}

func (ctx *Context) Redirect(url string, status int) {
	ctx.Rw.Header().Set("Location", url)
	ctx.Rw.WriteHeader(status)
}

func NewContext(rw http.ResponseWriter, r *http.Request) *Context {
	input := make(map[string]string)
	output := make(map[string]string)
	data := make(map[interface{}]interface{})

	return &Context{Input: input, Output: output, Data: data, Req: r, Rw: rw}
}
