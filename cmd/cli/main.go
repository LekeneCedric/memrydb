package main

import "github.com/LekeneCedric/memrydb/internal/server"

func main() {
	s := &server.Server{}
	s.Setup().Start()
}
