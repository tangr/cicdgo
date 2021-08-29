package service

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/tangr/cicdgo/app/dao"
)

var Script = scriptService{}

type scriptService struct{}

type ListScripts struct {
	Id          int    `json:"script_id"`
	Script_name string `json:"script_name"`
	Author      string `json:"author"`
	Updated_at  int    `json:"updated_at"`
}

func (s *scriptService) ListScripts() []ListScripts {
	scripts := ([]ListScripts)(nil)
	err := dao.CicdScript.Fields("id,script_name,author,updated_at").Structs(&scripts)
	if err != nil {
		glog.Error(err)
	}
	return scripts
}

func (s *scriptService) GetScriptName(script_id string) string {
	script_name, err := dao.CicdScript.Fields("script_name").Where("id=", script_id).Value()
	if err != nil {
		glog.Error(err)
	}
	return script_name.String()
}

func (s *scriptService) GetScriptBody(script_name string) string {
	script_body, err := dao.CicdScript.Fields("script_body").Where("script_name=", script_name).Value()
	if err != nil {
		glog.Error(err)
	}
	return script_body.String()
}

func (s *scriptService) New(script_name string, script_body string) int {
	new_script := g.Map{"script_name": script_name, "script_body": script_body, "updated_at": gtime.Now().Timestamp()}
	result, err := dao.CicdScript.Data(new_script).Save()
	if err != nil {
		glog.Error(err)
	}
	script_id, err := result.LastInsertId()
	if err != nil {
		glog.Error(err)
	}
	return int(script_id)
}

func (s *scriptService) Show(script_id int) g.Map {
	result, err := dao.CicdScript.Where("id=", script_id).One()
	if err != nil {
		glog.Error(err)
	}
	return result.Map()
}

func (s *scriptService) Update(script_id int, script_name string, script_body string) error {
	new_script := g.Map{"script_name": script_name, "script_body": script_body, "updated_at": gtime.Now().Timestamp()}
	_, err := dao.CicdScript.Data(new_script).Where(g.Map{"id": script_id}).Update()
	if err != nil {
		glog.Error(err)
	}
	return nil
}
