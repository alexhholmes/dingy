package dingy

import (
	"encoding/gob"
	"fmt"
	"net"
	"net/rpc"
)

const (
	VoteRequestType = iota
	AppendEntriesType
)

// Transport is a TCP implementation for Raft RPCs.
type Transport struct {
	listener net.Listener
	node     *Raft
	peers    map[NodeID]*rpc.Client
}

func NewTCPTransport(node *Raft, address string) (*Transport, error) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return &Transport{
		listener: ln,
		node:     node,
		peers:    make(map[NodeID]*rpc.Client),
	}, nil
}

// Serve starts accepting incoming RPC connections on the listener
func (t *Transport) Serve() error {
	if t.listener == nil {
		return fmt.Errorf("no listener configured")
	}

	for {
		conn, err := t.listener.Accept()
		if err != nil {
			// Check if listener was closed
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				continue
			}
			return err
		}

		// Handle each connection in a goroutine
		go t.handle(conn)
	}
}

func (t *Transport) handle(conn net.Conn) {
	defer conn.Close()

	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)

	var msgType uint8
	decoder.Decode(&msgType)

	switch msgType {
	case VoteRequestType:
		var req VoteRequest
		decoder.Decode(&req)

		resp := t.node.RequestVote(req)
		encoder.Encode(resp)

	case AppendEntriesType:
		var req AppendRequest
		decoder.Decode(&req)

		resp := t.node.AppendEntries(req)
		encoder.Encode(resp)
	}
}

func (t *Transport) RequestVote(target NodeID, req VoteRequest) (VoteResponse, error) {
	conn, err := net.Dial("tcp", target.Address)
	if err != nil {
		return VoteResponse{}, err
	}
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	// Send message type
	if err := encoder.Encode(VoteRequestType); err != nil {
		return VoteResponse{}, err
	}

	// Send request
	if err := encoder.Encode(req); err != nil {
		return VoteResponse{}, err
	}

	// Read response
	var resp VoteResponse
	if err := decoder.Decode(&resp); err != nil {
		return VoteResponse{}, err
	}

	return resp, nil
}

func (t *Transport) AppendEntries(target NodeID, req AppendRequest) (AppendResponse, error) {
	conn, err := net.Dial("tcp", target.Address)
	if err != nil {
		return AppendResponse{}, err
	}
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	// Send message type
	if err := encoder.Encode(AppendEntriesType); err != nil {
		return AppendResponse{}, err
	}

	// Send request
	if err := encoder.Encode(req); err != nil {
		return AppendResponse{}, err
	}

	// Read response
	var resp AppendResponse
	if err := decoder.Decode(&resp); err != nil {
		return AppendResponse{}, err
	}

	return resp, nil
}

func (t *Transport) InstallNode(node *Raft) {
	t.node = node
}
