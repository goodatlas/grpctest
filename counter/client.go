package counter

import (
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/grpc"
)

// Client ..."
type Client struct {
	CounterClient
	Name string
	Conn *grpc.ClientConn
}

// NewClient initializes a counter client"
func NewClient(upstreamaddr, kind string) (*Client, error) {
	conn, err := grpc.Dial(upstreamaddr, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	fmt.Printf("Connected to upstream: %s\n", upstreamaddr)

	return &Client{
		CounterClient: NewCounterClient(conn),
		Name:          fmt.Sprintf("%s-%s", kind, randID(5)),
		Conn:          conn,
	}, nil
}

func randID(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
