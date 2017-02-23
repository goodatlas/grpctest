package frontend

import (
	"fmt"
	"log"
	"net/http"

	"google.golang.org/grpc/grpclog"

	"golang.org/x/net/context"

	"github.com/goodatlas/grpctest/counter"
)

// Start starts frontend service
func Start(bindaddr, upstreamaddr string, dnslb bool) {
	c, err := counter.NewClient(upstreamaddr, dnslb, "frontend")

	if err != nil {
		log.Fatal(err.Error())
	}

	defer c.Conn.Close()

	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		ir, err := c.Increment(context.Background(), &counter.IncrementRequest{Name: c.Name})

		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Fprintf(w, "Name: %s\nCount: %d\n", c.Name, ir.Count)
	})

	grpclog.Printf("Binding to: %s\n", bindaddr)
	http.ListenAndServe(bindaddr, nil)
}
