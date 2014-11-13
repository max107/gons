package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"syscall"
)

const (
	NsRunDir  = "/var/run/netns"
	SelfNetNs = "/proc/self/ns/net"
)

func Create(name string) {
	netNsPath := path.Join(NsRunDir, name)
	os.Mkdir(NsRunDir, 0755)

	if err := syscall.Mount(NsRunDir, NsRunDir, "none", syscall.MS_BIND, ""); err != nil {
		log.Fatalf("Could not create Network namespace: %s", err)
	}
}

func Mount(name string) {
	netNsPath := path.Join(NsRunDir, name)
	if err := syscall.Mount(SelfNetNs, netNsPath, "none", syscall.MS_BIND, ""); err != nil {
		log.Fatalf("Could not Mount Network namespace: %s", err)
	}
}

func Unshare() {
	if err := syscall.Unshare(syscall.CLONE_NEWNET); err != nil {
		log.Fatalf("Could not clone new Network namespace: %s", err)
	}
}

func Unlink(name string) {
	netNsPath := path.Join(NsRunDir, name)
	if err := syscall.Unlink(netNsPath); err != nil {
		log.Fatalf("Could not Unlink new Network namespace: %s", err)
	}
}

func Unmount(name string) {
	netNsPath := path.Join(NsRunDir, name)
	if err := syscall.Unmount(netNsPath, syscall.MNT_DETACH); err != nil {
		log.Fatalf("Could not Unmount new Network namespace: %s", err)
	}
}

func List() {
	ifcs, _ := net.Interfaces()
	for _, ifc := range ifcs {
		fmt.Printf("%#v\n", ifc)
	}
}

func Open(name string) {
	netNsPath := path.Join(NsRunDir, name)
	fd, err := syscall.Open(netNsPath, syscall.O_RDONLY|syscall.O_CREAT|syscall.O_EXCL, 0)
	if err != nil {
		log.Fatalf("Could not create Network namespace: %s", err)
	}
	syscall.Close(fd)
}
