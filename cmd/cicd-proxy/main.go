package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"

	"github.com/gogf/gf/os/glog"
	"github.com/koding/websocketproxy"
)

var (
	flagBackend    = flag.String("backend", "ws://localhost:8070/", "Backend URL for proxying")
	flagListenPort = flag.String("listenport", "localhost:8080", "Listen Port for proxying")
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	glog.Info(*flagBackend)
	glog.Info(*flagListenPort)

	u, err := url.Parse(*flagBackend)
	if err != nil {
		log.Fatalln(err)
	}

	err = http.ListenAndServe(*flagListenPort, websocketproxy.NewProxy(u))
	if err != nil {
		log.Fatalln(err)
	}
}
