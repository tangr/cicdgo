package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"

	"github.com/gogf/gf/frame/g"
	"github.com/koding/websocketproxy"
)

var (
	flagBackend    = flag.String("backend", "ws://localhost:8070/", "Backend URL for proxying")
	flagListenPort = flag.String("listenport", "localhost:8080", "Listen Port for proxying")
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	g.Log().Info(*flagBackend)
	g.Log().Info(*flagListenPort)

	u, err := url.Parse(*flagBackend)
	if err != nil {
		log.Fatalln(err)
	}

	err = http.ListenAndServe(*flagListenPort, websocketproxy.NewProxy(u))
	if err != nil {
		log.Fatalln(err)
	}
}
