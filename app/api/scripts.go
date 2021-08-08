package api

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/tangr/cicdgo/app/service"
)

type scriptApi struct{}

var Script = &scriptApi{}

func (a *scriptApi) ListScripts(r *ghttp.Request) {
	scripts := service.Script.ListScripts()
	params := g.Map{
		"url":          UrlPrefix + "/scripts/",
		"scripts":      scripts,
		"newScriptUrl": UrlPrefix + "/scripts/new",
	}
	r.Response.WriteTpl("scripts/list.html", params)
}

func (a *scriptApi) NewScript(r *ghttp.Request) {
	params := g.Map{
		"url":          UrlPrefix + "/scripts/",
		"newScriptUrl": UrlPrefix + "/v1/scripts",
	}
	r.Response.WriteTpl("scripts/new.html", params)
}

func (a *scriptApi) ShowScript(r *ghttp.Request) {
	script_id := r.GetString("id")
	script_name := service.Script.GetScriptName(script_id)
	params := g.Map{
		"url":         UrlPrefix + "/scripts/",
		"apiurl":      UrlPrefix + "/v1/scripts/" + script_id,
		"script_id":   script_id,
		"script_name": script_name,
	}
	r.Response.WriteTpl("scripts/show.html", params)
}

func (a *scriptApi) New(r *ghttp.Request) {
	var script_name, script_body string

	script_name = r.GetString("script_name")
	script_body = r.GetString("script_body")
	script_body = strings.Replace(script_body, "\r\n", "\n", -1)
	script_id := service.Script.New(script_name, script_body)

	r.Response.RedirectTo(UrlPrefix + "/scripts/" + fmt.Sprint(script_id))
}

func (a *scriptApi) Show(r *ghttp.Request) {
	script_id := r.GetInt("id")

	result := service.Script.Show(script_id)
	r.Response.WriteJsonExit(result)
}

func (a *scriptApi) Update(r *ghttp.Request) {
	script_id := r.GetInt("id")
	script_name := r.GetFormString("script_name")
	script_body := r.GetFormString("script_body")
	script_body = strings.Replace(script_body, "\r\n", "\n", -1)

	err := service.Script.Update(script_id, script_name, script_body)
	if err != nil {
		glog.Error(err)
	}

	r.Response.RedirectTo(UrlPrefix+"/scripts/"+fmt.Sprint(script_id), 303)
}
