package api

import (
	"encoding/json"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"

	"github.com/tangr/cicdgo/app/model"
	"github.com/tangr/cicdgo/app/service"
)

type wsServer struct{}

var WsServer = &wsServer{}

func (a *wsServer) Wsci(r *ghttp.Request) {
	var (
		recvMsg *model.WsAgentSend
		sendMsg *model.WsServerSend
	)

	go func() {
		// defer close(done)
		for {
			time.Sleep(time.Second)
			// g.Log().Warningf("in Wsci")
			service.WsServer.SyncNewCIJob()
		}
	}()

	ws, err := r.WebSocket()
	if err != nil {
		g.Log().Error(err)
		r.Exit()
	}

	for {
		err := ws.ReadJSON(&recvMsg)
		if recvMsg == nil {
			time.Sleep(time.Second)
			continue
		}
		g.Log().Debug("CI1************************")
		clientip := r.GetClientIp()
		recvJson, _ := json.Marshal(*recvMsg)
		g.Log().Debugf("recvJson: %s", recvJson)
		if err != nil {
			g.Log().Error(err)
			return
		}

		sendMsg = service.WsServer.DoAgentCi(recvMsg, clientip)

		if err = ws.WriteJSON(sendMsg); err != nil {
			g.Log().Error(err)
		}
		sendJson, _ := json.Marshal(sendMsg)
		g.Log().Debugf("sendJson: %s", string(sendJson))
	}
}

func (a *wsServer) Wscd(r *ghttp.Request) {
	var (
		recvMsg *model.WsAgentSend
		sendMsg *model.WsServerSend
	)

	// done := make(chan struct{})
	go func() {
		// defer close(done)
		for {
			time.Sleep(time.Second)
			// g.Log().Warningf("in Wscd")
			service.WsServer.SyncNewCDJob()
		}
	}()

	ws, err := r.WebSocket()
	if err != nil {
		g.Log().Error(err)
		r.Exit()
	}
	for {
		err := ws.ReadJSON(&recvMsg)
		if recvMsg == nil {
			g.Log().Error(recvMsg)
			time.Sleep(time.Second)
			continue
		}
		g.Log().Debug("CD1************************")
		// log.Printf("GetClientIp: %s", r.GetClientIp())
		clientip := r.GetClientIp()
		recvJson, _ := json.Marshal(*recvMsg)
		g.Log().Debugf("recvJson: %s", recvJson)
		if err != nil {
			g.Log().Error(err)
			return
		}

		sendMsg = service.WsServer.DoAgentCd(recvMsg, clientip)

		if err = ws.WriteJSON(sendMsg); err != nil {
			g.Log().Error(err)
		}
		sendJson, _ := json.Marshal(sendMsg)
		g.Log().Debugf("sendJson: %s", string(sendJson))
	}
}

func (a *wsServer) GetAgentStatus(r *ghttp.Request) {
	var pipeline_id int = r.GetInt("pipeline_id")
	var job_id int = r.GetInt("job_id")
	agent_status := service.WsServer.GetAgentStatus(pipeline_id, job_id)
	r.Response.WriteExit(agent_status)
}

func (a *wsServer) RetryTask(r *ghttp.Request) {
	var task_id int = r.GetInt("task_id")
	var job_id int = r.GetInt("job_id")
	var clientip string = r.GetString("clientip")
	result := service.WsServer.RetryTask(task_id, job_id, clientip)
	r.Response.WriteExit(result)
}

func (a *wsServer) AbortTask(r *ghttp.Request) {
	var task_id int = r.GetInt("task_id")
	var job_id int = r.GetInt("job_id")
	var clientip string = r.GetString("clientip")
	result := service.WsServer.AbortTask(task_id, job_id, clientip)
	r.Response.WriteExit(result)
}
