package connections

func IsEphemeralPort(port uint16) bool {
	return port >= MinEphemeralPort && port <= MaxEphemeralPort
}
