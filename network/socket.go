package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	syscall "golang.org/x/sys/unix"
)

type Iface struct {
	Name   string
	Status NetworkStatus
}

type Ifaces = []*Iface

func NewIface(name string) (*Iface, error) {
	status, err := getIfaceStatus(name)
	if err != nil {
		return nil, err
	}

	port := Iface{
		Name:   name,
		Status: status,
	}
	return &port, nil

}

func NewIfaces() (Ifaces, error) {
	ifaces, err := ListPhisicalIfaces("/sys/devices")
	if err != nil {
		return nil, err
	}

	ports := Ifaces{}
	for _, iface := range ifaces {
		iface, err := NewIface(iface)
		if err != nil {
			return nil, err
		}
		ports = append(ports, iface)
	}

	return ports, nil

}

func ListPhisicalIfaces(root string) ([]string, error) {
	findList := []string{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() != "net" || strings.Contains(path, "virtual") {
			return nil
		}

		devices, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, dev := range devices {
			findList = append(findList, dev.Name())
		}

		return nil

	})
	sort.Slice(findList, func(i, j int) bool { return findList[i] < findList[j] })
	return findList, err
}

func getIfaceStatus(name string) (NetworkStatus, error) {
	file, err := os.Open(filepath.Join("/sys/class/net", name, "carrier"))
	if err != nil {
		return 0, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	switch string(content) {
	case "0":
		return Down, nil
	case "1":
		return Up, nil
	default:
		return 0, errors.New(fmt.Sprintf("Couldn't get the status of network interface \"%s\"", name))
	}
}

func NewNetlinkSocket() error {
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

	pfds := []syscall.PollFd{syscall.PollFd{
		Fd:      int32(fd),
		Events:  syscall.POLLIN,
		Revents: 0,
	}}

	_, err = syscall.Poll(pfds, TIMEOUT_SEC)
	if err != nil {
		return err
	}

	buff := make([]byte, 4096)
	oob := make([]byte, 0)

	n, _, _, _, err := syscall.Recvmsg(fd, buff, oob, 0)
	if err != nil {
		return err
	}

	for offset := 0; offset+nlmsgSize <= n; {
		nh := parseNetlinkMsg(buff, offset)

		end := offset + int(nh.NlmsgLen)
		fmt.Printf("from %d to %d of %d bytes\n", offset, end, n)
		if nh.NlmsgType == syscall.NLMSG_DONE {
			break
		}

		if offset == 0 {
			fmt.Printf("%v#\n", nh)
		}
		fmt.Printf("%X\n", buff[offset+nlmsgSize:end])

		offset += int(nh.NlmsgLen)

	}

	return nil
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
