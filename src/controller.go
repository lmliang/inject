package mantis

import (
	"net/http"
)

type Handler interface {
	Handler(rw http.ResponseWriter, r *http.Request)
}

type Controller interface {
	Get()
	Post()
	Patch()
	Put()
	Delete()
	Options()
	Head()
	Render()
	Handler
}

type controller struct {
	Ctx *Context
}

func (ctrl *controller) Get() {
	ctrl.Ctx.Bodys = "Method Not Allow"
}

func (ctrl *controller) Post() {
	ctrl.Ctx.Bodys = "Method Not Allow"
}

func (ctrl *controller) Patch() {
	ctrl.Ctx.Bodys = "Method Not Allow"
}

func (ctrl *controller) Put() {
	ctrl.Ctx.Bodys = "Method Not Allow"
}

func (ctrl *controller) Delete() {
	ctrl.Ctx.Bodys = "Method Not Allow"
}

func (ctrl *controller) Options() {
	ctrl.Ctx.Bodys = "Method Not Allow"
}

func (ctrl *controller) Head() {
	ctrl.Ctx.Bodys = "Method Not Allow"
}

func (ctrl *controller) Render() {

}

func (ctrl *controller) Handler(rw http.ResponseWriter, r *http.Request) {
	ctrl.Ctx = NewContext(rw, r)

	switch r.Method {
	case "Get":
		ctrl.Get()
	case "Post":
		ctrl.Post()
	case "Patch":
		ctrl.Patch()
	case "Put":
		ctrl.Put()
	case "Delete":
		ctrl.Delete()
	case "Options":
		ctrl.Options()
	case "Head":
		ctrl.Head()
	default:
		ctrl.Ctx.Bodys = "Unkown Method " + r.Method
	}

	if len(ctrl.Ctx.Tmpl) > 0 {
		ctrl.Render()
	}

	ctrl.Ctx.Run()
}
