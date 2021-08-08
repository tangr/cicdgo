package service

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/tangr/cicdgo/app/dao"
	"github.com/tangr/cicdgo/app/model"
)

var Pipeline = pipelineService{}

type pipelineService struct{}

func (s *pipelineService) ListPipelines() []model.ListPipelines {
	pipelines := ([]model.ListPipelines)(nil)
	err := dao.CicdPipeline.Fields("id,pipeline_name").Structs(&pipelines)
	if err != nil {
		glog.Error(err)
	}
	return pipelines
}

func (s *pipelineService) GetPipelineName(pipeline_id int) string {
	pipeline_name, err := dao.CicdPipeline.Fields("pipeline_name").Where("id=", pipeline_id).Value()
	if err != nil {
		glog.Error(err)
	}
	return pipeline_name.String()
}

func (s *pipelineService) GetGroupId(pipeline_id int) int {
	group_id, err := dao.CicdPipeline.Fields("group_id").Where("id=", pipeline_id).Value()
	if err != nil {
		glog.Error(err)
	}
	return group_id.Int()
}

func (s *pipelineService) GetPipelineBodyString(pipeline_id int) string {
	pipeline_body, err := dao.CicdPipeline.Fields("body").Where("id=", pipeline_id).Value()
	if err != nil {
		glog.Error(err)
	}
	return pipeline_body.String()
}

func (s *pipelineService) GetPipelineBody(pipeline_id int) (int, model.PipelineBody) {
	type Pipeline struct {
		Agent_id int                `json:"agent_int"`
		Body     model.PipelineBody `json:"pipeline_body"`
	}
	var pipeline Pipeline
	err := dao.CicdPipeline.Fields("agent_id, body").Where("id=", pipeline_id).Struct(&pipeline)
	if err != nil {
		glog.Error(err)
	}
	agent_id := pipeline.Agent_id
	pipeline_body := pipeline.Body
	return agent_id, pipeline_body
}

func (s *pipelineService) New(pipeline_name string, group_id int, agent_id int, pipeline_body string) int {
	len_orig_pipeline_name := len(pipeline_name)
	var hash_code string
	if len_orig_pipeline_name <= 50 {
		hash_code = Comm.RandSeq(60 - len_orig_pipeline_name)
	} else {
		hash_code = Comm.RandSeq(10)
	}
	pipeline_name = pipeline_name + ":" + hash_code

	new_pipeline := g.Map{"pipeline_name": pipeline_name, "group_id": group_id, "agent_id": agent_id, "body": pipeline_body, "updated_at": gtime.Now().Timestamp()}

	result, err := dao.CicdPipeline.Data(new_pipeline).Save()
	if err != nil {
		glog.Error(err)
	}

	pipeline_id, err := result.LastInsertId()
	if err != nil {
		glog.Error(err)
	}

	return int(pipeline_id)
}

func (s *pipelineService) Show(pipeline_id int) g.Map {
	result, err := dao.CicdPipeline.Where("id=", pipeline_id).One()
	if err != nil {
		glog.Error(err)
	}
	return result.Map()
}

func (s *pipelineService) Update(pipeline_id int, group_id int, agent_id int, pipeline_name string, pipeline_body string) error {
	new_pipeline := g.Map{"pipeline_name": pipeline_name, "body": pipeline_body, "agent_id": agent_id, "group_id": group_id}
	_, err := dao.CicdPipeline.Data(new_pipeline).Where(g.Map{"id": pipeline_id}).Update()
	if err != nil {
		glog.Error(err)
	}
	return nil
}
