package main

import (
	"time"

	_ "github.com/tangr/cicdgo/boot"
	_ "github.com/tangr/cicdgo/router"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gsession"
)

func main() {
	s := g.Server("console")

	s.SetConfigWithMap(g.Map{
		"SessionMaxAge":  time.Minute * 60 * 24 * 7,
		"SessionStorage": gsession.NewStorageFile(g.Cfg().GetString("server.console.SessionPath")),
	})

	s.Run()
}
