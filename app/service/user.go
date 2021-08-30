package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/tangr/cicdgo/app/dao"
	"github.com/tangr/cicdgo/app/model"
	"golang.org/x/crypto/bcrypt"
)

var User = userService{}

type userService struct{}

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

var AuthorEnable = g.Cfg().GetBool("server.console.AuthorEnable")
var adminGroupName string = g.Cfg().GetString("server.console.AdminGroup")
var adminGroupId int

func (s *userService) ListUsers() []ListUsers {
	users := ([]ListUsers)(nil)
	err := dao.CicdUser.Fields("id,email,updated_at").Structs(&users)
	if err != nil {
		g.Log().Error(err)
	}
	return users
}

func (s *userService) New(username string, groups []string, password string) int64 {
	var passhash []byte
	var err error
	if password != "" {
		passhash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			g.Log().Error(err)
			// g.Log().Debug(passhash)
		}
	} else {
		passhash = []byte(password)
	}

	newuser := g.Map{
		"email":      username,
		"password":   string(passhash),
		"group_id":   groups,
		"updated_at": gtime.Now().Timestamp(),
	}
	result, err := dao.CicdUser.Data(newuser).Save()
	if err != nil {
		g.Log().Error(err)
	}
	userid, err := result.LastInsertId()
	if err != nil {
		g.Log().Error(err)
	}
	return userid
}

func (s *userService) Update(user_id string, username string, groups []string, password string) string {
	var newuser g.Map
	if password != "" {
		passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			g.Log().Error(err)
		}
		g.Log().Info(passhash)
		newuser = g.Map{
			"email":      username,
			"password":   passhash,
			"group_id":   groups,
			"updated_at": gtime.Now().Timestamp(),
		}
	} else {
		newuser = g.Map{
			"email":      username,
			"group_id":   groups,
			"updated_at": gtime.Now().Timestamp(),
		}
	}

	_, err := dao.CicdUser.Data(newuser).Where("id=", user_id).Update()
	if err != nil {
		g.Log().Error(err)
	}
	return user_id
}

func (s *userService) GetUserName(user_id string) string {
	email_name, err := dao.CicdUser.Fields("email").Where("id=", user_id).Value()
	if err != nil {
		g.Log().Error(err)
	}
	return email_name.String()
}

func (s *userService) SetAdminGroupId() {
	group_name := adminGroupName
	group_id, err := dao.CicdGroup.Fields("id").Where("group_name=", group_name).Value()
	if err != nil {
		g.Log().Error(err)
	}
	adminGroupId = group_id.Int()
}

func (s *userService) CheckUserAdmin(user_id uint) bool {
	if adminGroupId == 0 {
		s.SetAdminGroupId()
	}
	group_ids_str, err := dao.CicdUser.Fields("group_id").Where("id=", user_id).Value()
	if err != nil {
		g.Log().Error(err)
	}
	type GroupIds []string
	group_ids := new(GroupIds)
	// g.Log().Debug(group_ids_str)
	// g.Log().Debug(group_ids_str.Array())
	// g.Log().Debug(group_ids_str.String())
	group_ids_byte := []byte(group_ids_str.String())
	err = json.Unmarshal(group_ids_byte, group_ids)
	if err != nil {
		g.Log().Error(err)
	}
	// g.Log().Debug(group_ids)

	adminGroupIdString := fmt.Sprint(adminGroupId)
	return stringInSlice(adminGroupIdString, *group_ids)
}

func (s *userService) GetUser(user_id string) *GetUser {
	var user *GetUser
	err := dao.CicdUser.Fields("email,group_id").Where("id=", user_id).Struct(&user)
	if err != nil {
		g.Log().Error(err)
	}
	return user
}

func (s *userService) GetGroupId(user_id int) []string {
	var group_ids *UserGroupIds
	err := dao.CicdUser.Fields("group_id").Where("id=", user_id).Struct(&group_ids)
	if err != nil {
		g.Log().Error(err)
	}
	new_group_ids := group_ids.GroupId
	new_group_slice := stringToSlice(new_group_ids)
	return new_group_slice
}

func (s *userService) SignIn(ctx context.Context, username, password string) error {
	var user *UserInfo
	var userSession = new(model.UserApiSession)
	err := dao.CicdUser.Fields("id,email,password").Where("email=?", username).Struct(&user)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not exist")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		g.Log().Debug("user or passwd error")
		return errors.New("user or passwd error")
	} else {
		g.Log().Debug("pw ok")
	}

	userSession.Id = user.Id
	userSession.Email = user.Email
	isAdmin := s.CheckUserAdmin(uint(user.Id))
	userSession.IsAdmin = isAdmin
	if err := Session.SetUser(ctx, userSession); err != nil {
		g.Log().Error(err)
		return err
	}
	Context.SetUser(ctx, &model.ContextUser{
		Id:    uint(user.Id),
		Email: user.Email,
	})
	return nil
}

func (s *userService) IsSignedIn(ctx context.Context) bool {
	if v := Context.Get(ctx); v != nil && v.User != nil {
		return true
	}
	return false
}

func (s *userService) IsAdmin(ctx context.Context) bool {
	if !AuthorEnable {
		return true
	}
	if v := Context.Get(ctx); v != nil && v.User != nil && s.CheckUserAdmin(v.User.Id) {
		return true
	}
	return false
}

func (s *userService) SignOut(ctx context.Context) error {
	return Session.RemoveUser(ctx)
}
