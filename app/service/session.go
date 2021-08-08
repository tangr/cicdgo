package service

import (
	"context"

	"github.com/tangr/cicdgo/app/model"
)

var Session = sessionService{}

type sessionService struct{}

const (
	sessionKeyUser = "SessionKeyUser"
)

func (s *sessionService) SetUser(ctx context.Context, user *model.UserApiSession) error {
	return Context.Get(ctx).Session.Set(sessionKeyUser, user)
}

func (s *sessionService) GetUser(ctx context.Context) *model.CicdUser {
	customCtx := Context.Get(ctx)
	if customCtx != nil {
		if v := customCtx.Session.GetVar(sessionKeyUser); !v.IsNil() {
			var user *model.CicdUser
			_ = v.Struct(&user)
			return user
		}
	}
	return nil
}

func (s *sessionService) RemoveUser(ctx context.Context) error {
	customCtx := Context.Get(ctx)
	if customCtx != nil {
		return customCtx.Session.Remove(sessionKeyUser)
	}
	return nil
}
