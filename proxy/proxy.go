package proxy

import (
	"log"
	"net"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/goodatlas/grpctest/counter"
)

type proxy struct {
	client *counter.Client
}

func (p *proxy) Increment(ctx context.Context, r *counter.IncrementRequest) (*counter.IncrementResponse, error) {
	_, err := p.client.Increment(context.Background(), &counter.IncrementRequest{Name: p.client.Name})

	if err != nil {
		return nil, err
	}

	ir, err := p.client.Increment(context.Background(), &counter.IncrementRequest{Name: r.Name})

	if err != nil {
		return nil, err
	}

	return &counter.IncrementResponse{Count: ir.Count}, nil
}

// Start starts frontend service
func Start(bindaddr, upstreamaddr string) {
	c, err := counter.NewClient(upstreamaddr, "proxy")

	if err != nil {
		log.Fatal(err.Error())
	}

	lis, err := net.Listen("tcp", bindaddr)

	if err != nil {
		log.Fatal(err.Error())
	}

	s := grpc.NewServer()

	counter.RegisterCounterServer(s, &proxy{client: c})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatal(err.Error())
	}
}
