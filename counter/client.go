package counter

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/goodatlas/grpctest/dnsresolver"

	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// Client ..."
type Client struct {
	CounterClient
	Name string
	Conn *grpc.ClientConn
}

// NewClient initializes a counter client"
func NewClient(upstreamaddr string, dnslb bool, kind string) (*Client, error) {
	var (
		conn     *grpc.ClientConn
		err      error
		dialOpts = []grpc.DialOption{grpc.WithInsecure()}
	)

	if dnslb {
		lb := grpc.RoundRobin(dnsresolver.NewResolver())
		dialOpts = append(dialOpts, grpc.WithBalancer(lb))
	}

	conn, err = grpc.Dial(upstreamaddr, dialOpts...)

	if err != nil {
		return nil, err
	}

	grpclog.Printf("Connected to upstream: %s\n", upstreamaddr)

	return &Client{
		CounterClient: NewCounterClient(conn),
		Name:          fmt.Sprintf("%s-%s", kind, getID()),
		Conn:          conn,
	}, nil
}

func getID() string {
	if hostname := os.Getenv("HOSTNAME"); hostname != "" {
		return hostname
	}

	return randID(12)
}

func randID(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	var anchars = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)

	for i := range b {
		b[i] = anchars[rand.Intn(len(anchars))]
	}

	return string(b)
}
