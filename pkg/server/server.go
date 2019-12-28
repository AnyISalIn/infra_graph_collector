package server

import (
	"encoding/json"
	"fmt"
	"github.com/anyisalin/infra_graph_collector/pkg/structs"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func buildSelf(host, addr string, c chan string) {
	c <- fmt.Sprintf(`MERGE (host:IP {
	addr:'%s',
	hostname:'%s',
	type: 'VirtualMachine'
	})`, addr, host)
}

func buildDep(addr string, peer structs.Peer, c chan string) {
	c <- fmt.Sprintf(`MERGE (n {
            addr:'%s'
            })
    `, peer.Addr)

	c <- fmt.Sprintf(`MATCH (a:IP),(b:IP)
    WHERE a.addr = '%s' AND b.addr = '%s'
    MERGE (a)-[r:PEER]->(b)
    `, addr, peer.Addr)
}

func writeData(uri, username, password string, c chan string) {
	var (
		err     error
		driver  neo4j.Driver
		session neo4j.Session
		result  neo4j.Result
	)

	driver, err = neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Close()

	session, err = driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	for {
		sql := <-c
		fmt.Println(sql)
		_, err = session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
			result, err = transaction.Run(sql, map[string]interface{}{})
			return nil, result.Err()
		})
	}
}

func Run(listen string) {
	c := make(chan string)
	go writeData("bolt://13.20.0.95:7687", "neo4j", "passwd1Q", c)
	http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		var payload structs.Payload
		json.Unmarshal(body, &payload)
		buildSelf(payload.Self.Hostname, payload.Self.Addr, c)
		for _, peer := range payload.Peers {
			go buildDep(payload.Self.Addr, peer, c)
		}
	})
	log.Fatal(http.ListenAndServe(listen, nil))
}
