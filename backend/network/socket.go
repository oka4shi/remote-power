package network

import (
	"encoding/binary"
	"errors"
	"log"

	syscall "golang.org/x/sys/unix"
)

func WatchNetlinkSocket(c chan []byte) error {
	fd, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, syscall.NETLINK_ROUTE)
	if err != nil {
		return err
	}
	if fd < 0 {
		return errors.New("Cannot get a socket")
	}

	sa := &syscall.SockaddrNetlink{
		Family: syscall.AF_NETLINK,
		Pad:    0,
		Pid:    0, // set 0 to be assigned by the kernel
		Groups: syscall.RTMGRP_LINK,
	}
	err = syscall.Bind(fd, sa)
	if err != nil {
		return err
	}

	pfds := []syscall.PollFd{{
		Fd:      int32(fd),
		Events:  syscall.POLLIN,
		Revents: 0,
	}}

	_, err = syscall.Poll(pfds, TIMEOUT_SEC)
	if err != nil {
		return err
	}

	for {
		buff := make([]byte, 4096)
		oob := make([]byte, 0)

		n, _, _, _, err := syscall.Recvmsg(fd, buff, oob, 0)
		log.Printf("Got a message from a netlink socket:")
		if err != nil {
			return err
		}

		for offset := 0; offset+nlmsgSize <= n; {
			nh := parseNetlinkMsg(buff, offset)

			end := offset + int(nh.NlmsgLen)
			log.Printf("Reading %d – %d (total: %d bytes)\n", offset, end, n)
			log.Printf("nlmsgLen: %d, nlmsgType: 0x%X, nlmsgFlags: 0x%X, nlmsgSeq: %d, nlmsgPid: %d\n",
				nh.NlmsgLen, nh.NlmsgType, nh.NlmsgFlags, nh.NlmsgSeq, nh.NlmsgPid)

			if nh.NlmsgType == syscall.NLMSG_DONE {
				break
			}

			payload := buff[offset+nlmsgSize : end]

			if nh.NlmsgType == syscall.NLMSG_MIN_TYPE {
				log.Printf("%X\n", payload)
				c <- payload
			}

			offset += int(nh.NlmsgLen)
		}
	}
}

func parseNetlinkMsg(b []byte, offset int) Nlmsghdr {
	nlmsghdr := Nlmsghdr{
		NlmsgLen:   bytesToUint32(b, offset, nlmsgLenStart, nlmsgLenSize),
		NlmsgType:  bytesToUint16(b, offset, nlmsgTypeStart, nlmsgTypeSize),
		NlmsgFlags: bytesToUint16(b, offset, nlmsgFlagsStart, nlmsgFlagsSize),
		NlmsgSeq:   bytesToUint32(b, offset, nlmsgSeqStart, nlmsgSeqSize),
		NlmsgPid:   bytesToUint32(b, offset, nlmsgPidStart, nlmsgPidSize),
	}
	return nlmsghdr
}

func bytesToUint16(b []byte, offset, start, size int) uint16 {
	return binary.LittleEndian.Uint16(b[offset+start : offset+start+size])
}

func bytesToUint32(b []byte, offset, start, size int) uint32 {
	return binary.LittleEndian.Uint32(b[offset+start : offset+start+size])
}
