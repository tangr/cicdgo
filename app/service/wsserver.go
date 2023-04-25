package service

import (
	"fmt"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/tangr/cicdgo/app/dao"
	"github.com/tangr/cicdgo/app/model"
)

var DeprecatedAfter int = g.Cfg().GetInt("server.console.DeprecatedAfter")

var WsServer = wsServer{}

type wsServer struct{}

// type AgentStatus map[string]int

type JobScript struct {
	Script JobScriptValue `json:"jobScript"`
}

type TaskStatus struct {
	TaskStatus string
	Output     string
}

// type AgentStatusMapV struct {
// 	Updated int // Active Timestamp
// 	Status  string
// 	JobId   int
// }

type CdAgentStatusMapV struct {
	Updated int // Active Timestamp
	Status  string
	JobId   int
}

// type CiAgentStatusMapV struct {
// 	Updated int // Active Timestamp
// 	Status  string
// 	JobId   []int
// }

// type AgentConcurrencyMapV struct {
// 	Number      int
// 	RunningList list.List
// }

type CiAgentNameClient struct {
	Name   string
	Ipaddr string
}

// var CiAgentMapIdName2 map[int]string = make(map[int]string)

// key: agentId
var CiAgentMapIdName map[int]*CiAgentNameClient = make(map[int]*CiAgentNameClient)

// var CiAgentJobs map[int]list.List = make(map[int]list.List)

var CdAgentMapIdName map[int]string = make(map[int]string)

// key: clientip
type CdAgentActivity map[string]*CdAgentStatusMapV

// type CiAgentActivity2 map[int]*CiAgentStatusMapV

type CiAgentActivity struct {
	Updated     int // Active Timestamp
	Status      string
	ClientIp    string
	RunningJobs map[int]string
}

type AgentJobRunning map[string]string

// All Activity Build Agents
// var CiAgentMapActivity map[int]*AgentStatusMapV = make(map[int]*AgentStatusMapV)
// var CiAgentMapActivity2 map[int]CiAgentActivity2 = make(map[int]CiAgentActivity2)

// key: agentId
var CiAgentMapActivity map[int]*CiAgentActivity = make(map[int]*CiAgentActivity)

// All Activity Pipeline Agents
// key: pipelineId
var CdAgentMapPipelineActivity map[int]CdAgentActivity = make(map[int]CdAgentActivity)

// Current Running Pipeline Agents
// key: pipelineId
var CdAgentMapPipelineRunning map[int]AgentJobRunning = make(map[int]AgentJobRunning)

// Current Running Pipeline Max Number
// var CdAgentMapPipelineRunningConcurrency map[int]int = make(map[int]int)

