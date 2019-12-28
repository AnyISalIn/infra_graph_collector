package main

import (
	"flag"
	"github.com/anyisalin/infra_graph_collector/pkg/server"
)

func main() {
	listen := flag.String("listen", ":8080", "http listen address, like :8080")
	flag.Parse()
	server.Run(*listen)
}
