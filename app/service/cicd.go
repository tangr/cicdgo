package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/tangr/cicdgo/app/dao"
	"github.com/tangr/cicdgo/app/model"
)

var Cicd = cicdService{}

type cicdService struct{}

func (s *cicdService) ListCicd(group_ids []string) []model.ListPipelines {
	pipelines := ([]model.ListPipelines)(nil)
	err := dao.CicdPipeline.Fields("id,pipeline_name").WhereIn("group_id", group_ids).Structs(&pipelines)
	if err != nil {
		glog.Error(err)
	}
	return pipelines
}

func (s *cicdService) GetPkgJobInfo(job_id int) string {
	type JobInfo struct {
		Comment string `json:"comment"`
		Author  string `json:"author"`
	}
	jobInfo := &JobInfo{}
	err := dao.CicdJob.Fields("comment,author").Where("id=", job_id).Struct(jobInfo)
	if err != nil {
		glog.Error(err)
	}
	comment := jobInfo.Comment
	author := jobInfo.Author
	retjobinfo := comment + " by " + author
	return retjobinfo
}

func timeDiffNow(timestamp_int int) string {
	timestamp := int64(timestamp_int)
	timeNow := time.Now().Unix()
	timediff := timeNow - timestamp
	if timediff < 60 {
		return fmt.Sprint(timediff) + " secs ago"
	} else if timediff < 3600 {
		return fmt.Sprint(timediff/60) + " mins ago"
	} else if timediff < 259200 {
		return fmt.Sprint(timediff/3600) + " hours ago"
	} else if timediff > 0 {
		return fmt.Sprint(timediff/86400) + " days ago"
	}
	return fmt.Sprint(timestamp)
}

func (s *cicdService) GetPipelinePkgs(pipeline_id int) string {
	pipeline_pkgs, err := dao.CicdPackage.Where("pipeline_id=", pipeline_id).OrderDesc("job_id").Limit(30).All()
	if err != nil {
		glog.Error(err)
	}
	type Pkg struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	var pkgs []Pkg = make([]Pkg, 0)
	for _, pkg := range pipeline_pkgs {
		job_id := pkg["job_id"].Int()
		job_info := s.GetPkgJobInfo(job_id)
		created_at := pkg["created_at"].Int()
		timediff := timeDiffNow(created_at)
		// glog.Debug(fmt.Sprint(job_id) + " " + job_info + " at " + timediff)
		pkgName := fmt.Sprint(job_id) + " " + job_info + " at " + timediff
		pkgValue := pkg["package_name"].String()
		newpkg := Pkg{Name: pkgName, Value: pkgValue}
		pkgs = append(pkgs, newpkg)
	}
	pipeline_pkgs_gmap := g.Map{
		"success": true,
		"results": pkgs,
	}
	ret_pipeline_pkgs, _ := json.Marshal(pipeline_pkgs_gmap)
	return string(ret_pipeline_pkgs)
}

func (s *cicdService) GetJobType(job_id int) string {
	job_type, err := dao.CicdJob.Fields("job_type").Where("id=", job_id).Value()
	if err != nil {
		glog.Error(err)
	}
	return job_type.String()
}

func (s *cicdService) GetJobs(pipeline_id int, pageIndex int, pageSize int) ([]model.ListJobs, int) {
	jobs := ([]model.ListJobs)(nil)
	offSet := pageSize * (pageIndex - 1)
	err := dao.CicdJob.Fields("id,pipeline_id,agent_id,job_type,job_status,comment,author,created_at").Order("id desc").Where("pipeline_id=", pipeline_id).Limit(offSet, pageSize).Structs(&jobs)
	if err != nil {
		glog.Error(err)
	}
	totalSize, err := dao.CicdJob.Fields("id").Where("pipeline_id=", pipeline_id).Count()
	if err != nil {
		glog.Error(err)
	}
	return jobs, totalSize
}

func (s *cicdService) New(pipeline_id int, envs map[string]interface{}, username string) int64 {
	var script_name, script_args string
	var jobtype string
	var job_envs map[string]string = Cicd.ParseEnvs(envs)
	var comment string = job_envs["comment"]
	var job_type string = job_envs["job_type"]

	pipeline_name, agent_id, pipeline_body := Pipeline.GetPipelineBody(pipeline_id)
	if job_type == "BUILD" {
		jobtype = job_type
		script_name = pipeline_body.StageCI.Script
		script_args = pipeline_body.StageCI.Args
		job_envs["PKGRDM"] = Comm.RandSeq(20)
	} else if job_type == "DEPLOY" {
		jobtype = job_type
		script_name = pipeline_body.StageCI.Script
		script_args = pipeline_body.StageCI.Args
	} else {
		glog.Errorf("unsupported job_type: %s", job_type)
	}
	job_envs["PIPELINEID"] = fmt.Sprint(pipeline_id)
	job_envs["PIPELINENAME"] = strings.Split(pipeline_name, ":")[0]
	job_envs["USERNAME"] = username
	script_body := Script.GetScriptBody(script_name)
	new_jobscript := new(model.JobScriptValue)
	new_jobscript.Envs = job_envs
	new_jobscript.Args = script_args
	new_jobscript.Body = script_body

	new_job := g.Map{
		"pipeline_id": pipeline_id,
		"agent_id":    agent_id,
		"job_type":    jobtype,
		"job_status":  "pending",
		"script":      new_jobscript,
		"comment":     comment,
		"author":      username,
		"created_at":  gtime.Now().Timestamp(),
	}
	result, err := dao.CicdJob.Data(new_job).Save()
	if err != nil {
		glog.Error(err)
	}
	job_id, err := result.LastInsertId()
	if err != nil {
		glog.Error(err)
	}
	return job_id
}

