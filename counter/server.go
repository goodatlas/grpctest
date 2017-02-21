package counter

import (
	"fmt"
	"log"
	"net"

	"sync"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	data  map[string]int32
	mutex *sync.Mutex
	out   *output
}

func newServer() *server {
	return &server{
		data:  make(map[string]int32),
		mutex: &sync.Mutex{},
		out:   &output{linesPrinted: 0},
	}
}

func (s *server) Increment(ctx context.Context, r *IncrementRequest) (*IncrementResponse, error) {
	s.mutex.Lock()

	c, ok := s.data[r.Name]

	if ok {
		s.data[r.Name]++
		c = s.data[r.Name]
	} else {
		s.data[r.Name], c = 1, 1
	}

	s.printData()

	s.mutex.Unlock()

	return &IncrementResponse{Count: c}, nil
}

func (s *server) printData() {

	s.out.reset()

	if s.out.linesPrinted == 0 {
		fmt.Println("Hits:")
	}

	for k, v := range s.data {
		fmt.Printf(" - %s: %d", k, v)
		fmt.Println("")
	}

	s.out.linesPrinted = len(s.data)
}

// Start runs RPC counter server
func Start(bindadrr string) {

	lis, err := net.Listen("tcp", bindadrr)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	RegisterCounterServer(s, newServer())

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type output struct {
	linesPrinted int
}

func (o *output) reset() {
	if o.linesPrinted > 0 {
		fmt.Printf("\x1b[%dA", o.linesPrinted)
	}
}
