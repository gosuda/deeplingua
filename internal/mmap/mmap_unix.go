//go:build unix || linux || darwin || freebsd || openbsd
// +build unix linux darwin freebsd openbsd

package mmap

import (
	"os"
	"syscall"
)

const (
	PROT_READ  = syscall.PROT_READ
	PROT_WRITE = syscall.PROT_WRITE

	MAP_SHARED  = syscall.MAP_SHARED
	MAP_PRIVATE = syscall.MAP_PRIVATE
)

func Map(fd uintptr, offset int, len int, prot int, flags int) ([]byte, error) {
	return syscall.Mmap(int(fd), int64(offset), len, prot, flags)
}

func UnMap(b []byte) error {
	return syscall.Munmap(b)
}

func OpenFile(name string, flag int, perm os.FileMode) (uintptr, error) {
	fd, err := syscall.Open(name, flag, uint32(perm))
	if err != nil {
		return 0, err
	}
	return uintptr(fd), nil
}

func CloseFile(fd uintptr) error {
	return syscall.Close(int(fd))
}
