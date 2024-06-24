package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

//TCCPeer represents the reomte node over the TCP established connection.

type TCPPeer struct {
	// conn is the underlying connection of the peer
	net.Conn

	// if we dial and retrieve a connection => outbound == true
	// if we accept and retrieve a conn => outbound == false
	outbound bool

	wg *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		wg:       &sync.WaitGroup{},
	}
}

func (p *TCPPeer) CloseStream() {
	p.wg.Done()
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {

	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC, 1024),
	}

}

// Addr implments the transport interface returning the address
// the transport is accepting connections.
func (t *TCPTransport) Addr() string {
	return t.ListenAddr
}

// Consume implemenmts the Transport interface, which will return
// read-only channel for reading the incoming messages received from another
// peer in the network.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

// Close implements the transport interface
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial implements the transport interface.
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	go t.handleConn(conn, true)

	return nil

}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	log.Printf("TCP transport listening on port %s\n", t.ListenAddr)

	return nil

}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		// fmt.Printf("New Incoming Connection %+v\n", conn)

		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)

	if err = t.HandshakeFunc(peer); err != nil {
		// conn.Close()
		// fmt.Printf("TCP handshake error: %s\n", err)
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	//Read Loop
	for {
		rpc := RPC{}

		err := t.Decoder.Decode(conn, &rpc)

		if err != nil {
			return
		}

		rpc.From = conn.RemoteAddr().String()

		if rpc.Stream {
			peer.wg.Add(1)
			fmt.Printf("[%s] incoming stream, waiting ...\n", conn.RemoteAddr())
			peer.wg.Wait()
			fmt.Printf("[%s] stream closed, resuming read loop ...\n", conn.RemoteAddr())
			continue

		}

		t.rpcch <- rpc

	}
}
