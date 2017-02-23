package main

import (
	"flag"
	"log"

	"github.com/goodatlas/grpctest/counter"
	"github.com/goodatlas/grpctest/frontend"
	"github.com/goodatlas/grpctest/proxy"
)

func main() {
	upstreamaddr := flag.String("upstream", "localhost:50051", "Upstream service address")
	bindaddr := flag.String("bind", "localhost:50051", "Address for binding service")
	dnslb := flag.Bool("dnslb", false, "Use client load balancer")

	flag.Parse()

	switch t := flag.Arg(0); t {
	case "frontend":
		frontend.Start(*bindaddr, *upstreamaddr, *dnslb)
	case "counter":
		counter.Start(*bindaddr)
	case "proxy":
		proxy.Start(*bindaddr, *upstreamaddr, *dnslb)
	case "":
		log.Fatal("Missing type argument")
	default:
		log.Fatalf("Unknown type argument: %s", t)
	}
}