func (s *wsServer) DoAgentCi(agentCiJobs *model.WsAgentSend, clientip string) *model.WsServerSend {
	var jobCiData model.WsServerSendMap
	var jobCiDatas model.WsServerSend
	var jobCiDatasNew model.WsServerSend
	if len(*agentCiJobs) == 0 {
		return &jobCiDatas
	}
	agentCiJobsP := *agentCiJobs
	agentId := agentCiJobsP[0].AgentId
	agentName := agentCiJobsP[0].AgentName
	// var jobCiData model.WsServerSendMap
	jobCiData.AgentId = agentId
	jobCiData.AgentName = agentName
	if !s.CheckAgentCI(agentId, agentName, clientip) {
		jobCiData.ErrMsg = "agentId: " + fmt.Sprint(agentId) + " and agentName: " + agentName + " or ipaddr not match."
		jobCiDatas = append(jobCiDatas, jobCiData)
		return &jobCiDatas
	}
	for _, ciJob := range *agentCiJobs {
		jobCiData = *s.HandleCIJob(&ciJob, clientip)
		jobCiDatas = append(jobCiDatas, jobCiData)
	}
	jobCiDataNew := *s.GetCIJob(agentId, clientip)
	// g.Log().Error(jobCiDataNew)
	jobCiDatas = append(jobCiDatas, jobCiDataNew)
	// var CiJobsMap map[int]int = make(map[int]int)
	// g.Log().Error(jobCiDatas)
	for _, ciJob := range jobCiDatas {
		jobId := ciJob.JobId
		if jobId == 0 {
			continue
		}
		jobStatus := ciJob.JobStatus
		if jobStatus == "success" || jobStatus == "failed" {
			continue
		}
		// if CiJobsMap[jobId] != 0 {
		// 	continue
		// }
		scriptBody := ciJob.Body
		if jobStatus == "pending" && scriptBody == "" {
			continue
		}
		// CiJobsMap[jobId] = jobId
		jobCiDatasNew = append(jobCiDatasNew, ciJob)
	}
	// g.Log().Error(jobCiDatasNew)

	if len(jobCiDatasNew) == 0 {
		// g.Log().Errorf("jobCiDatasjobCiDatas: %+v", jobCiDatas)
		var jobCiData model.WsServerSendMap
		var jobCiDatas model.WsServerSend
		jobCiData.AgentId = agentId
		jobCiData.AgentName = agentName
		jobCiDatas = append(jobCiDatas, jobCiData)
		return &jobCiDatas
	}
	g.Log().Error(jobCiDatasNew)
	return &jobCiDatasNew
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

	jobId := ciJob.JobId
	jobStatus := ciJob.JobStatus
	jobCiData.JobId = jobId
	jobCiData.JobStatus = jobStatus

	if CiAgentMapActivity[agentId] == nil {
		CiAgentMapActivity[agentId] = &CiAgentActivity{Status: ""}
	}

	CiAgentMapActivity[agentId].Updated = int(gtime.Now().Timestamp())
	CiAgentMapActivity[agentId].ClientIp = clientip

	// g.Log().Error(ciJob)

	if jobStatus == "success" || jobStatus == "failed" {
		if _, err := dao.CicdJob.Data(g.Map{"job_status": jobStatus}).Where("id", jobId).Update(); err != nil {
			g.Log().Error(err)
		}
		if _, err := dao.CicdPackage.Data(g.Map{"job_status": jobStatus}).Where("job_id", jobId).Update(); err != nil {
			g.Log().Error(err)
		}
		jobOutput := ciJob.JobOutput
		if jobOutput != "" {
			newlog := g.Map{"agent_id": agentId, "job_type": "BUILD", "job_id": jobId, "ipaddr": clientip, "task_status": jobStatus, "output": jobOutput, "updated_at": gtime.Now().Timestamp()}
			if _, err := dao.CicdLog.Data(newlog).Save(); err != nil {
				g.Log().Error(err)
			}
		} else {
			newlog := g.Map{"agent_id": agentId, "job_type": "BUILD", "job_id": jobId, "ipaddr": clientip, "task_status": jobStatus, "updated_at": gtime.Now().Timestamp()}
			if _, err := dao.CicdLog.Data(newlog).Save(); err != nil {
				g.Log().Error(err)
			}
		}
		CiAgentMapActivity[agentId].RunningJobs[jobId] = jobStatus
		// CiAgentMapActivity[agentId][jobId].Status = ""
		// CiAgentMapActivity[agentId][clientip].JobId = jobId
		// return s.GetCIJob(agentId, clientip)
	}
	if jobStatus == "running" {
		if _, err := dao.CicdJob.Data(g.Map{"job_status": jobStatus}).Where("id", jobId).Update(); err != nil {
			g.Log().Error(err)
		}

		var lastTaskStatus = &TaskStatus{}
		if err := dao.CicdLog.Fields("task_status,output").Where(g.Map{"job_id": jobId, "ipaddr": clientip}).Struct(lastTaskStatus); err != nil {
			g.Log().Debug(err)
		}
		// latestTaskStatus := lastTaskStatus.TaskStatus
		latestTaskStatus := CiAgentMapActivity[agentId].RunningJobs[jobId]
		if latestTaskStatus == "aborted" {
			jobCiData.JobStatus = "aborted"
		}
		lastOutput := lastTaskStatus.Output
		jobOutput := ciJob.JobOutput
		if jobOutput != lastOutput {
			newlog := g.Map{"agent_id": agentId, "job_type": "BUILD", "job_id": jobId, "ipaddr": clientip, "task_status": jobStatus, "output": jobOutput, "updated_at": gtime.Now().Timestamp()}
			if _, err := dao.CicdLog.Data(newlog).Save(); err != nil {
				g.Log().Error(err)
			}
		}
		CiAgentMapActivity[agentId].RunningJobs[jobId] = jobStatus
		// return &jobCiData
	}
	return &jobCiData
	// return s.GetCIJob(agentId, clientip)
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
		CdAgentMapPipelineActivity[pipelineId] = make(map[string]*CdAgentStatusMapV)
	}
	if CdAgentMapPipelineActivity[pipelineId][clientip] == nil {
		var AgentV *CdAgentStatusMapV = &CdAgentStatusMapV{Status: "init"}
		CdAgentMapPipelineActivity[pipelineId][clientip] = AgentV
	}
	CdAgentMapPipelineActivity[pipelineId][clientip].Updated = int(gtime.Now().Timestamp())

	jobId := cdJob.JobId
	jobStatus := cdJob.JobStatus
	jobCdData.JobId = jobId
	jobCdData.JobStatus = jobStatus

	if jobStatus == "success" || jobStatus == "failed" {
		jobOutput := cdJob.JobOutput
		if jobOutput != "" {
			newlog := g.Map{"pipeline_id": pipelineId, "job_type": "DEPLOY", "job_id": jobId, "ipaddr": clientip, "task_status": jobStatus, "output": jobOutput, "updated_at": gtime.Now().Timestamp()}
			if _, err := dao.CicdLog.Data(newlog).Save(); err != nil {
				g.Log().Error(err)
			}
		}
		if jobStatus == "failed" {
			if _, err := dao.CicdJob.Data(g.Map{"job_status": jobStatus}).Where("id", jobId).Update(); err != nil {
				g.Log().Error(err)
			}
		}
		CdAgentMapPipelineActivity[pipelineId][clientip].Status = jobStatus
		CdAgentMapPipelineActivity[pipelineId][clientip].JobId = jobId
		return s.GetCDJob(pipelineId, clientip)
	}
	if jobStatus == "running" {
		var lastTaskStatus = &TaskStatus{}
		if err := dao.CicdLog.Fields("task_status,output").Where(g.Map{"job_id": jobId, "ipaddr": clientip}).Struct(lastTaskStatus); err != nil {
			g.Log().Debug(err)
		}
		// latestTaskStatus := lastTaskStatus.TaskStatus
		latestTaskStatus := CdAgentMapPipelineActivity[pipelineId][clientip].Status
		if latestTaskStatus == "aborted" {
			jobCdData.JobStatus = "aborted"
		}
		lastOutput := lastTaskStatus.Output
		jobOutput := cdJob.JobOutput
		if jobOutput != lastOutput {
			newlog := g.Map{"pipeline_id": pipelineId, "job_type": "DEPLOY", "job_id": jobId, "ipaddr": clientip, "task_status": jobStatus, "output": jobOutput, "updated_at": gtime.Now().Timestamp()}
			if _, err := dao.CicdLog.Data(newlog).Save(); err != nil {
				g.Log().Error(err)
			}
		}
		return &jobCdData
	}
	return s.GetCDJob(pipelineId, clientip)
}