func (s *cicdService) ParseEnvs(envs map[string]interface{}) map[string]string {
	var new_envs map[string]string = make(map[string]string)
	for k, v := range envs {
		switch v := v.(type) {
		case string:
			new_envs[k] = v
		case []interface{}:
			new_v := ""
			for _, u := range v {
				if new_v == "" {
					new_v = u.(string)
					continue
				}
				new_v = new_v + "," + u.(string)
			}
			new_envs[k] = new_v
		}
	}
	return new_envs
}

func (s *cicdService) GetJobEnvs(pipeline_id int, job_id int) string {
	jobScript := &model.JobScriptValue{}
	job_map := g.Map{"id": job_id, "pipeline_id": pipeline_id}
	script, err := dao.CicdJob.Fields("script").Where(job_map).Value()
	if err != nil {
		glog.Error(err)
	}
	script_byte := script.Bytes()
	err = json.Unmarshal(script_byte, jobScript)
	if err != nil {
		glog.Error(err)
	}
	jobEnvs := jobScript.Envs
	jobEnvs_byte, err := json.Marshal(jobEnvs)
	if err != nil {
		glog.Error(err)
	}
	jobEnvs_json := string(jobEnvs_byte)
	return jobEnvs_json
}

func (s *cicdService) CheckJobid(pipeline_id int, job_id int) bool {
	num, err := dao.CicdJob.Where(g.Map{"pipeline_id": pipeline_id, "id": job_id}).Count()
	if err != nil {
		glog.Error(err)
		return false
	}
	if num < 1 {
		return false
	}
	return true
}

func (s *cicdService) CheckTaskid(pipeline_id int, task_id int) bool {
	r, err := dao.CicdLog.Fields("job_id").Where("id = ?", task_id).Value()
	if err != nil {
		glog.Error(err)
		return false
	}
	job_id := r.Int()
	return s.CheckJobid(pipeline_id, job_id)
}

func (s *cicdService) AbortJob(pipeline_id int, task_id int) error {
	if !s.CheckTaskid(pipeline_id, task_id) {
		return nil
	}
	_, err := dao.CicdLog.Data("job_status='aborted'").Where("id", task_id).Update()
	if err != nil {
		glog.Error(err)
	}
	return nil
}

func (s *cicdService) RetryJob(pipeline_id int, task_id int) error {
	if !s.CheckTaskid(pipeline_id, task_id) {
		return nil
	}
	_, err := dao.CicdLog.Data("job_status='pending'").Where("id", task_id).Update()
	if err != nil {
		glog.Error(err)
	}
	return nil
}

func (s *cicdService) GetJobTasks(pipeline_id int, job_id int) []model.ListTasks {
	tasks := ([]model.ListTasks)(nil)
	if !s.CheckJobid(pipeline_id, job_id) {
		return tasks
	}
	err := dao.CicdLog.Fields("id,job_id,job_status,ipaddr,updated_at").Order("id desc").Where(g.Map{"job_id": job_id}).Structs(&tasks)
	if err != nil {
		glog.Error(err)
	}
	return tasks
}

func (s *cicdService) GetJobProgress(pipeline_id int, job_id int) (string, string, bool) {
	tasks := ([]model.ListTasks)(nil)
	if !s.CheckJobid(pipeline_id, job_id) {
		return "", "", true
	}
	job_map := g.Map{"job_id": job_id}
	err := dao.CicdLog.Fields("id,job_id,job_status,ipaddr,updated_at").Where(job_map).Structs(&tasks)
	if err != nil {
		glog.Error(err)
	}
	task_count := len(tasks)
	task_count_pending, task_count_running, task_count_success, task_count_failed := 0, 0, 0, 0
	for _, task := range tasks {
		if task.Job_status == "pending" {
			task_count_pending = task_count_pending + 1
		} else if task.Job_status == "running" {
			task_count_running = task_count_running + 1
		} else if task.Job_status == "success" {
			task_count_success = task_count_success + 1
		} else if task.Job_status == "failed" {
			task_count_failed = task_count_failed + 1
		}
	}
	var job_finished bool
	if task_count_success+task_count_failed == task_count {
		job_finished = true
	} else {
		job_finished = false
	}
	task_total := fmt.Sprint(task_count)
	task_value := fmt.Sprint(task_count_success, ",", task_count_failed, ",", task_count_running, ",", task_count_pending)
	return task_total, task_value, job_finished
}

func (s *cicdService) GetOutput(pipeline_id int, log_id int) *model.GetOutput {
	output := (*model.GetOutput)(nil)
	if !s.CheckTaskid(pipeline_id, log_id) {
		return output
	}
	err := dao.CicdLog.Fields("job_status,updated_at,output").Where(g.Map{"id": log_id}).Struct(&output)
	if err != nil {
		glog.Error(err)
	}
	return output
}
