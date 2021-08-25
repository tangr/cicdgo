package router

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"

	"github.com/tangr/cicdgo/app/api"
	"github.com/tangr/cicdgo/app/service"
)

func init() {
	var UrlPrefix string = g.Cfg().GetString("server.console.UrlPrefix")
	s := g.Server("console")
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/", api.Cicd.ListCicd)
		group.Group(UrlPrefix, func(group *ghttp.RouterGroup) {
			group.GET("/login", api.User.LoginPage)
			group.Middleware(
				service.Middleware.Ctx,
				service.Middleware.Authen,
			)
			group.GET("/forbidden", api.User.ForbiddenPage)
			group.GET("/", api.Cicd.ListCicd)
			group.GET("/:pipeline_id", api.Cicd.ShowCicd)
			group.GET("/:pipeline_id/:job_id", api.Cicd.ShowJob)

			group.GET("/logout", api.User.SignOut)

			group.Middleware(
				service.Middleware.AuthorAdmin,
			)
			group.GET("/scripts", api.Script.ListScripts)
			group.GET("/scripts/new", api.Script.NewScript)
			group.GET("/scripts/:id", api.Script.ShowScript)

			group.GET("/agents", api.Agent.ListAgents)
			group.GET("/agents/new", api.Agent.NewAgent)
			group.GET("/agents/:id", api.Agent.ShowAgent)

			group.GET("/pipelines", api.Pipeline.ListPipelines)
			group.GET("/pipelines/new", api.Pipeline.NewPipeline)
			group.GET("/pipelines/:id", api.Pipeline.ShowPipeline)

			group.GET("/users", api.User.ListUsers)
			group.GET("/users/new", api.User.NewUser)
			group.GET("/users/:id", api.User.ShowUser)

			group.GET("/groups", api.Group.ListGroups)
			group.GET("/groups/new", api.Group.NewGroup)
			group.GET("/groups/:id", api.Group.ShowGroup)
		})

		group.Group(UrlPrefix+"/v1", func(group *ghttp.RouterGroup) {
			group.Middleware(
				service.Middleware.Ctx,
			)
			group.POST("/users/login", api.User.SignIn)
			group.Middleware(
				service.Middleware.Authen,
			)
			group.GET("/:pipeline_id/:task_id/log", api.Cicd.GetLog)
			group.GET("/:pipeline_id/:job_id/env", api.Cicd.GetJobEnvs)
			group.GET("/:pipeline_id/:job_id/progress", api.Cicd.GetJobProgress)
			group.GET("/:pipeline_id/body", api.Cicd.GetPipelineBody)
			group.GET("/:pipeline_id/pkgs", api.Cicd.GetPipelinePkgs)
			group.POST("/:pipeline_id/:job_id/concurrency", api.Cicd.PostJobConcurrency)
			group.POST("/:pipeline_id/:job_id/jobstatus", api.Cicd.PostJobStatus)
			group.POST("/:pipeline_id/:task_id/taskstatus", api.Cicd.PostTaskStatus)
			group.POST("/:pipeline_id/:job_id/:task_id/taskstatus", api.Cicd.PostTaskStatus)

			group.POST("/:pipeline_id/newjob", api.Cicd.NewJob)
			// group.POST("/:pipeline_id/:task_id/abort", api.Cicd.AbortJob)
			// group.POST("/:pipeline_id/:task_id/retry", api.Cicd.RetryJob)

			group.Middleware(
				service.Middleware.AuthorAdmin,
			)
			group.POST("/scripts", api.Script.New)
			group.GET("/scripts/:id", api.Script.Show)
			group.POST("/scripts/:id", api.Script.Update)

			group.POST("/agents", api.Agent.New)
			group.GET("/agents/:id", api.Agent.Show)
			group.POST("/agents/:id", api.Agent.Update)

			group.POST("/pipelines", api.Pipeline.New)
			group.GET("/pipelines/:id", api.Pipeline.Show)
			group.POST("/pipelines/:id", api.Pipeline.Update)

			group.POST("/users", api.User.New)
			group.GET("/users/:id", api.User.Show)
			group.POST("/users/:id", api.User.Update)

			group.POST("/groups", api.Group.New)
			group.GET("/groups/:id", api.Group.Show)
			group.POST("/groups/:id", api.Group.Update)

		})
	})

}

func init() {
	s := g.Server("wscicd")
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Group("/wsv1", func(group *ghttp.RouterGroup) {
			group.ALL("/wsci", api.WsServer.Wsci)
			group.ALL("/wscd", api.WsServer.Wscd)
			group.GET("/:pipeline_id/:job_id/status", api.WsServer.GetAgentStatus)
			group.POST("/:pipeline_id/:job_id/concurrency", api.WsServer.GetAgentStatus)
		})
	})
}