func (s *wsServer) GetPipelineId(job_id int) int {
	pipeline_id, err := dao.CicdJob.Fields("pipeline_id").Where("id=", job_id).Value()
	if err != nil {
		g.Log().Error(err)
	}
	return pipeline_id.Int()
}

func (s *wsServer) GetCIJob(agentId int, clientip string) *model.WsServerSendMap {
	var newJobScriptP = new(JobScript)
	var newJobCiDataP = new(model.WsServerSendMap)
	newJobCiDataP.AgentId = agentId
	newJobCiDataP.AgentName = CiAgentMapIdName[agentId].Name
	if len(CiAgentMapActivity[agentId].RunningJobs) == 0 {
		return newJobCiDataP
	}

	g.Log().Error(CiAgentMapActivity[agentId].RunningJobs)

	var jobId int
	for newJobId, status := range CiAgentMapActivity[agentId].RunningJobs {
		if status == "pending" {
			jobId = newJobId
			break
		}
	}

	if jobId == 0 {
		return newJobCiDataP
	}

	if err := dao.CicdJob.Fields("script").Where(g.Map{"id": jobId}).Struct(newJobScriptP); err != nil {
		g.Log().Debug(err)
	}
	newJobScript := *newJobScriptP
	jobStatus := "pending"
	newJobCiDataP.JobId = jobId
	newJobCiDataP.JobStatus = jobStatus
	newJobCiDataP.Body = newJobScript.Script.Body
	newJobCiDataP.Args = newJobScript.Script.Args
	newJobCiDataP.Envs = newJobScript.Script.Envs
	g.Log().Error(newJobCiDataP)
	if len(newJobScript.Script.Envs) != 0 {
		package_name := fmt.Sprint(jobId) + "_" + newJobScript.Script.Envs["PKGRDM"]
		newJobCiDataP.Envs["PKGRDM"] = package_name
		newJobCiDataP.Envs["IPADDR"] = clientip
		newJobCiDataP.Envs["JOBID"] = fmt.Sprint(jobId)
		pipeline_id := s.GetPipelineId(jobId)
		if _, err := dao.CicdPackage.Data(g.Map{"pipeline_id": pipeline_id, "job_id": jobId, "job_status": jobStatus, "package_name": package_name, "created_at": gtime.Now().Timestamp()}).Save(); err != nil {
			g.Log().Error(err)
		}
	}
	return newJobCiDataP
}

