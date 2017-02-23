package dnsresolver

import (
	"time"

	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/naming"
)

type nextResults struct {
	results []string
	err     error
}

type watcher struct {
	resolverLookup ResolverLookup
	currentAddrs   map[string]interface{}
	close          chan struct{}
	next           chan nextResults
	target         string
	interval       time.Duration
}

// Next blocks until an update or error happens. It may return one or more
// updates. The first call should get the full set of the results. It should
// return an error if and only if Watcher cannot recover.
func (w *watcher) Next() ([]*naming.Update, error) {
	var (
		updates []*naming.Update
		nr      = <-w.next
	)

	if nr.err != nil {
		return updates, nr.err
	}

	return w.generateUpdates(nr.results), nil
}

// Close closes the Watcher.
func (w *watcher) Close() {
	close(w.close)
}

func (w *watcher) generateUpdates(results []string) []*naming.Update {
	resultMap := make(map[string]interface{})

	for _, addr := range results {
		resultMap[addr] = nil
	}

	var (
		addOps    = w.updateOperations(resultMap, w.currentAddrs, naming.Add)
		deleteOps = w.updateOperations(w.currentAddrs, resultMap, naming.Delete)
	)

	w.currentAddrs = resultMap

	return append(addOps, deleteOps...)
}

func (w *watcher) updateOperations(a, b map[string]interface{}, op naming.Operation) []*naming.Update {
	updates := make([]*naming.Update, 0)

	for addr := range a {
		if _, ok := b[addr]; !ok {
			updates = append(updates, &naming.Update{Op: op, Addr: addr})

			action := "added to "
			if op == naming.Delete {
				action = "removed from"
			}

			grpclog.Printf("dnsresolver: Address %s %s %s", addr, action, w.target)
		}
	}

	return updates
}

func (w *watcher) getNextResults() nextResults {
	results, err := w.resolverLookup(w.target)
	return nextResults{results, err}
}

func (w *watcher) run() {
	for {
		select {
		case <-w.close:
			w.next <- nextResults{results: []string{}}
			close(w.next)
			return
		case <-time.After(w.interval):
			w.next <- w.getNextResults()
		}
	}
}
