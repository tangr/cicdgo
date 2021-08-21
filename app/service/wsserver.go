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

type AgentStatus map[string]int

type JobScript struct {
	Script model.JobScriptValue `json:"jobScript"`
}

type AgentStatusMapV struct {
	Updated int // Active Timestamp
	Status  string
	JobId   int
}

type AgentConcurrencyMapV struct {
	Number      int
	RunningList list.List
}

var CiAgentMapIdName map[int]string = make(map[int]string)
var CiAgentJobs map[int]list.List = make(map[int]list.List)

var CdAgentMapIdName map[int]string = make(map[int]string)

type AgentActivity map[string]*AgentStatusMapV
type AgentJobRunning map[string]string

// All Activity Build Agents
// var CiAgentMapActivity map[int]*AgentStatusMapV = make(map[int]*AgentStatusMapV)
var CiAgentMapActivity map[int]AgentActivity = make(map[int]AgentActivity)

// All Activity Pipeline Agents
var CdAgentMapPipelineActivity map[int]AgentActivity = make(map[int]AgentActivity)

// Current Running Pipeline Agents
var CdAgentMapPipelineRunning map[int]AgentJobRunning = make(map[int]AgentJobRunning)

// Current Running Pipeline Max Number
var CdAgentMapPipelineRunningConcurrent map[int]int = make(map[int]int)

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
	jobCiData.AgentId = agentId
	jobCiData.AgentName = agentName
	if !s.CheckAgentCI(agentId, agentName) {
		jobCiData.ErrMsg = "agentId: " + fmt.Sprint(agentId) + " and agentName: " + agentName + " not match."
		return &jobCiData
	}

	if CiAgentMapActivity[agentId] == nil {
		CiAgentMapActivity[agentId] = make(map[string]*AgentStatusMapV)
	}
	if CiAgentMapActivity[agentId][clientip] == nil {
		var AgentV *AgentStatusMapV = &AgentStatusMapV{Status: ""}
		CiAgentMapActivity[agentId][clientip] = AgentV
	}
	CiAgentMapActivity[agentId][clientip].Updated = int(gtime.Now().Timestamp())

	// if CiAgentMapActivity[agentId] == nil {
	// 	var AgentV *AgentStatusMapV = &AgentStatusMapV{Status: ""}
	// 	CiAgentMapActivity[agentId] = AgentV
	// }
	// CiAgentMapActivity[agentId].Updated = int(gtime.Now().Timestamp())

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
			// if _, err := dao.CicdJob.Data(g.Map{"job_status": jobStatus}).Where("id", jobId).Update(); err != nil {
			// 	glog.Error(err)
			// }
			if _, err := dao.CicdLog.Data(g.Map{"job_id": jobId, "ipaddr": clientip, "job_status": jobStatus, "output": jobOutput, "updated_at": gtime.Now().Timestamp()}).Save(); err != nil {
				glog.Error(err)
			}
		}
		// CiAgentMapActivity[agentId].Status = ""
		// CiAgentMapActivity[agentId].JobId = jobId
		CiAgentMapActivity[agentId][clientip].Status = ""
		CiAgentMapActivity[agentId][clientip].JobId = jobId
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
	jobCdData.AgentId = pipelineId
	jobCdData.AgentName = pipelineName
	if !s.CheckAgentCD(pipelineId, pipelineName) {
		jobCdData.ErrMsg = "pipelineId: " + fmt.Sprint(pipelineId) + " and pipelineName: " + pipelineName + " not match."
		return &jobCdData
	}

	if CdAgentMapPipelineActivity[pipelineId] == nil {
		CdAgentMapPipelineActivity[pipelineId] = make(map[string]*AgentStatusMapV)
	}
	if CdAgentMapPipelineActivity[pipelineId][clientip] == nil {
		var AgentV *AgentStatusMapV = &AgentStatusMapV{Status: ""}
		CdAgentMapPipelineActivity[pipelineId][clientip] = AgentV
	}
	CdAgentMapPipelineActivity[pipelineId][clientip].Updated = int(gtime.Now().Timestamp())

	jobId := cdJob.JobId
	jobStatus := cdJob.JobStatus
	jobCdData.JobID = jobId
	jobCdData.JobStatus = jobStatus

	if jobStatus == "success" || jobStatus == "failed" {
		// if _, err := dao.CicdJob.Data(g.Map{"job_status": jobStatus}).Where("id", jobId).Update(); err != nil {
		// 	glog.Error(err)
		// }
		jobOutput := cdJob.JobOutput
		if jobOutput != "" {
			// if _, err := dao.CicdJob.Data(g.Map{"job_status": jobStatus}).Where("id", jobId).Update(); err != nil {
			// 	glog.Error(err)
			// }
			if _, err := dao.CicdLog.Data(g.Map{"job_id": jobId, "ipaddr": clientip, "job_status": jobStatus, "output": jobOutput, "updated_at": gtime.Now().Timestamp()}).Save(); err != nil {
				glog.Error(err)
			}
		}
		CdAgentMapPipelineActivity[pipelineId][clientip].Status = ""
		CdAgentMapPipelineActivity[pipelineId][clientip].JobId = jobId
		return s.GetCDJob(cdJob.AgentId, clientip)
	}
	if jobStatus == "running" {
		jobOutput := cdJob.JobOutput
		// if _, err := dao.CicdJob.Data(g.Map{"job_status": jobStatus}).Where("id", jobId).Update(); err != nil {
		// 	glog.Error(err)
		// }
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

func (s *wsServer) GetCIJob(agent_id int, clientip string) *model.WsServerSendMap {
	var newJobScriptP = new(JobScript)
	var newJobCiDataP = new(model.WsServerSendMap)
	newJobCiDataP.AgentId = agent_id
	newJobCiDataP.AgentName = CiAgentMapIdName[agent_id]
	if CiAgentMapActivity[agent_id][clientip].Status != "pending" {
		return newJobCiDataP
	}
	jobId := CiAgentMapActivity[agent_id][clientip].JobId
	if err := dao.CicdJob.Fields("script").Where(g.Map{"id": jobId}).Struct(newJobScriptP); err != nil {
		glog.Debug(err)
	}
	newJobScript := *newJobScriptP
	jobStatus := "pending"
	newJobCiDataP.JobID = jobId
	newJobCiDataP.JobStatus = jobStatus
	newJobCiDataP.Body = newJobScript.Script.Body
	newJobCiDataP.Args = newJobScript.Script.Args
	newJobCiDataP.Envs = newJobScript.Script.Envs
	if len(newJobScript.Script.Envs) != 0 {
		package_name := fmt.Sprint(jobId) + "_" + newJobScript.Script.Envs["PKGRDM"]
		newJobCiDataP.Envs["PKGRDM"] = package_name
		newJobCiDataP.Envs["IPADDR"] = clientip
		newJobCiDataP.Envs["JOBID"] = fmt.Sprint(jobId)
		pipeline_id := s.GetPipelineId(jobId)
		if _, err := dao.CicdPackage.Data(g.Map{"pipeline_id": pipeline_id, "job_id": jobId, "job_status": jobStatus, "package_name": package_name, "created_at": gtime.Now().Timestamp()}).Save(); err != nil {
			glog.Error(err)
		}
	}
	return newJobCiDataP
}

func (s *wsServer) GetCDJob(pipeline_id int, clientip string) *model.WsServerSendMap {
	var newJobScriptP = new(JobScript)
	var newJobCdDataP = new(model.WsServerSendMap)
	newJobCdDataP.AgentId = pipeline_id
	newJobCdDataP.AgentName = CdAgentMapIdName[pipeline_id]
	if CdAgentMapPipelineActivity[pipeline_id][clientip].Status != "pending" {
		return newJobCdDataP
	}

	jobId := CdAgentMapPipelineActivity[pipeline_id][clientip].JobId
	deploy_job := g.Map{"id": jobId}
	if err := dao.CicdJob.Fields("script").Where(deploy_job).Struct(newJobScriptP); err != nil {
		glog.Debug(err)
	}
	newJobScript := *newJobScriptP
	newJobCdDataP.JobID = jobId
	newJobCdDataP.JobStatus = "pending"
	newJobCdDataP.Body = newJobScript.Script.Body
	newJobCdDataP.Args = newJobScript.Script.Args
	newJobCdDataP.Envs = newJobScript.Script.Envs
	if len(newJobCdDataP.Envs) != 0 {
		newJobCdDataP.Envs["IPADDR"] = clientip
		newJobCdDataP.Envs["JOBID"] = fmt.Sprint(jobId)
	}
	return newJobCdDataP
}

func (s *wsServer) SyncNewCIJob() {
	type NewJobBuild struct {
		ID      int `json:"jobid"`
		AgentId int `json:"agent_id"`
	}
	var newJobs = new([]NewJobBuild)

	if err := dao.CicdJob.Fields("id,agent_id").Where("job_type", "BUILD").WhereIn("job_status", g.Slice{"pending", "running"}).Structs(newJobs); err != nil {
		glog.Debug(err)
	}

	for _, newJob := range *newJobs {
		buidAgentId := newJob.AgentId
		jobId := newJob.ID

		job_status_v, err := dao.CicdLog.Fields("job_status").Where("job_id", jobId).Value()
		if err != nil {
			glog.Error(err)
		}
		job_status := job_status_v.String()
		if _, err := dao.CicdJob.Data(g.Map{"job_status": job_status}).Where("id", jobId).Update(); err != nil {
			glog.Error(err)
		}

		// fill up new running jobs
		// if k_clientip, ok := CiAgentMapActivity[buidAgentId]; ok {
		// 	if jobId > CiAgentMapActivity[buidAgentId][k_clientip].JobId {
		// 		CiAgentMapActivity[buidAgentId][k_clientip].Status = "pending"
		// 		CiAgentMapActivity[buidAgentId][k_clientip].JobId = jobId
		// 	}
		// }
		if AgentActivity, ok := CiAgentMapActivity[buidAgentId]; ok {
			for k_clientip := range AgentActivity {
				if jobId > CiAgentMapActivity[buidAgentId][k_clientip].JobId {
					CiAgentMapActivity[buidAgentId][k_clientip].Status = "pending"
					CiAgentMapActivity[buidAgentId][k_clientip].JobId = jobId
				}
			}
		}

	}
}

func (s *wsServer) SyncNewCDJob() {
	type NewJobDeploy struct {
		ID         int `json:"jobid"`
		PipelineId int `json:"pipelineid"`
	}
	var newJobs = new([]NewJobDeploy)

	if err := dao.CicdJob.Fields("id,pipeline_id").Where("job_type", "DEPLOY").WhereIn("job_status", g.Slice{"pending", "running"}).Structs(newJobs); err != nil {
		glog.Debug(err)
	}

	for _, newJob := range *newJobs {
		pipelineId := newJob.PipelineId
		jobId := newJob.ID

		finished_jobnum, err := dao.CicdLog.Where("job_id", jobId).WhereIn("job_status", g.Slice{"success", "failed"}).Count()
		if err != nil {
			glog.Error(err)
		}
		// glog.Debugf("finished_jobnum41", finished_jobnum)
		// glog.Debugf("finished_jobnum41", len(CdAgentMapPipelineActivity[pipelineId]))
		if finished_jobnum >= len(CdAgentMapPipelineActivity[pipelineId]) {
			if _, err := dao.CicdJob.Data(g.Map{"job_status": "success"}).Where("id", jobId).Update(); err != nil {
				glog.Error(err)
			}
		} else {
			if _, err := dao.CicdJob.Data(g.Map{"job_status": "running"}).Where("id", jobId).Update(); err != nil {
				glog.Error(err)
			}
		}

		if CdAgentMapPipelineRunning[pipelineId] == nil {
			CdAgentMapPipelineRunning[pipelineId] = make(map[string]string)
		}

		// clear run finished jobs
		if AgentJobRunning, ok := CdAgentMapPipelineRunning[pipelineId]; ok {
			for k_clientip := range AgentJobRunning {
				if CdAgentMapPipelineActivity[pipelineId][k_clientip].Status != "pending" {
					delete(CdAgentMapPipelineRunning[pipelineId], k_clientip)
				}
			}
		}

		// fill up new running jobs
		newJobRunningCapacity := 1 - len(CdAgentMapPipelineRunning[pipelineId])
		glog.Warning("newJobRunningCapacity: ", newJobRunningCapacity)
		glog.Warning("CdAgentMapPipelineRunning: ", CdAgentMapPipelineRunning[pipelineId])
		glog.Warning("CdAgentMapPipelineActivity: ", CdAgentMapPipelineActivity[pipelineId])
		if newJobRunningCapacity > 0 {
			for i := 0; i < newJobRunningCapacity; i++ {
				glog.Infof("i: %d", i)
				if AgentActivity, ok := CdAgentMapPipelineActivity[pipelineId]; ok {
					for k_clientip := range AgentActivity {
						if jobId > CdAgentMapPipelineActivity[pipelineId][k_clientip].JobId {
							CdAgentMapPipelineActivity[pipelineId][k_clientip].Status = "pending"
							CdAgentMapPipelineActivity[pipelineId][k_clientip].JobId = jobId
							CdAgentMapPipelineRunning[pipelineId][k_clientip] = ""
							break
						}
					}
				}
			}
		}

	}
}

func (s *wsServer) CheckAgentCI(agentid int, agentname string) bool {
	if name, ok := CiAgentMapIdName[agentid]; ok {
		return name == agentname
	}
	if i, err := dao.CicdAgent.Where("id", agentid).Where("agent_name", agentname).Count(); err != nil {
		glog.Error(err)
		return false
	} else {
		if i != 0 {
			CiAgentMapIdName[agentid] = agentname
			return true
		}
		glog.Error(false)
		return false
	}
}

func (s *wsServer) CheckAgentCD(pielineid int, pipelinename string) bool {
	if name, ok := CdAgentMapIdName[pielineid]; ok {
		return name == pipelinename
	}
	if i, err := dao.CicdPipeline.Where("id", pielineid).Where("pipeline_name", pipelinename).Count(); err != nil {
		glog.Error(err)
		return false
	} else {
		if i != 0 {
			CdAgentMapIdName[pielineid] = pipelinename
			return true
		}
		return false
	}
}

func (s *wsServer) GetAgentStatus(pipeline_id int, job_id int) map[string]int {
	type JobType struct {
		PipelineId int
		AgentId    int
		JobType    string
	}
	var job_type_struct = &JobType{}
	var newAgentStatus = make(map[string]int)

	job_map := g.Map{"pipeline_id": pipeline_id, "id": job_id}
	err := dao.CicdJob.Fields("pipeline_id,agent_id,job_type").Where(job_map).Struct(job_type_struct)
	if err != nil {
		glog.Error(err)
	}
	job_type := job_type_struct.JobType

	if job_type == "BUILD" {
		build_agent_id := job_type_struct.AgentId
		deploy_agents := CiAgentMapActivity[build_agent_id]
		for clientip, agent_map := range deploy_agents {
			mapk := fmt.Sprint(build_agent_id, "-", clientip)
			updated := agent_map.Updated
			newAgentStatus[mapk] = updated
			// return newAgentStatus
		}
		return newAgentStatus
	}
	deploy_agents := CdAgentMapPipelineActivity[pipeline_id]
	for clientip, pipeline_agent_map := range deploy_agents {
		mapk := fmt.Sprint(pipeline_id, "-", clientip)
		updated := pipeline_agent_map.Updated
		newAgentStatus[mapk] = updated
		// return newAgentStatus
	}
	return newAgentStatus
}
