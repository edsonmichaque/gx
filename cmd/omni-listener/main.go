package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/edsonmichaque/libomni"
	"github.com/edsonmichaque/libomni/registry"
	"github.com/edsonmichaque/omni/internal/listener/dispatcher"
	"github.com/google/uuid"
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
		session := libomni.Session{
			ID:        uuid.NewString(),
			Timestamp: time.Now(),
		}

		c, err := l.Accept()
		if err != nil {
			fmt.Println(session.ID, err)
			continue
		}

		log.Println(session.ID, "accepting a new request")

		dispatcher := dispatcher.New(c, registry.TCPProviders())

		log.Println(session.ID, "dispatching a new request")

		go func() {
			defer close(session, c)

			if err := dispatcher.Dispatch(session); err != nil {
				log.Println(session.ID, "aborting request due to error", err)
			}
		}()
	}
}

func close(session libomni.Session, c io.Closer) {
	log.Println(session.ID, "closing the connection")
	log.Println(session.ID, "session duration", time.Since(session.Timestamp))

	if err := c.Close(); err != nil {
		log.Println(session.ID, "could not close the connection")
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
