package service

import (
	"context"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/tangr/cicdgo/app/dao"
)

var Group = groupService{}

type groupService struct{}

type ListGroups struct {
	Id         int    `json:"id"`
	Group_name string `json:"groupName"`
}

func (s *groupService) ListGroups() []ListGroups {
	groups := ([]ListGroups)(nil)
	err := dao.CicdGroup.Fields("id,group_name").Structs(&groups)
	if err != nil {
		glog.Error(err)
	}
	return groups
}

func (s *groupService) GetGroupNames() gdb.Result {
	result, err := dao.CicdGroup.Fields("id,group_name").All()
	if err != nil {
		glog.Error(err)
	}
	return result
}

func (s *groupService) New(groupname string) int64 {
	newgroup := g.Map{
		"group_name": groupname,
		"parent_id":  0,
	}
	result, err := dao.CicdGroup.Data(newgroup).Save()
	if err != nil {
		glog.Error(err)
	}
	groupid, err := result.LastInsertId()
	if err != nil {
		glog.Error(err)
	}
	return groupid
}

func (s *groupService) Update(groupid string, groupname string) string {
	newgroup := g.Map{
		"group_name": groupname,
	}
	glog.Debug(newgroup)
	_, err := dao.CicdGroup.Data(newgroup).Where("id=", groupid).Update()
	if err != nil {
		glog.Error(err)
	}
	return groupid
}

func (s *groupService) GetGroupName(group_id string) string {
	group_name, err := dao.CicdGroup.Fields("group_name").Where("id=", group_id).Value()
	if err != nil {
		glog.Error(err)
	}
	return group_name.String()
}

func (s *groupService) GetScriptBody(script_name string) string {
	script_body, err := dao.CicdScript.Fields("script_body").Where("script_name=", script_name).Value()
	if err != nil {
		glog.Error(err)
	}
	glog.Error(script_body)
	return script_body.String()
}

func (s *groupService) IsSignedIn(ctx context.Context) bool {
	if v := Context.Get(ctx); v != nil && v.User != nil {
		return true
	}
	return false
}

func (s *groupService) SignOut(ctx context.Context) error {
	return Session.RemoveUser(ctx)
}
