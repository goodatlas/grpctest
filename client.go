package grpctest

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"golang.org/x/net/context"

	"time"

	"google.golang.org/grpc"
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

// StartClient starts client
func StartClient(hostaddr, bindaddr string) {
	conn, err := grpc.Dial(hostaddr, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := NewCounterClient(conn)
	name := "client-" + randID(5)

	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		ir, err := c.Increment(context.Background(), &IncrementRequest{Name: name})

		if err != nil {
			log.Fatalf("could not increment: %v", err)
		}

		fmt.Fprintf(w, "Name: %s\nCount: %d\n", name, ir.Count)
	})

	http.ListenAndServe(bindaddr, nil)
}
