package api

import (
	"fmt"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/tangr/cicdgo/app/model"
	"github.com/tangr/cicdgo/app/service"
	"github.com/tangr/cicdgo/library/response"
)

type userApi struct{}

var User = &userApi{}

func (a *userApi) ListUsers(r *ghttp.Request) {
	users := service.User.ListUsers()
	params := g.Map{
		"url":         UrlPrefix + "/users/",
		"users":       users,
		"newUsersUrl": UrlPrefix + "/users/new",
	}
	r.Response.WriteTpl("users/list.html", params)
}

func (a *userApi) NewUser(r *ghttp.Request) {
	groups := service.Group.ListGroups()
	params := g.Map{
		"url":        UrlPrefix + "/users/",
		"newUserUrl": UrlPrefix + "/v1/users",
		"groups":     groups,
	}
	r.Response.WriteTpl("users/new.html", params)
}

func (a *userApi) ShowUser(r *ghttp.Request) {
	user_id := r.GetString("id")
	user_name := service.User.GetUserName(user_id)
	groups := service.Group.ListGroups()
	params := g.Map{
		"url":       UrlPrefix + "/users/",
		"apiurl":    UrlPrefix + "/v1/users/" + user_id,
		"user_name": user_name,
		"groups":    groups,
	}
	r.Response.WriteTpl("users/show.html", params)
}

func (a *userApi) New(r *ghttp.Request) {
	username := r.GetFormString("username")
	groups := r.GetFormStrings("groups")
	password := r.GetFormString("password")
	userid := service.User.New(username, groups, password)
	r.Response.RedirectTo(UrlPrefix + "/users/" + fmt.Sprint(userid))
}

func (a *userApi) Show(r *ghttp.Request) {
	user_id := r.GetString("id")
	user := service.User.GetUser(user_id)
	params := g.Map{
		"username": user.Email,
		"groups":   user.Group_Id,
	}
	r.Response.WriteJsonExit(params)
}

func (a *userApi) Update(r *ghttp.Request) {
	userid := r.GetString("id")
	username := r.GetFormString("username")
	groups := r.GetFormStrings("groups")
	password := r.GetFormString("password")
	_ = service.User.Update(userid, username, groups, password)
	r.Response.RedirectTo(UrlPrefix+"/users/"+userid, 303)
}

func (a *userApi) LoginPage(r *ghttp.Request) {
	params := g.Map{
		"apiurl": UrlPrefix + "/v1/users/login",
	}
	r.Response.WriteTpl("users/login.html", params)
}

func (a *userApi) ForbiddenPage(r *ghttp.Request) {
	r.Response.WriteTpl("users/forbidden.html")
}

func (a *userApi) SignIn(r *ghttp.Request) {
	var (
		data *model.UserApiSignInReq
	)

	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	if err := service.User.SignIn(r.Context(), data.Email, data.Password); err != nil {
		response.JsonExit(r, 1, err.Error())
	} else {
		r.Response.RedirectTo(UrlPrefix + "/")
	}
}

func (a *userApi) SignOut(r *ghttp.Request) {
	if err := service.User.SignOut(r.Context()); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	r.Response.RedirectTo(UrlPrefix + "/login")
}
