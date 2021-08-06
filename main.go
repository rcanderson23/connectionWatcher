package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rcanderson23/connectionWatcher/connections"
)

const (
	// TCP is the fs path to the tcp file
	TCP = "/proc/net/tcp"
	// TTL is the length of time in seconds that connections need to be tracked
	TTL = int64(60)
	// WaitPeriod is the amount of time in seconds before every observation of TCP
	WaitPeriod = time.Second * 10
)

// shortcut: configurable settings, such as prometheus listening port, TTL, period between observations, etc
// shortcut: no timeout on the blocking of a remote host, probably best to set some kind of TTL on host blocking instead of indefinitely
// shortcut: this is all done sequentially which probably isn't the most efficient. It _may_ be better to use channels, especially if
// reading directly from a pcap instead of polling. Channels generally make it more complex to read, we would want to ensure
// channels are actually faster and worth it.
func main() {
	blocker := connections.NewIPBlocker()
	cw := connections.NewConnectionWatcher(blocker)

	// seed connection watcher data, better UX to see logs right away rather than after the first ticker loop
	cw.Observe(TCP, time.Now().Unix())

	ticker := time.NewTicker(WaitPeriod)

	// create channel to gracefully terminate
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		for {
			select {
			case <-ticker.C:
				t := time.Now().Unix()
				cw.Observe(TCP, t)
				cw.Blocker.RemoveOldConnections(t, TTL)
				hosts := cw.Blocker.HostsToBlock()
				errs := cw.Blocker.BlockHosts(hosts)
				if len(errs) != 0 {
					for err := range errs {
						log.Printf("failed to block host: %v", err)
					}
				}
			case <-done:
				break
			}
		}
	}()

	go func() {
		// shortcut: hardcode path and port for metrics, production you may want to make these configurable
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":9090", nil))
	}()

	<-done

	// cleanup iptables added during runtime
	cw.Blocker.CleanUp()
}