func (s *wsServer) GetCDJob(pipelineId int, clientip string) *model.WsServerSendMap {
	var newJobScriptP = new(JobScript)
	var newJobCdDataP = new(model.WsServerSendMap)
	newJobCdDataP.AgentId = pipelineId
	newJobCdDataP.AgentName = CdAgentMapIdName[pipelineId]
	var jobId int
	if CdAgentMapPipelineActivity[pipelineId][clientip].Status == "init" {
		deploy_job := g.Map{"pipeline_id": pipelineId, "job_type": "DEPLOY", "job_status": "success"}
		job_id_var, err := dao.CicdJob.Fields("id").Where(deploy_job).OrderDesc("id").Limit(1).Value()
		if err != nil {
			g.Log().Debug(err)
		}
		jobId = job_id_var.Int()
		newJobCdDataP.JobStatus = "init"
	} else {
		if CdAgentMapPipelineActivity[pipelineId][clientip].Status != "pending" {
			return newJobCdDataP
		}
		jobId = CdAgentMapPipelineActivity[pipelineId][clientip].JobId
		newJobCdDataP.JobStatus = "pending"
	}

	if jobId == 0 {
		return newJobCdDataP
	}

	deploy_job := g.Map{"id": jobId}
	if err := dao.CicdJob.Fields("script").Where(deploy_job).Struct(newJobScriptP); err != nil {
		g.Log().Debug(err)
	}
	newJobScript := *newJobScriptP
	newJobCdDataP.JobId = jobId
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

	// if err := dao.CicdJob.Fields("id,agent_id").Where("job_type", "BUILD").WhereIn("job_status", g.Slice{"pending", "running"}).Structs(newJobs); err != nil {
	// g.Log().Error("SyncNewCIJob1")
	if err := dao.CicdJob.Fields("id,agent_id").Where("job_type", "BUILD").WhereIn("job_status", g.Slice{"pending"}).Structs(newJobs); err != nil {
		g.Log().Debug(err)
	}
	// g.Log().Error("SyncNewCIJob2")

	NowTimestamp := int(gtime.Now().Timestamp())

	for _, newJob := range *newJobs {
		agentId := newJob.AgentId
		jobId := newJob.ID

		job_status_v, err := dao.CicdLog.Fields("task_status").Where("job_id", jobId).Value()
		if err != nil {
			g.Log().Error(err)
		}
		job_status := job_status_v.String()
		if job_status != "" {
			if _, err := dao.CicdJob.Data(g.Map{"job_status": job_status}).Where("id", jobId).Update(); err != nil {
				g.Log().Error(err)
			}
		}

		if CiAgentActivity, ok := CiAgentMapActivity[agentId]; ok {
			// clear not activity agents
			if NowTimestamp-CiAgentMapActivity[agentId].Updated > DeprecatedAfter {
				delete(CiAgentMapActivity, agentId)
				continue
			}
			g.Log().Debugf("1CiAgentActivity.RunningJobs[jobId] %s, %d", CiAgentActivity.RunningJobs[jobId], jobId)

			if CiAgentMapActivity[agentId].RunningJobs == nil {
				CiAgentMapActivity[agentId].RunningJobs = make(map[int]string)
			}

			if _, ok := CiAgentActivity.RunningJobs[jobId]; ok {
				if CiAgentActivity.RunningJobs[jobId] != "pending" {
					g.Log().Error(CiAgentActivity.RunningJobs[jobId])
					delete(CiAgentMapActivity[agentId].RunningJobs, jobId)
					// continue
				}
			} else {
				CiAgentMapActivity[agentId].RunningJobs[jobId] = "pending"
				// break
			}

			g.Log().Debugf("2CiAgentActivity.RunningJobs[jobId] %s, %d", CiAgentActivity.RunningJobs[jobId], jobId)

			g.Log().Debugf("agentId %d, jobId: %d", agentId, jobId)
			g.Log().Debugf("CiAgentActivity %v, CiAgentActivity.RunningJobs: %v", CiAgentActivity, CiAgentActivity.RunningJobs)

		}

	}
}

