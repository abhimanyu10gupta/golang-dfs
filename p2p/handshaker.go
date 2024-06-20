package p2p

// HandshakeFunc is an interface that
type HandshakeFunc func(Peer) error

// NOPHandshakeFunc
func NOPHandshakeFunc(Peer) error {
	return nil
}
