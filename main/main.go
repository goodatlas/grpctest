package main

import (
	"errors"
	"flag"
	"log"

	"github.com/goodatlas/grpctest"
)

func getType(args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("Missing type argument")
	}

	if args[0] != "server" && args[0] != "client" {
		return "", errors.New("Unknown type argument: " + args[0])
	}

	return args[0], nil
}

func main() {
	hostaddr := flag.String("hostaddr", "localhost:50051", "Host service address")
	bindaddr := flag.String("bindaddr", "localhost:50051", "Address for binding service")
	flag.Parse()

	t, err := getType(flag.Args())

	if err != nil {
		log.Fatal(err.Error())
	}

	switch t {
	case "client":
		grpctest.StartClient(*hostaddr, *bindaddr)
	case "server":
		grpctest.StartServer(*bindaddr)
	}
}
