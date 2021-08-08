# connectionWatcher
connectionWatcher polls `/proc/net/tcp` every 10 seconds to parse TCP connections. ConnectionWatcher implements a naive 
port scanner where if detects multiple connections from the same remote IP on multiple local ports. This block is done 
by inserting rules via iptables. Functionality is limited by what `/proc/net/tcp` provides, as such, it is possible that 
incoming connections are blocked because of outgoing connections made by the host. The best fix for this issue is to
watch packets and make decisions based on the TCP handshake rather than watch `/proc/net/tcp`. 

## Requirements
* Linux x86_64
* Run as root
* Go 1.16.x (if running from source, have not tested other versions)

## Usage
### Source
1. [Install Go](https://golang.org/doc/install)
```
git clone https://github.com/rcanderson23/connectionWatcher
cd connectionWatcher
make build
bin/connectionWatcher
```
