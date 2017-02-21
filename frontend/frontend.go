package frontend

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"golang.org/x/net/context"

	"time"

	"google.golang.org/grpc"

	"github.com/goodatlas/grpctest/counter"
)

func randID(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Start starts frontend service
func Start(bindaddr, upstreamaddr string) {
	conn, err := grpc.Dial(upstreamaddr, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	fmt.Printf("Connected to upstream: %s\n", upstreamaddr)

	defer conn.Close()
	c := counter.NewCounterClient(conn)
	name := "client-" + randID(5)

	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		ir, err := c.Increment(context.Background(), &counter.IncrementRequest{Name: name})

		if err != nil {
			log.Fatalf("could not increment: %v", err)
		}

		fmt.Fprintf(w, "Name: %s\nCount: %d\n", name, ir.Count)
	})

	fmt.Printf("Binding to: %s\n", bindaddr)
	http.ListenAndServe(bindaddr, nil)
}