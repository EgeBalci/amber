package amber

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	pe "github.com/EgeBalci/debug/pe"
)

// Analyze returns a Blueprint structure for the given PE file name
func (bp *Blueprint) Analyze() error {

	abs, err := filepath.Abs(bp.FileName)
	if err != nil {
		return err
	}
	// Set absolute file name
	bp.FullFileName = abs

	peFile, err := pe.Open(bp.FullFileName)
	if err != nil {
		return err
	}
	bp.file = peFile

	arch, err := getPEfileArchitecture(bp.FullFileName)
	if err != nil {
		return err
	}
	bp.Architecture = arch

	if bp.CustomStubName != "" {
		arch, err := getPEfileArchitecture(bp.CustomStubName)
		if err != nil {
			return err
		}
		if arch != bp.Architecture {
			return errors.New("custom stub architecture mismatch")
		}

		stub, err := ioutil.ReadFile(bp.CustomStubName)
		if err != nil {
			return err
		}

		bp.CustomStub = stub
	} else {
		bp.CustomStub = nil
	}

	// Fetch OptionalHeader values to blueprint
	switch hdr := (peFile.OptionalHeader).(type) {
	case *pe.OptionalHeader32:
		// cast those back to a uint32 before use in 32bit
		bp.ImageBase = uint64(hdr.ImageBase)
		bp.Subsystem = hdr.Subsystem
		bp.SizeOfImage = hdr.SizeOfImage

		bp.IsDLL = peFile.Characteristics == (peFile.Characteristics | pe.IMAGE_FILE_DLL)
		bp.HasRelocData = hdr.DataDirectory[5].Size != 0x00
		bp.HasBoundedImports = hdr.DataDirectory[11].Size != 0x00
		bp.HasDelayedImports = hdr.DataDirectory[13].Size != 0x00
		bp.IsCLR = hdr.DataDirectory[14].Size != 0x00

		bp.ExportTable = uint64(hdr.DataDirectory[0].VirtualAddress + uint32(hdr.ImageBase))
		bp.ImportTable = uint64(hdr.DataDirectory[1].VirtualAddress + uint32(hdr.ImageBase))
		bp.RelocTable = uint64(hdr.DataDirectory[5].VirtualAddress + uint32(hdr.ImageBase))
		bp.ImportAdressTable = uint64(hdr.DataDirectory[12].VirtualAddress + uint32(hdr.ImageBase))
		break
	case *pe.OptionalHeader64:
		bp.ImageBase = hdr.ImageBase
		bp.Subsystem = hdr.Subsystem
		bp.SizeOfImage = hdr.SizeOfImage

		bp.IsDLL = peFile.Characteristics == (peFile.Characteristics | pe.IMAGE_FILE_DLL)
		bp.HasRelocData = hdr.DataDirectory[5].Size != 0x00
		bp.HasBoundedImports = hdr.DataDirectory[11].Size != 0x00
		bp.HasDelayedImports = hdr.DataDirectory[13].Size != 0x00
		bp.IsCLR = hdr.DataDirectory[14].Size != 0x00

		bp.ExportTable = uint64(hdr.DataDirectory[0].VirtualAddress + uint32(hdr.ImageBase))
		bp.ImportTable = uint64(hdr.DataDirectory[1].VirtualAddress + uint32(hdr.ImageBase))
		bp.RelocTable = uint64(hdr.DataDirectory[5].VirtualAddress + uint32(hdr.ImageBase))
		bp.ImportAdressTable = uint64(hdr.DataDirectory[12].VirtualAddress + uint32(hdr.ImageBase))
		break
	}

	// Check for TLS callbacks !!!

	fileSize, err := GetFileSize(bp.FullFileName)
	if err != nil {
		return err
	}
	bp.FileSize = fileSize

	return nil
}

func getPEfileArchitecture(fileName string) (int, error) {
	file, err := pe.Open(fileName)
	if err != nil {
		return 0, err
	}

	switch file.FileHeader.Machine {
	case pe.IMAGE_FILE_MACHINE_I386:
		return 32, nil
	case pe.IMAGE_FILE_MACHINE_AMD64:
		return 64, nil
	default:
		return 0, errors.New("unsupported PE file architecture")
	}
}
