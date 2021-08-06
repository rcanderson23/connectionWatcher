# connectionWatcher
connectionWatcher polls `/proc/net/tcp` every 10 seconds to parse TCP connections. ConnectionWatcher implements a naive 
port scanner where if detects multiple connections from the same remote IP on multiple local ports. This block is done 
by inserting rules via iptables. Functionality is limited by what `/proc/net/tcp` provides, as such, it is possible that 
incoming connections are blocked because of outgoing connections made by the host. The best fix for this issue is to
watch packets and make decisions based on the TCP handshake rather than watch `/proc/net/tcp`. 

## Requirements
* Linux x86_64
* NET_ADMIN privileges (if you want IP blocking abilities) 

## Usage
### Source
```
git clone https://github.com/rcanderson23/connectionWatcher
cd connectionWatcher
make build
bin/connectionWatcher
```
