package dingy

import (
	"sync"
	"sync/atomic"
)

type RaftState struct {
	// The current term, cache of StableStore
	currentTerm atomic.Uint64

	// Highest committed log entry
	commitIndex atomic.Uint64

	// Last applied log to the FSM
	lastApplied atomic.Uint64

	// Protects next four fields
	// Cache the latest snapshot index/term and latest log from LogStore
	lastLock          sync.Mutex
	lastSnapshotIndex uint64
	lastSnapshotTerm  uint64
	lastLogIndex      uint64
	lastLogTerm       uint64

	// Tracks running goroutines
	routines sync.WaitGroup

	State
}

func (r *RaftState) getState() State {
	return r.State
}

func (r *RaftState) setState(state *State) {
	r.State = *state
}

// Command represents an operation to be executed by the state machine
type Command struct {
	Type    CommandType
	Payload interface{}
}

type CommandType uint8

const (
	CommandNoop CommandType = iota
	CommandConfiguration
	CommandLogEntry
)

// VoteRequest is sent by candidates to gather votes.
// Contains candidate's log state for up-to-date check.
type VoteRequest struct {
	Term         uint64 // candidate's term
	CandidateID  string // candidate requesting vote
	LastLogIndex uint64 // index of candidate's last log entry
	LastLogTerm  uint64 // term of candidate's last log entry
}

// VoteResponse contains vote decision and responder's term.
// Higher term causes candidate to step down.
type VoteResponse struct {
	Term        uint64 // current term, for candidate to update
	VoteGranted bool   // true if candidate received vote
}

// AppendRequest serves as both heartbeat and log replication.
// PrevLog fields ensure log consistency before appending.
type AppendRequest struct {
	Term         uint64     // leader's term
	LeaderID     string     // so follower can redirect clients
	PrevLogIndex uint64     // index of log entry immediately preceding new ones
	PrevLogTerm  uint64     // term of prevLogIndex entry
	Entries      []LogEntry // log entries to store (empty for heartbeat)
	LeaderCommit uint64     // leader's commitIndex
}

// AppendResponse indicates success and current term.
// Failure means log inconsistency at PrevLogIndex.
type AppendResponse struct {
	Term    uint64 // current term for leader to update
	Success bool   // true if follower contained entry matching prevLog
}

// LogEntry represents a single entry in the replicated log
type LogEntry struct {
	Index uint64  // position in log
	Term  uint64  // term when entry was received by leader
	Type  LogType // command, configuration, or no-op
	Data  []byte  // serialized command
}

type LogType uint8

const (
	LogCommand LogType = iota
	LogConfiguration
	LogNoop
)

type State interface {
	ProcessTimeout() (State, []Command)
	ProcessVote(VoteRequest) (State, VoteResponse)
	ProcessAppend(AppendRequest) (State, AppendResponse)
}

var (
	_ State = (*LeaderState)(nil)
	_ State = (*FollowerState)(nil)
	_ State = (*CandidateState)(nil)
)

type LeaderState struct {
}

func (l LeaderState) ProcessTimeout() (State, []Command) {
	// TODO implement me
	panic("implement me")
}

func (l LeaderState) ProcessVote(v VoteRequest) (State, VoteResponse) {
	// TODO implement me
	panic("implement me")
}

func (l LeaderState) ProcessAppend(a AppendRequest) (State, AppendResponse) {
	// TODO implement me
	panic("implement me")
}

type FollowerState struct {
}

func (f FollowerState) ProcessTimeout() (State, []Command) {
	// TODO implement me
	panic("implement me")
}

func (f FollowerState) ProcessVote(v VoteRequest) (State, VoteResponse) {
	// TODO implement me
	panic("implement me")
}

func (f FollowerState) ProcessAppend(a AppendRequest) (State, AppendResponse) {
	// TODO implement me
	panic("implement me")
}

type CandidateState struct {
}

func (c CandidateState) ProcessTimeout() (State, []Command) {
	// TODO implement me
	panic("implement me")
}

func (c CandidateState) ProcessVote(v VoteRequest) (State, VoteResponse) {
	// TODO implement me
	panic("implement me")
}

func (c CandidateState) ProcessAppend(a AppendRequest) (State, AppendResponse) {
	// TODO implement me
	panic("implement me")
}
