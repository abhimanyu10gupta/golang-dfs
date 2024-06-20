package p2p

import "net"

// Message holds any arbitrary data that is being sent
// over the Transport between two nodes in the network.
type Message struct {
	From    net.Addr
	Payload []byte
}
