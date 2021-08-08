package api

import (
	"encoding/json"
	"time"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"

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

	ws, err := r.WebSocket()
	if err != nil {
		glog.Error(err)
		r.Exit()
	}

	for {
		err := ws.ReadJSON(&recvMsg)
		if recvMsg == nil {
			time.Sleep(time.Second)
			continue
		}
		glog.Debug("CI************************")
		clientip := r.GetClientIp()
		recvJson, _ := json.Marshal(*recvMsg)
		glog.Debugf("recvJson: %s", recvJson)
		if err != nil {
			glog.Error(err)
			return
		}

		sendMsg = service.WsServer.DoAgentCi(recvMsg, clientip)

		if err = ws.WriteJSON(sendMsg); err != nil {
			glog.Error(err)
		}
		sendJson, _ := json.Marshal(sendMsg)
		glog.Debugf("sendJson: %s", string(sendJson))
	}
}

func (a *wsServer) Wscd(r *ghttp.Request) {
	var (
		recvMsg *model.WsAgentSend
		sendMsg *model.WsServerSend
	)

	ws, err := r.WebSocket()
	if err != nil {
		glog.Error(err)
		r.Exit()
	}
	for {
		err := ws.ReadJSON(&recvMsg)
		if recvMsg == nil {
			glog.Error(recvMsg)
			time.Sleep(time.Second)
			continue
		}
		glog.Debug("CD************************")
		// log.Printf("GetClientIp: %s", r.GetClientIp())
		clientip := r.GetClientIp()
		recvJson, _ := json.Marshal(*recvMsg)
		glog.Debugf("recvJson: %s", recvJson)
		if err != nil {
			glog.Error(err)
			return
		}

		sendMsg = service.WsServer.DoAgentCd(recvMsg, clientip)

		if err = ws.WriteJSON(sendMsg); err != nil {
			glog.Error(err)
		}
		sendJson, _ := json.Marshal(sendMsg)
		glog.Debugf("sendJson: %s", string(sendJson))
	}
}
