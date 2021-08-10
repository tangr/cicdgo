package service

import (
	"container/list"
	"fmt"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/tangr/cicdgo/app/dao"
	"github.com/tangr/cicdgo/app/model"
)

var WsServer = wsServer{}

type wsServer struct{}

var CiAgents map[int]string = make(map[int]string)
var CiAgentJobs map[int]list.List = make(map[int]list.List)

var CdAgents map[int]string = make(map[int]string)
var CdAgentJobs map[int]list.List = make(map[int]list.List)

func (s *wsServer) DoAgentCi(agentCiJobs *model.WsAgentSend, clientip string) *model.WsServerSend {
	var jobCiData model.WsServerSendMap
	var jobCiDatas model.WsServerSend
	if len(*agentCiJobs) == 0 {
		return &jobCiDatas
	}
	for _, ciJob := range *agentCiJobs {
		jobCiData = *s.HandleCIJob(&ciJob, clientip)
		jobCiDatas = append(jobCiDatas, jobCiData)
	}
	return &jobCiDatas
}

func (s *wsServer) DoAgentCd(agentCdJobs *model.WsAgentSend, clientip string) *model.WsServerSend {
	var jobCdData model.WsServerSendMap
	var jobCdDatas model.WsServerSend
	if len(*agentCdJobs) == 0 {
		return &jobCdDatas
	}
	for _, cdJob := range *agentCdJobs {
		jobCdData = *s.HandleCDJob(&cdJob, clientip)
		jobCdDatas = append(jobCdDatas, jobCdData)
	}
	return &jobCdDatas
}

func (s *wsServer) HandleCIJob(ciJob *model.WsAgentSendMap, clientip string) *model.WsServerSendMap {
	var jobCiData model.WsServerSendMap
	agentId := ciJob.AgentId
	agentName := ciJob.AgentName
	if !s.CheckAgentCI(agentId, agentName) {
		jobCiData.AgentId = agentId
		jobCiData.AgentName = agentName
		jobCiData.ErrMsg = "agentId: " + fmt.Sprint(agentId) + " and agentName: " + agentName + " not match."
		return &jobCiData
	}

	jobCiData.AgentId = ciJob.AgentId
	jobCiData.AgentName = ciJob.AgentName
	jobId := ciJob.JobId
	jobStatus := ciJob.JobStatus
	jobCiData.JobID = jobId
	jobCiData.JobStatus = jobStatus

	if jobStatus == "success" || jobStatus == "failed" {
		if _, err := dao.CicdJob.Data(g.Map{"job_status": jobStatus}).Where("id", jobId).Update(); err != nil {
			glog.Error(err)
		}
		if _, err := dao.CicdPackage.Data(g.Map{"job_status": jobStatus}).Where("job_id", jobId).Update(); err != nil {
			glog.Error(err)
		}
		jobOutput := ciJob.JobOutput
		if jobOutput != "" {
			if _, err := dao.CicdJob.Data(g.Map{"job_status": jobStatus}).Where("id", jobId).Update(); err != nil {
				glog.Error(err)
			}
			if _, err := dao.CicdLog.Data(g.Map{"job_id": jobId, "ipaddr": clientip, "job_status": jobStatus, "output": jobOutput, "updated_at": gtime.Now().Timestamp()}).Save(); err != nil {
				glog.Error(err)
			}
		}
		return s.GetCIJob(ciJob.AgentId, clientip)
	}
	if jobStatus == "running" {
		jobOutput := ciJob.JobOutput
		if _, err := dao.CicdJob.Data(g.Map{"job_status": jobStatus}).Where("id", jobId).Update(); err != nil {
			glog.Error(err)
		}
		if _, err := dao.CicdLog.Data(g.Map{"job_id": jobId, "ipaddr": clientip, "job_status": jobStatus, "output": jobOutput, "updated_at": gtime.Now().Timestamp()}).Save(); err != nil {
			glog.Error(err)
		}
		return &jobCiData
	}
	return s.GetCIJob(ciJob.AgentId, clientip)
}

