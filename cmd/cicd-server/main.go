package main

import (
	_ "github.com/tangr/cicdgo/boot"
	_ "github.com/tangr/cicdgo/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	s := g.Server("wscicd")
	s.Run()
}
