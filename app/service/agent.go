package service

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/tangr/cicdgo/app/dao"
	"github.com/tangr/cicdgo/app/model"
)

var Agent = agentService{}

type agentService struct{}

func (s *agentService) ListAgents() []model.ListAgents {
	agents := ([]model.ListAgents)(nil)
	err := dao.CicdAgent.Fields("id,agent_name,updated_at").Structs(&agents)
	if err != nil {
		glog.Error(err)
	}
	return agents
}

func (s *agentService) GetAgentNames() gdb.Result {
	result, err := dao.CicdAgent.Fields("id,agent_name").All()
	if err != nil {
		glog.Error(err)
	}
	return result
}

func (s *agentService) GetAgentInfo(agent_id string) (string, string) {
	type AgentInfo struct {
		AgentName   string `json:"agent_name"`
		AgentIpaddr string `json:"ipaddr"`
	}
	var newAgentInfo = &AgentInfo{}
	err := dao.CicdAgent.Fields("agent_name,ipaddr").Where("id=", agent_id).Struct(newAgentInfo)
	if err != nil {
		glog.Error(err)
	}
	agent_name := newAgentInfo.AgentName
	agent_ipaddr := newAgentInfo.AgentIpaddr
	return agent_name, agent_ipaddr
}

func (s *agentService) New(agent_name string, ipaddr string) int {
	len_orig_agent_name := len(agent_name)
	var hash_code string
	if len_orig_agent_name <= 50 {
		hash_code = Comm.RandSeq(60 - len_orig_agent_name)
	} else {
		hash_code = Comm.RandSeq(10)
	}
	agent_name = agent_name + ":" + hash_code
	new_agent := g.Map{"agent_name": agent_name, "ipaddr": ipaddr, "updated_at": gtime.Now().Timestamp()}
	result, err := dao.CicdAgent.Data(new_agent).Save()
	if err != nil {
		glog.Error(err)
	}
	agent_id, err := result.LastInsertId()
	if err != nil {
		glog.Error(err)
	}
	return int(agent_id)
}

func (s *agentService) Show(agent_id int) g.Map {
	result, err := dao.CicdAgent.Where("id=", agent_id).One()
	if err != nil {
		glog.Error(err)
	}
	return result.Map()
}

func (s *agentService) Update(agent_id int, agent_name string, ipaddr string) error {
	new_agent := g.Map{"agent_name": agent_name, "ipaddr": ipaddr, "updated_at": gtime.Now().Timestamp()}
	_, err := dao.CicdAgent.Data(new_agent).Where(g.Map{"id": agent_id}).Update()
	if err != nil {
		glog.Error(err)
	}
	return nil
}
