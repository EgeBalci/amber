package amber

import (
	"errors"

	"github.com/EgeBalci/debug/pe"
)

// Analyze returns a Blueprint structure for the given PE file name
func Analyze(fileName string) (*Blueprint, error) {
	var new Blueprint

	peFile, err := pe.Open(fileName)
	if err != nil {
		return nil, err
	}
	new.file = peFile
	// Set file name
	new.FileName = fileName
	// Set architecture
	switch peFile.FileHeader.Machine {
	case IMAGE_FILE_MACHINE_I386:
		new.Architecture = 32
	case IMAGE_FILE_MACHINE_AMD64:
		new.Architecture = 64
	default:
		return nil, errors.New("file architecture not supported")
	}

	switch hdr := (peFile.OptionalHeader).(type) {
	case *pe.OptionalHeader32:
		// cast those back to a uint32 before use in 32bit
		new.ImageBase = uint64(hdr.ImageBase)
		new.Subsystem = hdr.Subsystem
		new.SizeOfImage = hdr.SizeOfImage

		new.IsDLL = peFile.Characteristics == (peFile.Characteristics | IMAGE_FILE_DLL)
		new.HasRelocData = hdr.DataDirectory[5].Size != 0x00
		new.HasBoundedImports = hdr.DataDirectory[11].Size != 0x00
		new.HasDelayedImports = hdr.DataDirectory[13].Size != 0x00
		new.IsCLR = hdr.DataDirectory[14].Size != 0x00

		new.ExportTable = uint64(hdr.DataDirectory[0].VirtualAddress + uint32(hdr.ImageBase))
		new.ImportTable = uint64(hdr.DataDirectory[1].VirtualAddress + uint32(hdr.ImageBase))
		new.RelocTable = uint64(hdr.DataDirectory[5].VirtualAddress + uint32(hdr.ImageBase))
		new.ImportAdressTable = uint64(hdr.DataDirectory[12].VirtualAddress + uint32(hdr.ImageBase))
		break
	case *pe.OptionalHeader64:
		new.ImageBase = hdr.ImageBase
		new.Subsystem = hdr.Subsystem
		new.SizeOfImage = hdr.SizeOfImage

		new.IsDLL = peFile.Characteristics == (peFile.Characteristics | IMAGE_FILE_DLL)
		new.HasRelocData = hdr.DataDirectory[5].Size != 0x00
		new.HasBoundedImports = hdr.DataDirectory[11].Size != 0x00
		new.HasDelayedImports = hdr.DataDirectory[13].Size != 0x00
		new.IsCLR = hdr.DataDirectory[14].Size != 0x00

		new.ExportTable = uint64(hdr.DataDirectory[0].VirtualAddress + uint32(hdr.ImageBase))
		new.ImportTable = uint64(hdr.DataDirectory[1].VirtualAddress + uint32(hdr.ImageBase))
		new.RelocTable = uint64(hdr.DataDirectory[5].VirtualAddress + uint32(hdr.ImageBase))
		new.ImportAdressTable = uint64(hdr.DataDirectory[12].VirtualAddress + uint32(hdr.ImageBase))
		break
	}

	fileSize, err := GetFileSize(fileName)
	if err != nil {
		return nil, err
	}
	new.FileSize = fileSize
	// Check for TLS callbacks !!!

	return &new, nil
}
