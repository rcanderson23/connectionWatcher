# connectionWatcher
connectionWatcher polls `/proc/net/tcp` every 10 seconds to parse TCP connections. ConnectionWatcher implements a naive 
port scanner where if detects multiple connections from the same remote IP on multiple local ports. This block is done 
by inserting rules via iptables. Functionality is limited by what `/proc/net/tcp` provides, as such, it is possible that 
incoming connections are blocked because of outgoing connections made by the host. There is logic to check if the 
local port is in the ephemeral range and treat this as an outgoing connection. This isn't 100% accurate for this 
determination. The best fix for this issue is to watch packets and make decisions based on the TCP handshake rather 
than watch `/proc/net/tcp`. 

## Requirements
* Linux x86_64
* Root privileges
* Go 1.16.x (if running from source, have not tested other versions)

## Usage
### Docker
1. [Install Docker](https://docs.docker.com/engine/install/)
```
docker run --network host --cap-drop ALL --cap-add NET_ADMIN --cap-add NET_RAW ghcr.io/rcanderson23/connectionwatcher:v0.1.1
```
### Binary
``` 
wget https://github.com/rcanderson23/connectionWatcher/releases/download/v0.1.1/connectionWatcher
chmod +x connectionWatcher
./connectionWatcher
```
### Source
The `main` branch should be treated as development and can be unstable. Use tagged branches or pre-compiled binaries
for production.
1. [Install Go](https://golang.org/doc/install)
```
git clone https://github.com/rcanderson23/connectionWatcher
cd connectionWatcher
make build
bin/connectionWatcher
```