func (s *wsServer) SyncNewCDJob() {
	type NewJobDeploy struct {
		ID          int `json:"jobid"`
		PipelineId  int `json:"pipelineid"`
		Concurrency int `json:"concurrency"`
	}
	var newJobs = new([]NewJobDeploy)

	if err := dao.CicdJob.Fields("id,pipeline_id,concurrency").Where("job_type", "DEPLOY").Where("job_status", "running").Structs(newJobs); err != nil {
		g.Log().Debug(err)
	}
	var JobConcurrency map[int]int = make(map[int]int)
	for _, newJob := range *newJobs {
		jobId := newJob.ID
		JobConcurrency[jobId] = newJob.Concurrency
	}

	type NewTaskDeploy struct {
		JobId      int `json:"jobid"`
		PipelineId int `json:"pipeline_id"`
	}
	var newTasks = new([]NewTaskDeploy)
	if err := dao.CicdLog.Fields("job_id,pipeline_id").Where("job_type", "DEPLOY").WhereIn("task_status", g.Slice{"pending", "running"}).Structs(newTasks); err != nil {
		g.Log().Debug(err)
	}

	for _, newTask := range *newTasks {
		var newJob = new(NewJobDeploy)
		newJob.ID = newTask.JobId
		newJob.PipelineId = newTask.PipelineId
		*newJobs = append(*newJobs, *newJob)
	}

	g.Log().Debugf("newJobs: %+v", newJobs)

	NowTimestamp := int(gtime.Now().Timestamp())

	for _, newJob := range *newJobs {
		pipelineId := newJob.PipelineId
		jobId := newJob.ID
		// concurrency := newJob.Concurrency
		concurrency := JobConcurrency[jobId]

		finished_jobnum, err := dao.CicdLog.Where("job_id", jobId).Where("task_status", "success").Count()
		if err != nil {
			g.Log().Error(err)
		}
		if finished_jobnum >= len(CdAgentMapPipelineActivity[pipelineId]) {
			if _, err := dao.CicdJob.Data(g.Map{"job_status": "success"}).Where("id", jobId).Update(); err != nil {
				g.Log().Error(err)
			}
		} else {
			if _, err := dao.CicdJob.Data(g.Map{"job_status": "running"}).Where("id", jobId).Update(); err != nil {
				g.Log().Error(err)
			}
		}

		if CdAgentMapPipelineRunning[pipelineId] == nil {
			CdAgentMapPipelineRunning[pipelineId] = make(map[string]string)
		}

		// clear run finished jobs
		if AgentJobRunning, ok := CdAgentMapPipelineRunning[pipelineId]; ok {
			for clientip := range AgentJobRunning {
				if CdAgentMapPipelineActivity[pipelineId][clientip] != nil {
					if CdAgentMapPipelineActivity[pipelineId][clientip].Status != "pending" {
						delete(CdAgentMapPipelineRunning[pipelineId], clientip)
					}
					// clear not activity agents
					if NowTimestamp-CdAgentMapPipelineActivity[pipelineId][clientip].Updated > DeprecatedAfter {
						delete(CdAgentMapPipelineActivity[pipelineId], clientip)
						continue
					}
				}
			}
		}

		// fill up new running jobs
		newJobRunningCapacity := concurrency - len(CdAgentMapPipelineRunning[pipelineId])
		if newJobRunningCapacity > 0 {
			for i := 0; i < newJobRunningCapacity; i++ {
				if CdAgentActivity, ok := CdAgentMapPipelineActivity[pipelineId]; ok {
					for clientip := range CdAgentActivity {
						if jobId > CdAgentMapPipelineActivity[pipelineId][clientip].JobId {
							CdAgentMapPipelineActivity[pipelineId][clientip].Status = "pending"
							CdAgentMapPipelineActivity[pipelineId][clientip].JobId = jobId
							CdAgentMapPipelineRunning[pipelineId][clientip] = ""
							break
						}
					}
				}
			}
		}

	}
}

