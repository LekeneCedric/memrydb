package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"github.com/LekeneCedric/memrydb/internal/config"
	"github.com/LekeneCedric/memrydb/internal/protocol"
	"github.com/LekeneCedric/memrydb/internal/storage"
)

type Server struct {
	config  *config.Config
	storage storage.Engine
}

func (s *Server) Setup() *Server {
	s.config = loadConfig()
	s.storage = storage.NewSharedMap(s.config.NumberOfShard)
	return s
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Println("> server start", listener.Addr())
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("[Server error]: Failed to establish connexion\n %s", err.Error())
			continue
		}
		go handleConnexion(conn, s)
	}
}

func handleConnexion(conn net.Conn, s *Server) {
	defer conn.Close()
	buffer := make([]byte, 1024*1000)
	n, err := conn.Read(buffer)
	if err != nil {
		conn.Write([]byte(err.Error()))
		return
	}
	request, err := protocol.DecryptQuery(buffer[:n])
	if err != nil {
		conn.Write([]byte(err.Error()))
		return
	}
	switch protocol.RequestType(request.Method) {
	case protocol.GET:
		res := s.storage.Get(request.Key)
		conn.Write(res)
	case protocol.SET:
		s.storage.Set(request.Key, request.Value)
		conn.Write([]byte("ok"))
	case protocol.DEL:
		s.storage.Remove(request.Key)
		conn.Write([]byte("ok"))
	default:
		return
	}
}

func loadConfig() *config.Config {
	fconf, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer fconf.Close()
	conf, err := config.NewConfig(fconf)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}
