package agent

import (
	"bytes"
	"encoding/json"
	"github.com/anyisalin/infra_graph_collector/pkg/structs"
	"github.com/cakturk/go-netstat/netstat"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func getPeers(selfIP string, allowPrefixList []string) ([]structs.Peer, error) {
	peers := make([]structs.Peer, 0)

	socks, err := netstat.TCPSocks(netstat.NoopFilter)
	if err != nil {
		return nil, err
	}
	for _, e := range socks {
		isAllow := false
		remoteIP := e.RemoteAddr.IP.String()
		if remoteIP == selfIP {
			continue
		}
		for _, allowPrefix := range allowPrefixList {
			if strings.HasPrefix(remoteIP, allowPrefix) {
				isAllow = true
				break
			}
		}
		if isAllow {
			peers = append(peers, structs.Peer{Addr: remoteIP, Port: e.RemoteAddr.Port})
		}
	}

	return peers, nil
}

func Run(endpoint, selfIPPrefix, allowPrefixListString string) {
	hostname, _ := os.Hostname()
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	var selfIP string
	allowPrefixList := strings.Split(allowPrefixListString, ",")
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil && strings.HasPrefix(ip.String(), selfIPPrefix) {
				selfIP = ip.String()
			}
		}
	}
	if selfIP == "" {
		log.Fatal("can not find self ip")
	}
	peers, err := getPeers(selfIP, allowPrefixList)
	if err != nil {
		log.Fatal(err)
	}
	payload := &structs.Payload{
		Self:  structs.Peer{Hostname: hostname, Addr: selfIP},
		Peers: peers,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	_, err = http.Post(endpoint, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
}
