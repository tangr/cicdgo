package main

import (
	_ "github.com/tangr/cicdgo/boot"
	_ "github.com/tangr/cicdgo/router"

	"github.com/tangr/cicdgo/app/service"
)

func main() {
	service.AgentCICD.AgentRun()
}
