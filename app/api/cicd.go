package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gpage"
	"github.com/tangr/cicdgo/app/service"
)

type cicdApi struct{}

var Cicd = &cicdApi{}

func (a *cicdApi) ListCicd(r *ghttp.Request) {
	group_ids := service.GetUserGroupIds(r.Context())
	pipelines := service.Cicd.ListCicd(group_ids)
	params := g.Map{
		"url":            UrlPrefix + "/",
		"pipelines":      pipelines,
		"newPipelineUrl": UrlPrefix + "/pipelines/new",
	}
	r.Response.WriteTpl("cicd/list.html", params)
}

func pageContent(page *gpage.Page) string {
	page.NextPageTag = `<i class="angle right icon"></i>`
	page.PrevPageTag = `<i class="angle left icon"></i>`
	pageStr := page.PrevPage()
	pageStr += fmt.Sprint(page.CurrentPage)
	pageStr += page.NextPage()
	return pageStr
}

func (a *cicdApi) ShowCicd(r *ghttp.Request) {
	var pipeline_id int = r.GetInt("pipeline_id")
	if !service.CheckAuthor(r.Context(), pipeline_id) {
		r.Response.RedirectTo(UrlPrefix + "/forbidden")
	}
	pageid := r.GetQueryInt("page")
	jobs, totalSize := service.Cicd.GetJobs(pipeline_id, pageid, 10)
	pipeline_name := service.Pipeline.GetPipelineName(pipeline_id)
	page := r.GetPage(totalSize, 10)
	params := g.Map{
		"url":           UrlPrefix + "/" + fmt.Sprint(pipeline_id),
		"apiurl":        UrlPrefix + "/v1/" + fmt.Sprint(pipeline_id) + "/body",
		"newJobUrl":     UrlPrefix + "/v1/" + fmt.Sprint(pipeline_id) + "/newjob",
		"pkgurl":        UrlPrefix + "/v1/" + fmt.Sprint(pipeline_id) + "/pkgs",
		"pipeline_name": pipeline_name,
		"pipeline_id":   pipeline_id,
		"jobs":          jobs,
		"page":          pageContent(page),
		"envurl":        UrlPrefix + "/v1/" + fmt.Sprint(pipeline_id) + "/",
	}
	r.Response.WriteTpl("cicd/show.html", params)
}

func (a *cicdApi) ShowJob(r *ghttp.Request) {
	var pipeline_id int = r.GetInt("pipeline_id")
	if !service.CheckAuthor(r.Context(), pipeline_id) {
		r.Response.RedirectTo(UrlPrefix + "/forbidden")
	}
	var job_id int = r.GetInt("job_id")
	tasks := service.Cicd.GetJobTasks(pipeline_id, job_id)
	pipeline_name := service.Pipeline.GetPipelineName(pipeline_id)
	concurrency, job_type, job_status := service.Cicd.GetJobInfo(job_id)
	params := g.Map{
		"url":           UrlPrefix + "/" + fmt.Sprint(pipeline_id) + "/",
		"apiurl":        UrlPrefix + "/v1/" + fmt.Sprint(pipeline_id, "/", job_id),
		"pipeline_name": pipeline_name,
		"pipeline_id":   pipeline_id,
		"job_id":        job_id,
		"concurrency":   concurrency,
		"job_type":      job_type,
		"job_status":    job_status,
		"tasks":         tasks,
		"taskurl":       UrlPrefix + "/v1/" + fmt.Sprint(pipeline_id) + "/",
	}
	if job_type == "BUILD" {
		r.Response.WriteExit(params)
		r.Response.WriteTpl("cicd/job_build.html", params)
	} else {
		params["progressurl"] = UrlPrefix + "/v1/" + fmt.Sprint(pipeline_id, "/", job_id) + "/progress"
		r.Response.WriteTpl("cicd/job_deploy.html", params)
	}
}

func (a *cicdApi) GetLog(r *ghttp.Request) {
	var pipeline_id int = r.GetInt("pipeline_id")
	if !service.CheckAuthor(r.Context(), pipeline_id) {
		r.Response.WriteStatus(http.StatusForbidden)
	}
	var log_id int = r.GetInt("task_id")
	output := service.Cicd.GetOutput(pipeline_id, log_id)
	r.Response.WriteExit(output)
}

func (a *cicdApi) GetJobEnvs(r *ghttp.Request) {
	var pipeline_id int = r.GetInt("pipeline_id")
	if !service.CheckAuthor(r.Context(), pipeline_id) {
		r.Response.WriteStatus(http.StatusForbidden)
	}
	var job_id int = r.GetInt("job_id")
	envs := service.Cicd.GetJobEnvs(pipeline_id, job_id)
	r.Response.WriteExit(envs)
}

