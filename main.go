package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/pwzgorilla/miniraft/handler"
	"github.com/pwzgorilla/miniraft/proto"
	"github.com/pwzgorilla/miniraft/raft"
)

type NullWriter int

func (NullWriter) Write([]byte) (int, error) {
	return 0, nil
}

func ReadConfig(path string) (*raft.ClusterConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf raft.ClusterConfig
	err = json.Unmarshal(data, &conf)
	if err == nil && len(conf.Servers) < 1 {
		err = errors.New("No Server Configuration found")
	}
	return &conf, err
}

func main() {
	//log.SetOutput(new(NullWriter))

	//TODO: Read the config.json file to get all the server configurations
	//clusterConfig, err := ReadConfig("config.json")
	//if err != nil {
	//	log.Println("Error parsing config file : ", err.Error())
	//	return
	//}

	//starting the server
	//id := flag.Int("id", 1, "an int")
	listen := flag.String("listen", "", "listen address")
	members := flag.String("members", "localhost:3588", "members address")
	flag.Parse()

	memberList := strings.Split(*members, ",")
	var Members []*raft.Member
	for _, member := range memberList {
		host, port, err := net.SplitHostPort(member)

		if err != nil {
			log.Println("[ERROR]", err.Error())
			return
		}

		p, err := strconv.Atoi(port)
		if err != nil {
			log.Println("[ERROR]", err.Error())
			return
		}

		Members = append(Members, &raft.Member{
			IP:   host,
			Port: p,
		})
	}

	//addrs, _ := net.InterfaceAddrs()

	//var myip string
	//for _, member := range memberList {
	//	for _, addr := range addrs {
	//		if member == addr.(*net.IPNet).IP.String() {
	//			myip = member
	//		}
	//	}
	//}

	//if myip == "" {
	//	log.Println("Error: members not contain local address")
	//	return
	//}

	//id, _ := strconv.Atoi(strings.Replace(myip, ".", "", -1))
	host, _, err := net.SplitHostPort(*listen)
	if err != nil {
		log.Println("[ERROR]", err.Error())
		return
	}
	//id, _ := strconv.Atoi(strings.Replace(host, ".", "", -1))
	log.Println("Starting sevrer with ID ", host)

	commitCh := make(chan proto.LogEntry, 10000)
	raftInstance, err := raft.NewRaft(Members, host, commitCh)
	if err != nil {
		log.Println("Error creating server instance : ", err.Error())
	}

	raftInstance.EventInCh = make(chan proto.Event, 1000)

	//First entry in the ClusterConfig will be the default leader
	var clientPort int
	for _, server := range raftInstance.Members {
		//Initialize the connection handler module
		if server.IP == raftInstance.ServerID {
			clientPort = server.Port
			raftInstance.CurrentState = raft.FOLLOWER
		}
	}

	if clientPort <= 0 {
		log.Println("Server's client port not valid")
	} else {
		go handler.StartConnectionHandler(raftInstance.ServerID, clientPort, raftInstance.EventInCh, raftInstance.Mutex)
	}

	//Inititialize the KV Store Module
	go raft.InitializeKVStore(raftInstance.ServerID, raftInstance.CommitCh)

	//Initialize the server
	raftInstance.InitServer()

	log.Println("Started raft Instance for server ID ", raftInstance.ServerID)
}
