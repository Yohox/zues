package main

import (
	"zues/src/server"
)

var servers = []*server.Server{
	server.ClientServer,
	server.InternalServer,
}

func StartServers(){
	for _, s := range servers {
		go s.Start()
	}
}

func StopServer() {
	for _, s := range servers {
		s.Stop()
	}
}

func main(){
	StartServers()
	select {}
}