func (a *cicdApi) GetJobProgress(r *ghttp.Request) {
	var pipeline_id int = r.GetInt("pipeline_id")
	if !service.CheckAuthor(r.Context(), pipeline_id) {
		r.Response.WriteStatus(http.StatusForbidden)
	}
	var job_id int = r.GetInt("job_id")
	// tasks := service.Cicd.GetJobTasks(pipeline_id, job_id)
	task_total, task_value, job_finished := service.Cicd.GetJobProgress(pipeline_id, job_id)
	job_progress := g.Map{
		"total":    task_total,
		"value":    task_value,
		"finished": job_finished,
	}
	r.Response.WriteExit(job_progress)
	// job_tasks := service.Cicd.GetJobProgress(pipeline_id, job_id)
	// r.Response.WriteExit(job_tasks)
}

func (a *cicdApi) PostJobConcurrency(r *ghttp.Request) {
	var pipeline_id int = r.GetInt("pipeline_id")
	if !service.CheckAuthor(r.Context(), pipeline_id) {
		r.Response.WriteStatus(http.StatusForbidden)
	}
	var job_id int = r.GetInt("job_id")
	var concurrency int = r.GetFormInt("concurrency")
	service.Cicd.PostJobConcurrency(pipeline_id, job_id, concurrency)
	r.Response.RedirectTo(UrlPrefix + "/" + fmt.Sprint(pipeline_id, "/", job_id))
}

func (a *cicdApi) PostJobStatus(r *ghttp.Request) {
	var pipeline_id int = r.GetInt("pipeline_id")
	if !service.CheckAuthor(r.Context(), pipeline_id) {
		r.Response.WriteStatus(http.StatusForbidden)
	}
	var job_id int = r.GetInt("job_id")
	var job_status string = r.GetString("status")
	service.Cicd.PostJobStatus(pipeline_id, job_id, job_status)
	r.Response.RedirectTo(UrlPrefix + "/" + fmt.Sprint(pipeline_id, "/", job_id))
}

func (a *cicdApi) AbortTask(r *ghttp.Request) {
	var pipeline_id int = r.GetInt("pipeline_id")
	if !service.CheckAuthor(r.Context(), pipeline_id) {
		r.Response.WriteStatus(http.StatusForbidden)
	}
	var task_id int = r.GetInt("task_id")
	var job_id int = r.GetInt("job_id")
	var clientip string = r.GetString("clientip")
	service.Cicd.AbortTask(pipeline_id, task_id, job_id, clientip)
	r.Response.RedirectTo(UrlPrefix + "/" + fmt.Sprint(pipeline_id, "/", job_id))
}

func (a *cicdApi) RetryTask(r *ghttp.Request) {
	var pipeline_id int = r.GetInt("pipeline_id")
	if !service.CheckAuthor(r.Context(), pipeline_id) {
		r.Response.WriteStatus(http.StatusForbidden)
	}
	var task_id int = r.GetInt("task_id")
	var job_id int = r.GetInt("job_id")
	var clientip string = r.GetString("clientip")
	service.Cicd.RetryTask(pipeline_id, task_id, job_id, clientip)
	r.Response.RedirectTo(UrlPrefix + "/" + fmt.Sprint(pipeline_id, "/", job_id))
}

func (a *cicdApi) GetPipelineBody(r *ghttp.Request) {
	var pipeline_id int = r.GetInt("pipeline_id")
	if !service.CheckAuthor(r.Context(), pipeline_id) {
		r.Response.WriteStatus(http.StatusForbidden)
	}
	pipeline_body := service.Pipeline.GetPipelineBodyString(pipeline_id)
	r.Response.WriteExit(pipeline_body)
}

func (a *cicdApi) GetPipelinePkgs(r *ghttp.Request) {
	var pipeline_id int = r.GetInt("pipeline_id")
	if !service.CheckAuthor(r.Context(), pipeline_id) {
		r.Response.WriteStatus(http.StatusForbidden)
	}
	pkgs := service.Cicd.GetPipelinePkgs(pipeline_id)
	r.Response.WriteExit(pkgs)
}

func (a *cicdApi) EmailPrefix(email string) string {
	if email == "" {
		return email
	}
	subemail := strings.Split(email, "@")
	email_prefix := subemail[0]
	return email_prefix
}

func (a *cicdApi) NewJob(r *ghttp.Request) {
	var pipeline_id int = r.GetInt("pipeline_id")
	if !service.CheckAuthor(r.Context(), pipeline_id) {
		r.Response.RedirectTo(UrlPrefix + "/forbidden")
	}
	var email string = service.Session.GetUser(r.Context()).Email
	username := a.EmailPrefix(email)
	envs := r.GetFormMap()
	job_id := service.Cicd.New(pipeline_id, envs, username)
	r.Response.RedirectTo(UrlPrefix + "/" + fmt.Sprint(pipeline_id) + "/" + strconv.FormatInt(job_id, 10))
}
