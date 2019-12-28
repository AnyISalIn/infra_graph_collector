package main

import (
	"flag"
	"fmt"
	"github.com/anyisalin/infra_graph_collector/pkg/agent"
	"strings"
)

func main() {
	endpoint := flag.String("endpoint", "http://localhost:8080", "http api endpoint url")
	selfIPPrefix := flag.String("self_ip_prefix", "", "self ip prefix, like 192.168")
	allowPrefixListString := flag.String("allow_prefix_list", "", "allow prefix list, like 192.168,13.20,172.16")
	flag.Parse()

	if !strings.HasPrefix(*endpoint, "http://") && !strings.HasPrefix(*endpoint, "https://") {
		panic(fmt.Sprintf("%s valid faild, http endpoint must starswith http", *endpoint))
	}

	agent.Run(*endpoint, *selfIPPrefix, *allowPrefixListString)
}
