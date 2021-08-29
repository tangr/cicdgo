// ==========================================================================
// This is auto-generated by gf cli tool. Fill this file as you wish.
// ==========================================================================

package model

// Fill with you ideas below.

type ListAgents struct {
	Id         int    `json:"agent_id"`
	Agent_name string `json:"agent_name"`
	Updated_at int    `json:"updated_at"`
}

type ListGroups struct {
	Id         int    `json:"id"`
	Group_name string `json:"groupName"`
}

type ListJobs struct {
	Id          int    `json:"job_id"`
	Pipeline_id int    `json:"pipeline_id"`
	Agent_id    int    `json:"agent_id"`
	Job_type    string `json:"job_type"`
	Job_status  string `json:"job_status"`
	Comment     string `json:"comment"`
	Author      string `json:"author"`
	Created_at  int    `json:"created_at"`
}

type AgentStatusMapV struct {
	Updated int
	Status  string
}

type WsAgentSend []WsAgentSendMap

type WsAgentSendMap struct {
	AgentId   int    `json:"agentId"`
	AgentName string `json:"agentName"`
	JobId     int    `json:"jobId"`
	JobStatus string `json:"jobStatus"`
	JobOutput string `json:"jobOutput"`
}

// type WsAgentSendMap2 struct {
// 	IdId      int    `json:"idid"`
// 	Name      string `json:"name"`
// 	JobId     int    `json:"jobId"`
// 	JobStatus string `json:"jobStatus"`
// 	JobOutput string `json:"jobOutput"`
// }

type WsServerSend []WsServerSendMap

type WsServerSendMap struct {
	AgentId   int               `json:"agentId"`
	AgentName string            `json:"agentName"`
	JobId     int               `json:"jobId"`
	JobStatus string            `json:"jobStatus"`
	Body      string            `json:"scriptBody"`
	Envs      map[string]string `json:"scriptEnvs"`
	Args      string            `json:"scriptArgs"`
	ErrMsg    string            `json:"errmsg"`
}

// type WsAgentCiSend []WsAgentCiSendMap

// type WsAgentCiSendMap struct {
// 	AgentId   int    `json:"agentId" v:"required#AgentId不能为空"`
// 	AgentName string `json:"agentName" v:"required#AgentName不能为空"`
// 	JobId     int    `json:"jobId"`
// 	JobStatus string `json:"jobStatus"`
// 	JobOutput string `json:"jobOutput"`
// }

// type WsAgentCdSend []WsAgentCdSendMap

// type WsAgentCdSendMap struct {
// 	PipelineId   int    `json:"pipelineId" v:"required#pipelineId不能为空"`
// 	PipelineName string `json:"pipelineName" v:"required#pipelineName不能为空"`
// 	JobId        int    `json:"jobId"`
// 	JobStatus    string `json:"jobStatus"`
// 	JobOutput    string `json:"jobOutput"`
// }

// type WsAgentCiSend2 struct {
// 	AgentID   int       `json:"agentId" v:"required#AgentId不能为空"`
// 	JobID     int       `json:"jobId"`
// 	JobStatus string    `json:"jobStatus"`
// 	JobScript JobScript `json:"jobScript"`
// }

// type WsServerCiSend2 struct {
// 	ErrCode int                 `json:"errcode"`
// 	ErrMsg  string              `json:"errmsg"`
// 	Data    []WsServerCiSendMap `json:"data"`
// }

// type WsServerCiSend []WsServerCiSendMap

// type WsServerCiSendMap struct {
// 	AgentId   int               `json:"agentId"`
// 	AgentName string            `json:"agentName"`
// 	JobID     int               `json:"jobId"`
// 	JobStatus string            `json:"jobStatus"`
// 	Body      string            `json:"scriptBody"`
// 	Envs      map[string]string `json:"scriptEnvs"`
// 	Args      string            `json:"scriptArgs"`
// 	ErrMsg    string            `json:"errmsg"`
// }

// type WsServerCdSend []WsServerCdSendMap

// type WsServerCdSendMap struct {
// 	PipelineId   int               `json:"pipelineId"`
// 	PipelineName string            `json:"pipelineName"`
// 	JobID        int               `json:"jobId"`
// 	JobStatus    string            `json:"jobStatus"`
// 	Body         string            `json:"scriptBody"`
// 	Envs         map[string]string `json:"scriptEnvs"`
// 	Args         string            `json:"scriptArgs"`
// 	ErrMsg       string            `json:"errmsg"`
// }

type JobScriptValue struct {
	Body string            `json:"scriptBody"`
	Envs map[string]string `json:"scriptEnvs"`
	Args string            `json:"scriptArgs"`
}

// type AgentCIs struct {
// 	ID   int    `json:"id"`
// 	Name string `json:"name"`
// }

// type AgentCDs struct {
// 	ID   int    `json:"id"`
// 	Name string `json:"name"`
// }

type JobMeta struct {
	ID        int    `json:"jobid"`
	JobStatus string `json:"status"`
}

// type JobRunning struct {
// 	ID        int    `json:"jobid"`
// 	JobStatus string `json:"status"`
// 	Output    string `json:"output"`
// }

// type JobScript struct {
// 	ID        int            `json:"jobid"`
// 	JobStatus string         `json:"status"`
// 	Script    JobScriptValue `json:"jobScript"`
// }

// type JobScriptData struct {
// 	ID        int               `json:"jobid"`
// 	JobStatus string            `json:"status"`
// 	Body      string            `json:"scriptBody"`
// 	Envs      map[string]string `json:"scriptEnvs"`
// 	Args      string            `json:"scriptArgs"`
// }

// type ListTasks struct {
// 	Id         int    `json:"log_id"`
// 	Job_id     int    `json:"job_id"`
// 	Job_status string `json:"job_status"`
// 	Ipaddr     string `json:"ipaddr"`
// 	Updated_at int    `json:"updated_at"`
// }

// type GetOutput struct {
// 	Job_status string `json:"status"`
// 	Updated_at int    `json:"updated_at"`
// 	Output     string `json:"output"`
// }

type ListPipelines struct {
	Id            int    `json:"pipeline_id"`
	Pipeline_name string `json:"pipeline_name"`
}

type Script struct {
	Args   string `json:"script_args"`
	Script string `json:"script_name"`
}

type PipelineBody struct {
	StageCI Script `json:"stageCI"`
	StageCD Script `json:"stageCD"`
}

type ListScripts struct {
	Id          int    `json:"script_id"`
	Script_name string `json:"script_name"`
	Author      string `json:"author"`
	Updated_at  int    `json:"updated_at"`
}

type ListUsers struct {
	Id         int    `json:"id"`
	Email      string `json:"email"`
	Groups     string `json:"groups"`
	Updated_at int    `json:"updated_at"`
}

type UserInfo struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserGroupIds struct {
	GroupId string `json:"group_id"`
}

type GetUser struct {
	Email    string `json:"email"`
	Group_Id string `json:"groups"`
}

type UserApiSession struct {
	Id      int    `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool
}

type UserApiSignInReq struct {
	Email    string `v:"required#账号不能为空"`
	Password string `v:"required#密码不能为空"`
}