func (s *wsServer) CheckAgentCI(agentid int, agentname string, clientip string) bool {
	// if name, ok := CiAgentMapIdName[agentid]; ok {
	// 	return name == agentname
	// }
	if CiAgentMapIdName[agentid] != nil {
		if CiAgentMapIdName[agentid].Name == agentname && CiAgentMapIdName[agentid].Ipaddr == clientip {
			return true
		}
	}
	if i, err := dao.CicdAgent.Where("id", agentid).Where(g.Map{"agent_name": agentname, "ipaddr": clientip}).Count(); err != nil {
		g.Log().Error(err)
		return false
	} else {
		if i != 0 {
			if CiAgentMapIdName[agentid] == nil {
				CiAgentMapIdName[agentid] = &CiAgentNameClient{}
			}
			CiAgentMapIdName[agentid].Name = agentname
			CiAgentMapIdName[agentid].Ipaddr = clientip
			return true
		}
		g.Log().Error(false)
		return false
	}
}

func (s *wsServer) CheckAgentCD(pielineid int, pipelinename string) bool {
	if name, ok := CdAgentMapIdName[pielineid]; ok {
		return name == pipelinename
	}
	if i, err := dao.CicdPipeline.Where("id", pielineid).Where("pipeline_name", pipelinename).Count(); err != nil {
		g.Log().Error(err)
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
		g.Log().Error(err)
	}
	job_type := job_type_struct.JobType

	if job_type == "BUILD" {
		build_agent_id := job_type_struct.AgentId
		if CiAgentMapActivity[build_agent_id] == nil {
			return newAgentStatus
		}
		deploy_agents := CiAgentMapActivity[build_agent_id]
		updated := deploy_agents.Updated
		clientip := deploy_agents.ClientIp
		mapk := fmt.Sprint(build_agent_id, "-", clientip)
		newAgentStatus[mapk] = updated
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

func (s *wsServer) RetryTask(task_id int, job_id int, clientip string) bool {
	type TaskInfo struct {
		PipelineId int    `json:"pipeline_id"`
		JoyType    string `json:"job_type"`
		JoyId      int    `json:"job_id"`
		TaskStatus string `json:"task_status"`
		Ipaddr     string `json:"ipaddr"`
	}
	lastTaskInfo := &TaskInfo{}
	err := dao.CicdLog.Fields("pipeline_id,job_type,job_id,task_status,ipaddr").Where(g.Map{"id": task_id}).Struct(&lastTaskInfo)
	if err != nil {
		g.Log().Error(err)
	}
	pipelineId := lastTaskInfo.PipelineId
	lastTaskStatus := lastTaskInfo.TaskStatus
	job_type := lastTaskInfo.JoyType
	lastJobId := lastTaskInfo.JoyId
	lastTaskClientip := lastTaskInfo.Ipaddr

	if lastTaskClientip != clientip {
		return false
	}
	if lastTaskStatus != "success" && lastTaskStatus != "failed" {
		return false
	}
	if lastJobId != job_id {
		return false
	}
	if job_type != "DEPLOY" {
		return false
	}
	if CdAgentMapPipelineActivity[pipelineId] == nil {
		return false
	}
	if _, ok := CdAgentMapPipelineActivity[pipelineId]; ok {
		CdAgentMapPipelineActivity[pipelineId][clientip].Status = "pending"
		return true
	}
	return false
}

func (s *wsServer) AbortTask(task_id int, job_id int, clientip string) bool {
	type TaskInfo struct {
		PipelineId int    `json:"pipeline_id"`
		AgentId    int    `json:"agent_id"`
		JoyType    string `json:"job_type"`
		JoyId      int    `json:"job_id"`
		TaskStatus string `json:"task_status"`
		Ipaddr     string `json:"ipaddr"`
	}
	lastTaskInfo := &TaskInfo{}
	err := dao.CicdLog.Fields("pipeline_id,agent_id,job_type,job_id,task_status,ipaddr").Where(g.Map{"id": task_id}).Struct(&lastTaskInfo)
	if err != nil {
		g.Log().Error(err)
	}
	jobId := job_id
	lastTaskStatus := lastTaskInfo.TaskStatus
	job_type := lastTaskInfo.JoyType
	lastJobId := lastTaskInfo.JoyId
	lastTaskClientip := lastTaskInfo.Ipaddr

	if lastTaskClientip != clientip {
		return false
	}
	if lastTaskStatus != "running" {
		return false
	}
	if lastJobId != jobId {
		return false
	}
	if job_type == "BUILD" {
		agentId := lastTaskInfo.AgentId
		if CiAgentMapActivity[agentId] == nil {
			return false
		}
		if _, ok := CiAgentMapActivity[agentId]; ok {
			if CiAgentMapActivity[agentId].RunningJobs[jobId] == "running" {
				CiAgentMapActivity[agentId].RunningJobs[jobId] = "aborted"
				return true
			}
			if _, err := dao.CicdLog.Data(g.Map{"task_status": "failed"}).Where("id", task_id).Update(); err != nil {
				g.Log().Error(err)
			}
			if _, err := dao.CicdJob.Data(g.Map{"job_status": "failed"}).Where("id", job_id).Update(); err != nil {
				g.Log().Error(err)
			}
			return true
		}
		return false
	}
	if job_type == "DEPLOY" {
		pipelineId := lastTaskInfo.PipelineId
		if CdAgentMapPipelineActivity[pipelineId] == nil {
			return false
		}
		if _, ok := CdAgentMapPipelineActivity[pipelineId]; ok {
			if CdAgentMapPipelineActivity[pipelineId][clientip].Status == "running" {
				CdAgentMapPipelineActivity[pipelineId][clientip].Status = "aborted"
				return true
			}
			if _, err := dao.CicdLog.Data(g.Map{"task_status": "failed"}).Where("id", task_id).Update(); err != nil {
				g.Log().Error(err)
			}
			if _, err := dao.CicdJob.Data(g.Map{"job_status": "failed"}).Where("id", job_id).Update(); err != nil {
				g.Log().Error(err)
			}
		}
		return false
	}
	return false
}
