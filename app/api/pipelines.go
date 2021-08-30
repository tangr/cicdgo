package api

import (
	"fmt"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/tangr/cicdgo/app/service"
)

type pipelineApi struct{}

var Pipeline = &pipelineApi{}

func (a *pipelineApi) ListPipelines(r *ghttp.Request) {
	pipelines := service.Pipeline.ListPipelines()
	params := g.Map{
		"url":            UrlPrefix + "/pipelines/",
		"pipelines":      pipelines,
		"newPipelineUrl": UrlPrefix + "/pipelines/new",
	}
	r.Response.WriteTpl("pipelines/list.html", params)
}

func (a *pipelineApi) NewPipeline(r *ghttp.Request) {
	agents := service.Agent.GetAgentNames()
	groups := service.Group.GetGroupNames()

	params := g.Map{
		"url":    UrlPrefix + "/pipelines/",
		"apiurl": UrlPrefix + "/v1/pipelines/",
		"groups": groups,
		"agents": agents,
	}
	r.Response.WriteTpl("pipelines/new.html", params)
}

func (a *pipelineApi) ShowPipeline(r *ghttp.Request) {
	var pipeline_id int = r.GetInt("id")
	agents := service.Agent.GetAgentNames()
	groups := service.Group.GetGroupNames()

	pipeline_name := service.Pipeline.GetPipelineName(pipeline_id)
	params := g.Map{
		"url":           UrlPrefix + "/pipelines/",
		"apiurl":        UrlPrefix + "/v1/pipelines/" + fmt.Sprint(pipeline_id),
		"pipeline_name": pipeline_name,
		"pipeline_id":   pipeline_id,
		"agents":        agents,
		"groups":        groups,
	}
	r.Response.WriteTpl("pipelines/edit.html", params)
}

func (a *pipelineApi) New(r *ghttp.Request) {
	var pipeline_name, pipeline_body string

	pipeline_name = r.GetFormString("pipeline_name")
	group_id := r.GetFormInt("group_id")
	agent_id := r.GetFormInt("agent_id")
	concurrency := r.GetFormInt("concurrency")
	pipeline_body = r.GetFormString("pipeline_body")
	pipeline_id := service.Pipeline.New(pipeline_name, group_id, agent_id, concurrency, pipeline_body)

	r.Response.RedirectTo(UrlPrefix + "/pipelines/" + fmt.Sprint(pipeline_id))
}

func (a *pipelineApi) Show(r *ghttp.Request) {
	pipeline_id := r.GetInt("id")
	result := service.Pipeline.Show(pipeline_id)
	r.Response.WriteJsonExit(result)
}

func (a *pipelineApi) Update(r *ghttp.Request) {
	pipeline_id := r.GetInt("id")
	group_id := r.GetFormInt("group_id")
	agent_id := r.GetFormInt("agent_id")
	concurrency := r.GetFormInt("concurrency")
	pipeline_name := r.GetFormString("pipeline_name")
	pipeline_body := r.GetFormString("pipeline_body")
	err := service.Pipeline.Update(pipeline_id, group_id, agent_id, concurrency, pipeline_name, pipeline_body)
	if err != nil {
		g.Log().Error(err)
	}
	r.Response.RedirectTo(UrlPrefix+"/pipelines/"+fmt.Sprint(pipeline_id), 303)
}
