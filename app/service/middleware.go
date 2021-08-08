package service

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/tangr/cicdgo/app/model"
)

var Middleware = middlewareService{}

type middlewareService struct{}

var UrlPrefix string = g.Cfg().GetString("server.console.UrlPrefix")

func (s *middlewareService) Ctx(r *ghttp.Request) {
	customCtx := &model.Context{
		Session: r.Session,
	}
	Context.Init(r, customCtx)
	if user := Session.GetUser(r.Context()); user != nil {
		customCtx.User = &model.ContextUser{
			Id:    user.Id,
			Email: user.Email,
		}

	}
	r.Middleware.Next()
}

func (s *middlewareService) Authen(r *ghttp.Request) {
	if User.IsSignedIn(r.Context()) {
		r.Middleware.Next()
	} else {
		r.Response.RedirectTo(UrlPrefix + "/login")
	}
}

func (s *middlewareService) Author(r *ghttp.Request) {
	if User.IsAdmin(r.Context()) {
		r.Middleware.Next()
	} else {
		r.Response.RedirectTo(UrlPrefix + "/forbidden")
	}
}

func (s *middlewareService) AuthorAdmin(r *ghttp.Request) {
	if User.IsAdmin(r.Context()) {
		r.Middleware.Next()
	} else {
		r.Response.RedirectTo(UrlPrefix + "/forbidden")
	}
}

func (s *middlewareService) CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}
