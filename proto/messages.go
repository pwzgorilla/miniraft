package proto

//Constants for request type
const (
	TypeAppendEntryRequest = iota
	TypeAppendEntryResponse
	TypeVoteRequest
	TypeVoteReply
	TypeTimeout
	TypeClientAppendRequest
	TypeHeartBeat
	TypeHeartBeatResponse
)

const (
	Broadcast = "broadcast"
)

type Lsn uint64      //Log sequence number, unique for all time.
type ErrRedirect int // See Log.Append. Implements Error interface.
type LogEntry interface {
	Lsn() Lsn
	Data() []byte
	IsCommitted() bool
	CurrentTerm() uint64
}

//LogEntry interface implementation
type LogEntryObj struct {
	LogSeqNumber Lsn
	DataBytes    []byte
	Committed    bool
	Term         uint64
}

func (entry LogEntryObj) Lsn() Lsn {
	return entry.LogSeqNumber
}

func (entry LogEntryObj) Data() []byte {
	return entry.DataBytes
}

func (entry LogEntryObj) IsCommitted() bool {
	return entry.Committed
}

func (entry LogEntryObj) CurrentTerm() uint64 {
	return entry.Term
}

type Event struct {
	Type int
	Data interface{}
}

//This struct will be sent as RPC message between the replicas
type AppendEntryRequest struct {
	LogEntries        []LogEntryObj
	LeaderID          string
	LeaderCommitIndex Lsn
	PreviousLogIndex  Lsn
	Term              uint64
	PreviousLogTerm   uint64
}

type AppendEntryResponse struct {
	//Reply strucrure.
	ServerID         string
	Term             uint64
	Success          bool
	PreviousLogIndex Lsn
	ExpectedIndex    Lsn
	RequestedIndex   Lsn
}

type VoteRequest struct {
	Term         uint64
	CandidateID  string
	LastLogIndex Lsn
	LastLogTerm  uint64
}

type VoteReply struct {
	Term     uint64
	Result   bool
	ServerID string
}

type Timeout struct{}

type ClientAppendRequest struct {
	Data       []byte
	ResponseCh *chan string
}

type HeartBeatRequest struct {
	Term     uint64
	LeaderID string
}

type HeartBeatReply struct {
	Term     uint64
	LeaderID string
	Success  bool
}
