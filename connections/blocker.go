package connections

import (
	"fmt"
	"log"
	"net"
	"sort"
	"strings"

	"github.com/coreos/go-iptables/iptables"
)

const (
	// Filter is the table to target with blocking rules
	Filter = "filter"
	// Chain is the chain to target with blocking rules
	Chain = "INPUT"
)

// IPBlocker is used to track and block remote IPs based on the number of connections made in 30 seconds
// shortcut: ideally the number of connections and timeframe are configurable by the user
type IPBlocker struct {
	// Store `LocalIP:RemoteIP` mapped to the port and timestamp(unix epoch)
	IPPortTime   map[string]map[uint16]int64
	BlockedHosts []net.IP
	IP4Table     *iptables.IPTables
}

// NewIPBlocker returns a pointer to a newly constructed IPBlocker
func NewIPBlocker() *IPBlocker {
	return &IPBlocker{
		IPPortTime: make(map[string]map[uint16]int64),
		IP4Table:   NewIPv4Table(),
	}
}

// NewIPv4Table returns an IPTables to be used for host blocking. Checks that the ACCEPT and INPUT are present to be used
// shortcut: we are assuming the ACCEPT table is present along with the INPUT chain
func NewIPv4Table() *iptables.IPTables {
	ip4t, err := iptables.New()
	if err != nil {
		log.Printf("Failed to create iptables: %v. Host blocking is disabled.", err)
		return nil
	}

	return ip4t
}

// RemoveOldConnections checks the map IPPortTime and removes any entries that are older than the provided ttl and the unix time provided
// now and ttl are measured in seconds
func (ipb *IPBlocker) RemoveOldConnections(now int64, ttl int64) []uint16 {
	var removedPorts []uint16

	for ip, portMap := range ipb.IPPortTime {
		for port, ts := range portMap {
			lifetime := now - ts
			if lifetime >= ttl {
				removedPorts = append(removedPorts, port)
				delete(ipb.IPPortTime[ip], port)
			}
		}
	}

	return removedPorts
}

// AddPort updates the IPPortTime map with the port and time(unix epoch)
func (ipb *IPBlocker) AddPort(localIP string, remoteIP string, port uint16, t int64) {
	key := fmt.Sprintf("%s:%s", localIP, remoteIP)
	if _, present := ipb.IPPortTime[key]; !present {
		ipb.IPPortTime[key] = make(map[uint16]int64)
	}

	ipb.IPPortTime[key][port] = t
}

// HostsToBlock checks for any remote hosts that has 3 or more ports connected from its IP
// returns a slice of RemoteHost.
func (ipb *IPBlocker) HostsToBlock() []RemoteHost {
	var hosts []RemoteHost

	// key `LocalIP:RemoteIP` formatted string
	for key, portMap := range ipb.IPPortTime {
		if len(portMap) >= 3 {

			var ports []uint16
			for port := range portMap {
				ports = append(ports, port)

				delete(ipb.IPPortTime[key], port)
			}

			ips := strings.Split(key, ":")
			hosts = append(hosts, RemoteHost{
				RemoteIP: net.ParseIP(ips[1]),
				LocalIP:  net.ParseIP(ips[0]),
				Ports:    ports,
			})
		}
	}

	return hosts
}

// BlockHosts inserts iptables entry to block hosts
// Note: we are checking to make sure we don't block incoming from 0.0.0.0 or 127.0.0.1 but there may be a better set...
// shortcut: assuming iptables is in use here and not nftables
// shortcut: assuming default INPUT chain is available to use
func (ipb *IPBlocker) BlockHosts(hosts []RemoteHost) []error {
	var errs []error

	for _, host := range hosts {
		remoteIP := host.RemoteIP.String()

		if remoteIP != "0.0.0.0" && remoteIP != "127.0.0.1" && !ipb.isBlocked(host.RemoteIP) {
			log.Printf("Port scan detected: %s\n", host.String())

			if ipb.IP4Table != nil {
				err := ipb.insertRule(host.RemoteIP)
				if err != nil {
					log.Printf("Failed to insert block rule for %s: %v", remoteIP, err)
				}
			}
		}
	}

	return errs
}

func (ipb *IPBlocker) isBlocked(addr net.IP) bool {
	for _, ip := range ipb.BlockedHosts {
		if ip.Equal(addr) {
			return true
		}
	}
	return false
}

func (ipb *IPBlocker) insertRule(ip net.IP) error {
	exist, err := ipb.IP4Table.Exists(Filter, Chain, "-s", ip.String(), "-j", "DROP")
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	err = ipb.IP4Table.Insert("filter", "INPUT", 1, "-s", ip.String(), "-j", "DROP")
	if err != nil {
		return err
	}

	ipb.BlockedHosts = append(ipb.BlockedHosts, ip)

	return nil
}

// CleanUp is meant to clean up any iptable rules on the host
func (ipb *IPBlocker) CleanUp() {
	log.Printf("Cleaning up iptable entries made by connection watcher")
	for _, ip := range ipb.BlockedHosts {
		err := ipb.IP4Table.Delete(Filter, Chain, "-s", ip.String(), "-j", "DROP")
		if err != nil {
			log.Printf("Failed to remove iptable block for %s: %v", ip.String(), err)
		}
	}
}

// RemoteHost stores information needed to block and report blocked IPs
type RemoteHost struct {
	// Remote IP making connection to the local host
	RemoteIP net.IP

	// Local IP connected to from RemoteIP
	LocalIP net.IP

	// Ports the remote IP has been observed making connections to
	Ports []uint16
}

func (rh *RemoteHost) String() string {
	sort.Slice(rh.Ports, func(i, j int) bool { return rh.Ports[i] < rh.Ports[j] })
	return fmt.Sprintf("%s -> %s on ports %s", rh.RemoteIP, rh.LocalIP, portsToString(rh.Ports))
}

func portsToString(ports []uint16) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ports)), ","), "[]")
}
