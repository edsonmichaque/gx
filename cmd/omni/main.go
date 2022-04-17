package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/edsonmichaque/omni/internal/dispatcher"
	"github.com/edsonmichaque/omni/libomni"
	"github.com/edsonmichaque/omni/libomni/registry"
)

type TCPServer struct {
	Providers map[string]libomni.Omni
	Port      string
}

type Server struct {
	TCPServer TCPServer
}

func (s Server) ServeTCP() error {
	l, err := net.Listen("tcp4", fmt.Sprintf(":%s", s.TCPServer.Port))
	if err != nil {
		log.Println(err)
		return err
	}

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		log.Println("accepting a new request")

		dispatcher := dispatcher.New(c, registry.TCPProviders())

		log.Println("dispatching a new request")
		go dispatcher.Dispatch()
	}
}

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	srv := Server{
		TCPServer: TCPServer{
			Providers: registry.TCPProviders(),
			Port:      args[1],
		},
	}

	log.Println("starting servers...")
	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		log.Println("starting TCP server on port", srv.TCPServer.Port)
		if err := srv.ServeTCP(); err != nil {
			log.Println(err)
		}
	}(&wg)

	wg.Wait()
}
