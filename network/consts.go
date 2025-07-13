package network

type NetworkStatus int

const (
	Unknown NetworkStatus = -1
	Down    NetworkStatus = 0
	Up      NetworkStatus = 1
)

const (
	TIMEOUT_SEC = -1
)

type Nlmsghdr struct {
	NlmsgLen   uint32 // Size of message including header
	NlmsgType  uint16 // Type of message content
	NlmsgFlags uint16 // Additional flags
	NlmsgSeq   uint32 // Sequence number
	NlmsgPid   uint32 // Sender port ID
}

const byteSize = 8
const (
	nlmsgLenSize   = 32 / byteSize
	nlmsgTypeSize  = 16 / byteSize
	nlmsgFlagsSize = 16 / byteSize
	nlmsgSeqSize   = 32 / byteSize
	nlmsgPidSize   = 32 / byteSize
)

const nlmsgSize = nlmsgLenSize + nlmsgTypeSize + nlmsgFlagsSize + nlmsgSeqSize + nlmsgPidSize

type nlmsgItemStart uint

const (
	nlmsgLenStart   = 0
	nlmsgTypeStart  = nlmsgLenStart + nlmsgLenSize
	nlmsgFlagsStart = nlmsgTypeStart + nlmsgTypeSize
	nlmsgSeqStart   = nlmsgFlagsStart + nlmsgFlagsSize
	nlmsgPidStart   = nlmsgSeqStart + nlmsgSeqSize
)