func (s *wsServer) HandleCDJob(cdJob *model.WsAgentSendMap, clientip string) *model.WsServerSendMap {
	var jobCdData model.WsServerSendMap
	pipelineId := cdJob.AgentId
	pipelineName := cdJob.AgentName
	if !s.CheckAgentCD(pipelineId, pipelineName) {
		jobCdData.AgentId = pipelineId
		jobCdData.AgentName = pipelineName
		jobCdData.ErrMsg = "pipelineId: " + fmt.Sprint(pipelineId) + " and pipelineName: " + pipelineName + " not match."
		return &jobCdData
	}

	jobCdData.AgentId = cdJob.AgentId
	jobCdData.AgentName = cdJob.AgentName
	jobId := cdJob.JobId
	jobStatus := cdJob.JobStatus
	jobCdData.JobID = jobId
	jobCdData.JobStatus = jobStatus

	if jobStatus == "success" || jobStatus == "failed" {
		if _, err := dao.CicdJob.Data(g.Map{"job_status": jobStatus}).Where("id", jobId).Update(); err != nil {
			glog.Error(err)
		}
		jobOutput := cdJob.JobOutput
		if jobOutput != "" {
			if _, err := dao.CicdJob.Data(g.Map{"job_status": jobStatus}).Where("id", jobId).Update(); err != nil {
				glog.Error(err)
			}
			if _, err := dao.CicdLog.Data(g.Map{"job_id": jobId, "ipaddr": clientip, "job_status": jobStatus, "output": jobOutput, "updated_at": gtime.Now().Timestamp()}).Save(); err != nil {
				glog.Error(err)
			}
		}
		return s.GetCDJob(cdJob.AgentId, clientip)
	}
	if jobStatus == "running" {
		jobOutput := cdJob.JobOutput
		if _, err := dao.CicdJob.Data(g.Map{"job_status": jobStatus}).Where("id", jobId).Update(); err != nil {
			glog.Error(err)
		}
		if _, err := dao.CicdLog.Data(g.Map{"job_id": jobId, "ipaddr": clientip, "job_status": jobStatus, "output": jobOutput, "updated_at": gtime.Now().Timestamp()}).Save(); err != nil {
			glog.Error(err)
		}
		return &jobCdData
	}
	return s.GetCDJob(cdJob.AgentId, clientip)
}

func (s *wsServer) GetPipelineId(job_id int) int {
	pipeline_id, err := dao.CicdJob.Fields("pipeline_id").Where("id=", job_id).Value()
	if err != nil {
		glog.Error(err)
	}
	return pipeline_id.Int()
}

func (s *wsServer) GetCIJob(id int, clientip string) *model.WsServerSendMap {
	var newJobScriptP = new(model.JobScript)
	var newJobCiDataP = new(model.WsServerSendMap)
	if err := dao.CicdJob.Fields("id,job_status,script").Where("job_type", "BUILD").Where("agent_id", id).Where("job_status", "pending").Limit(1).Order("id asc").Struct(newJobScriptP); err != nil {
		glog.Debug(err)
	}
	newJobScript := *newJobScriptP
	jobId := newJobScript.ID
	jobStatus := newJobScript.JobStatus
	newJobCiDataP.AgentId = id
	newJobCiDataP.AgentName = CiAgents[id]
	newJobCiDataP.JobID = jobId
	newJobCiDataP.JobStatus = jobStatus
	newJobCiDataP.Body = newJobScript.Script.Body
	newJobCiDataP.Args = newJobScript.Script.Args
	newJobCiDataP.Envs = newJobScript.Script.Envs
	if len(newJobScript.Script.Envs) != 0 {
		package_name := fmt.Sprint(jobId) + "_" + newJobScript.Script.Envs["PKGRDM"]
		newJobCiDataP.Envs["PKGRDM"] = package_name
		newJobCiDataP.Envs["IPADDR"] = clientip
		pipeline_id := s.GetPipelineId(jobId)
		if _, err := dao.CicdPackage.Data(g.Map{"pipeline_id": pipeline_id, "job_id": jobId, "job_status": jobStatus, "package_name": package_name, "created_at": gtime.Now().Timestamp()}).Save(); err != nil {
			glog.Error(err)
		}
	}
	return newJobCiDataP
}

func (s *wsServer) GetCDJob(id int, clientip string) *model.WsServerSendMap {
	var newJobScriptP = new(model.JobScript)
	var newJobCdDataP = new(model.WsServerSendMap)
	if err := dao.CicdJob.Fields("id,job_status,script").Where("job_type", "DEPLOY").Where("pipeline_id", id).Where("job_status", "pending").Limit(1).Order("id desc").Struct(newJobScriptP); err != nil {
		glog.Debug(err)
	}
	newJobScript := *newJobScriptP
	newJobCdDataP.AgentId = id
	newJobCdDataP.AgentName = CdAgents[id]
	newJobCdDataP.JobID = newJobScript.ID
	newJobCdDataP.JobStatus = newJobScript.JobStatus
	newJobCdDataP.Body = newJobScript.Script.Body
	newJobCdDataP.Args = newJobScript.Script.Args
	newJobCdDataP.Envs = newJobScript.Script.Envs
	return newJobCdDataP
}

func (s *wsServer) CheckAgentCI(id int, agentname string) bool {
	if name, ok := CiAgents[id]; ok {
		return name == agentname
	}
	if i, err := dao.CicdAgent.Where("id", id).Where("agent_name", agentname).Count(); err != nil {
		glog.Error(err)
		return false
	} else {
		if i != 0 {
			CiAgents[id] = agentname
			return true
		}
		glog.Error(false)
		return false
	}
}

func (s *wsServer) CheckAgentCD(id int, pipelinename string) bool {
	if name, ok := CdAgents[id]; ok {
		return name == pipelinename
	}
	if i, err := dao.CicdPipeline.Where("id", id).Where("pipeline_name", pipelinename).Count(); err != nil {
		glog.Error(err)
		return false
	} else {
		if i != 0 {
			CdAgents[id] = pipelinename
			return true
		}
		return false
	}
}
