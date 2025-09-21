package dingy

import "github.com/google/uuid"

type Raft struct {
	id    NodeID
	state RaftState
}

type NodeID struct {
	// Identity component (stable across restarts)
	UUID uuid.UUID

	// Network component (can change)
	Address string // "host:port" for TCP

	// Generation component (monotonic)
	Epoch uint64 // Increments on rejoin
}

func NewRaft() *Raft {
	return &Raft{}
}

func (r *Raft) RequestVote(req VoteRequest) VoteResponse {
	// Get current state atomically
	currentState := r.state.getState()
	if currentState == nil {
		return VoteResponse{Term: 0, VoteGranted: false}
	}

	// Call the state's ProcessVote method
	newState, resp := (*currentState).ProcessVote(req)

	// Update state if it changed
	if newState != nil && newState != *currentState {
		r.state.setState(&newState)
	}

	return resp
}

func (r *Raft) AppendEntries(req AppendRequest) AppendResponse {
	// Get current state atomically
	currentState := r.state.getState()
	if currentState == nil {
		return AppendResponse{Term: 0, Success: false}
	}

	// Call the state's ProcessAppend method
	newState, resp := (*currentState).ProcessAppend(req)

	// Update state if it changed
	if newState != nil && newState != *currentState {
		r.state.setState(&newState)
	}

	return resp
}
