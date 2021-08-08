package api

import (
	"fmt"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/tangr/cicdgo/app/service"
)

type groupApi struct{}

var Group = &groupApi{}

func (a *groupApi) ListGroups(r *ghttp.Request) {
	groups := service.Group.ListGroups()
	params := g.Map{
		"url":         UrlPrefix + "/groups/",
		"groups":      groups,
		"newGroupUrl": UrlPrefix + "/groups/new",
	}
	r.Response.WriteTpl("groups/list.html", params)
}

func (a *groupApi) NewGroup(r *ghttp.Request) {
	params := g.Map{
		"url":         UrlPrefix + "/groups/",
		"newGroupUrl": UrlPrefix + "/v1/groups",
	}
	r.Response.WriteTpl("groups/new.html", params)
}

func (a *groupApi) ShowGroup(r *ghttp.Request) {
	group_id := r.GetString("id")
	group_name := service.Group.GetGroupName(group_id)
	params := g.Map{
		"url":        UrlPrefix + "/groups/",
		"apiurl":     UrlPrefix + "/v1/groups/" + group_id,
		"group_name": group_name,
	}
	r.Response.WriteTpl("groups/show.html", params)
}

func (a *groupApi) New(r *ghttp.Request) {
	var groupname string = r.GetFormString("groupname")
	groupid := service.Group.New(groupname)
	r.Response.RedirectTo(UrlPrefix + "/groups/" + fmt.Sprint(groupid))
}

func (a *groupApi) Show(r *ghttp.Request) {
	group_id := r.GetString("id")
	groupname := service.Group.GetGroupName(group_id)
	params := g.Map{
		"groupname": groupname,
	}
	r.Response.WriteJsonExit(params)
}

func (a *groupApi) Update(r *ghttp.Request) {
	groupid := r.GetString("id")
	groupname := r.GetFormString("groupname")
	_ = service.Group.Update(groupid, groupname)
	r.Response.RedirectTo(UrlPrefix+"/groups/"+groupid, 303)
}
