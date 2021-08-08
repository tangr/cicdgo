package api

import (
	"fmt"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/tangr/cicdgo/app/service"
)

type agentApi struct{}

var Agent = &agentApi{}

func (a *agentApi) ListAgents(r *ghttp.Request) {
	agents := service.Agent.ListAgents()
	params := g.Map{
		"url":         UrlPrefix + "/agents/",
		"agents":      agents,
		"newAgentUrl": UrlPrefix + "/agents/new",
	}
	r.Response.WriteTpl("agents/list.html", params)
}

func (a *agentApi) NewAgent(r *ghttp.Request) {
	params := g.Map{
		"url":         UrlPrefix + "/agents",
		"newAgentUrl": UrlPrefix + "/v1/agents",
	}
	r.Response.WriteTpl("agents/new.html", params)
}

func (a *agentApi) ShowAgent(r *ghttp.Request) {
	agent_id := r.GetString("id")
	agent_name := service.Agent.GetAgentName(agent_id)

	params := g.Map{
		"url":        UrlPrefix + "/agents/",
		"apiurl":     UrlPrefix + "/v1/agents/" + agent_id,
		"agent_name": agent_name,
	}
	r.Response.WriteTpl("agents/edit.html", params)
}

func (a *agentApi) New(r *ghttp.Request) {
	var agent_name string = r.GetString("agent_name")
	agent_id := service.Agent.New(agent_name)
	r.Response.RedirectTo(UrlPrefix + "/agents/" + fmt.Sprint(agent_id))
}

func (a *agentApi) Show(r *ghttp.Request) {
	agent_id := r.GetInt("id")
	result := service.Agent.Show(agent_id)
	r.Response.WriteJsonExit(result)
}

func (a *agentApi) Update(r *ghttp.Request) {
	agent_id := r.GetInt("id")
	agent_name := r.GetFormString("agent_name")
	err := service.Agent.Update(agent_id, agent_name)
	if err != nil {
		glog.Error(err)
	}
	r.Response.RedirectTo(UrlPrefix+"/agents/"+fmt.Sprint(agent_id), 303)
}
