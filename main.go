package main

import "github.com/LekeneCedric/memrydb/cmd/server"

func main() {
	s := &server.Server{}
	s.Setup().Start()
}
