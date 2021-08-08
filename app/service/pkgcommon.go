package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gogf/gf/os/glog"
)

var Comm = commService{}

type commService struct{}

func (s *commService) RandSeq(randlen int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, randlen)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func stringToSlice(str string) []string {
	var newslice []string = make([]string, 0)
	err := json.Unmarshal([]byte(str), &newslice)
	if err != nil {
		glog.Error(err)
	}
	return newslice
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GetUserGroupIds(ctx context.Context) []string {
	var group_ids []string = make([]string, 0)
	if v := Context.Get(ctx); v != nil && v.User != nil {
		user_id := int(v.User.Id)
		group_ids := User.GetGroupId(user_id)
		return group_ids
	}
	return group_ids
}

func CheckAuthor(ctx context.Context, pipeline_id int) bool {
	if v := Context.Get(ctx); v != nil && v.User != nil {
		user_id := int(v.User.Id)
		group_id_user := User.GetGroupId(user_id)
		group_id_pipeline := Pipeline.GetGroupId(pipeline_id)
		return stringInSlice(fmt.Sprint(group_id_pipeline), group_id_user)
	}
	return false
}
