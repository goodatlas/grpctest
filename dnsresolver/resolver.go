package dnsresolver

import (
	"fmt"
	"net"
	"strings"

	"time"

	"google.golang.org/grpc/naming"
)

// Resolver load balacer resolver
type Resolver struct {
	WatcherInterval time.Duration
	ResolverLookup  ResolverLookup
	watcher         *watcher
}

// NewResolver creates resolver instance
func NewResolver() naming.Resolver {
	return &Resolver{
		WatcherInterval: 200 * time.Millisecond,
		ResolverLookup:  defaultLookup,
	}
}

// Resolve creates a Watcher for target.
func (r *Resolver) Resolve(target string) (naming.Watcher, error) {
	var (
		next             = make(chan nextResults)
		firstLookup, err = r.ResolverLookup(target)
	)

	if err != nil {
		return nil, err
	}

	// Initial read
	go func() {
		next <- nextResults{results: firstLookup}
	}()

	r.watcher = &watcher{
		resolverLookup: r.ResolverLookup,
		currentAddrs:   make(map[string]interface{}),
		close:          make(chan struct{}),
		next:           next,
		interval:       r.WatcherInterval,
		target:         target,
	}

	go r.watcher.run()

	return r.watcher, nil
}

// A DNSResolver resolves addresses
type DNSResolver interface {
	Lookup(target string) ([]string, error)
}

// ResolverLookup function used to lookup dns records
type ResolverLookup func(target string) ([]string, error)

func defaultLookup(target string) ([]string, error) {
	var (
		host string
		port string
		fr   []string
	)

	if spl := strings.Split(target, ":"); len(spl) == 2 {
		host = spl[0]
		port = fmt.Sprintf(":%s", spl[1])
	} else {
		host = target
		port = ""
	}

	lr, err := net.LookupHost(host)

	if err != nil {
		return nil, err
	}

	for _, v := range lr {
		fr = append(fr, fmt.Sprintf("%s%s", v, port))
	}

	return fr, nil
}
