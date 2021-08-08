package connections

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/rcanderson23/connectionWatcher/metrics"
)

const (
	MinEphemeralPort = uint16(32768)
	MaxEphemeralPort = uint16(60999)
)

// Connection stores the connection tuple
type Connection struct {
	LocalIP    net.IP
	LocalPort  uint16
	RemoteIP   net.IP
	RemotePort uint16
}

// ConnectionWatcher holds connection state to be compared against and updated at every observation
type ConnectionWatcher struct {
	Connections map[string]Connection
	Blocker     *IPBlocker
	IgnoredIPs  []net.IP
}

// NewConnectionWatcher returns a pointer to a new ConnectionWatcher that includes the provided IPBlocker
// as well as a set of IPs that should not be inserted into the IPBlocker
func NewConnectionWatcher(blocker *IPBlocker) *ConnectionWatcher {
	return &ConnectionWatcher{
		Connections: make(map[string]Connection),
		Blocker:     blocker,
	}
}

// Observe opens the provided file path and populates the ConnectionWatcher structure with the contents of the file
func (cw *ConnectionWatcher) Observe(path string, t int64) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("failed to open file: %v", err)
		return
	}

	obsConns, err := getConnections(f)
	if err != nil {
		log.Printf("failed to check new connections: %v", err)
		return
	}

	cw.updateIPBlocker(obsConns, t)
	printNewConnections(obsConns, cw.Connections)

	// Connections are now equal to what was observed
	cw.Connections = obsConns
}

func (cw *ConnectionWatcher) updateIPBlocker(conns map[string]Connection, t int64) {
	for _, conn := range conns {
		cw.Blocker.AddPort(conn.LocalIP.String(), conn.RemoteIP.String(), conn.LocalPort, t)
	}
}

// shortcut: we are assuming that if the local port is in the default ephemeral range, it is the local host connecting
// out. This isn't guaranteed but increases the accuracy of the printed logs.
func printNewConnections(obs map[string]Connection, past map[string]Connection) {
	for i := range obs {
		if _, present := past[i]; !present {
			if obs[i].LocalPort >= MinEphemeralPort && obs[i].LocalPort <= MaxEphemeralPort {
				log.Printf("New connection %s:%d -> %s:%d\n", obs[i].LocalIP, obs[i].LocalPort, obs[i].RemoteIP, obs[i].RemotePort)
			} else {
				log.Printf("New connection %s:%d -> %s:%d\n", obs[i].RemoteIP, obs[i].RemotePort, obs[i].LocalIP, obs[i].LocalPort)
			}
			metrics.NewConnections.Inc()
		}
	}
}

// getConnections accepts an io.Reader and returns a map of the string of the connection tuple along with the data
// structure of the connection
func getConnections(r io.Reader) (map[string]Connection, error) {
	var lineCount int

	scanner := bufio.NewScanner(r)
	scanner.Scan() // skipping first line

	newConns := make(map[string]Connection)
	for scanner.Scan() {
		lineCount++
		line := strings.Split(strings.TrimSpace(scanner.Text()), " ")

		local := line[1]
		localIP, localPort, err := parseEndpoint(local)
		if err != nil {
			log.Printf("failed to parse local endpoint: %v", err)
			continue
		}

		remote := line[2]
		remoteIP, remotePort, err := parseEndpoint(remote)
		if err != nil {
			log.Printf("failed to parse destination endpoint: %v", err)
			continue
		}

		key := fmt.Sprintf("%s:%d:%s:%d", localIP.String(), localPort, remoteIP.String(), remotePort)
		newConns[key] = Connection{
			LocalIP:    localIP,
			LocalPort:  localPort,
			RemoteIP:   remoteIP,
			RemotePort: remotePort,
		}
	}

	err := scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("error scanning file: %v", err)
	}

	// if line counter is empty,
	if lineCount == 0 {
		return nil, fmt.Errorf("/proc/net/tcp was empty")
	}

	return newConns, nil
}
