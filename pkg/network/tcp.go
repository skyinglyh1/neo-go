package network

import (
	"io"
	"net"
)

func listenTCP(s *Server, port string) error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go handleConnection(s, conn, true)
	}
}

func connectToRemoteNode(s *Server, address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		s.logger.Printf("failed to connects to remote node %s", address)
		if conn != nil {
			conn.Close()
		}
		return
	}
	s.logger.Printf("connected to %s", conn.RemoteAddr())
	go handleConnection(s, conn, false)
}

func connectToSeeds(s *Server, addrs []string) {
	for _, addr := range addrs {
		go connectToRemoteNode(s, addr)
	}
}

func handleConnection(s *Server, conn net.Conn, initiated bool) {
	peer := NewPeer(conn, initiated)
	s.register <- peer

	// remove the peer from connected peers and cleanup the connection.
	defer func() {
		s.unregister <- peer
		conn.Close()
	}()

	// Start a goroutine that will handle all writes to the registered peer.
	go peer.writeLoop()

	// Read from the connection and decode it into an RPCMessage and
	// tell the server there is message available for proccesing.
	for {
		msg := &Message{}
		if err := msg.decode(conn); err != nil {
			// remote connection probably closed.
			if err == io.EOF {
				s.logger.Printf("conn read error: %s", err)
				break
			}
			// remove this node on any decode errors.
			s.logger.Printf("RPC :: decode error %s", err)
			break
		}
		s.message <- messageTuple{peer, msg}
	}
}
