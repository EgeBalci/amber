package amber

import (
	"bytes"
	"encoding/binary"
	"errors"
	"path/filepath"

	"github.com/EgeBalci/amber/utils"
	pe "github.com/EgeBalci/debug/pe"
)

const (
	PE_DOS_STUB = "This program cannot be run in DOS mode"
)

var (
	ErrUnsupportedArch  = errors.New("unsupported PE file architecture")
	ErrInvalidPeSpecs   = errors.New("unsupported PE file specs")
	ErrInvalidPeHeaders = errors.New("invalid PE headers")
)

// Blueprint structure contains PE specs, tool parameters and
// OS spesific info
type PE struct {
	Name            string
	FullName        string
	FileSize        int
	IAT             bool
	Resource        bool
	IgnoreIntegrity bool
	IatResolver     bool
	SyscallLoader   bool
	ScrapeHeaders   bool
	// PE specs...
	Architecture      int
	SizeOfImage       uint32
	ImageBase         uint64
	AddressOfEntry    uint32
	Subsystem         uint16
	ImportTable       uint64
	ExportTable       uint64
	RelocTable        uint64
	ImportAdressTable uint64
	HasBoundedImports bool
	HasDelayedImports bool
	HasTLSCallbacks   bool
	HasRelocData      bool
	IsCLR             bool
	IsDLL             bool

	// PE File
	file *pe.File
}

func Open(fileName string) (bp *PE, err error) {
	bp = new(PE)
	bp.Name = fileName
	bp.FullName, err = filepath.Abs(fileName)
	if err != nil {
		return
	}

	bp.file, err = pe.Open(bp.FullName)
	if err != nil {
		return
	}

	switch bp.file.FileHeader.Machine {
	case pe.IMAGE_FILE_MACHINE_I386:
		bp.Architecture = 32
	case pe.IMAGE_FILE_MACHINE_AMD64:
		bp.Architecture = 64
	default:
		return nil, ErrUnsupportedArch
	}

	// Fetch OptionalHeader values to blueprint
	switch hdr := (bp.file.OptionalHeader).(type) {
	case *pe.OptionalHeader32:
		// cast those back to a uint32 before use in 32bit
		bp.ImageBase = uint64(hdr.ImageBase)
		bp.Subsystem = hdr.Subsystem
		bp.SizeOfImage = hdr.SizeOfImage

		bp.IsDLL = bp.file.Characteristics == (bp.file.Characteristics | pe.IMAGE_FILE_DLL)
		bp.HasRelocData = hdr.DataDirectory[5].Size != 0x00
		bp.HasBoundedImports = hdr.DataDirectory[11].Size != 0x00
		bp.HasDelayedImports = hdr.DataDirectory[13].Size != 0x00
		bp.IsCLR = hdr.DataDirectory[14].Size != 0x00

		bp.ExportTable = uint64(hdr.DataDirectory[0].VirtualAddress + uint32(hdr.ImageBase))
		bp.ImportTable = uint64(hdr.DataDirectory[1].VirtualAddress + uint32(hdr.ImageBase))
		bp.RelocTable = uint64(hdr.DataDirectory[5].VirtualAddress + uint32(hdr.ImageBase))
		bp.ImportAdressTable = uint64(hdr.DataDirectory[12].VirtualAddress + uint32(hdr.ImageBase))
	case *pe.OptionalHeader64:
		bp.ImageBase = hdr.ImageBase
		bp.Subsystem = hdr.Subsystem
		bp.SizeOfImage = hdr.SizeOfImage

		bp.IsDLL = bp.file.Characteristics == (bp.file.Characteristics | pe.IMAGE_FILE_DLL)
		bp.HasRelocData = hdr.DataDirectory[5].Size != 0x00
		bp.HasBoundedImports = hdr.DataDirectory[11].Size != 0x00
		bp.HasDelayedImports = hdr.DataDirectory[13].Size != 0x00
		bp.IsCLR = hdr.DataDirectory[14].Size != 0x00

		bp.ExportTable = uint64(hdr.DataDirectory[0].VirtualAddress + uint32(hdr.ImageBase))
		bp.ImportTable = uint64(hdr.DataDirectory[1].VirtualAddress + uint32(hdr.ImageBase))
		bp.RelocTable = uint64(hdr.DataDirectory[5].VirtualAddress + uint32(hdr.ImageBase))
		bp.ImportAdressTable = uint64(hdr.DataDirectory[12].VirtualAddress + uint32(hdr.ImageBase))
	}

	bp.FileSize, err = utils.GetFileSize(bp.FullName)
	return
}

// AssemblePayload generates the binary stub bla bla...
func (pe *PE) AssembleLoader() ([]byte, error) {

	var (
		rawFile = pe.file.RawBytes
		err     error
	)

	if pe.ScrapeHeaders {
		rawFile, err = pe.ScrapePeHeaders()
		if err != nil {
			return nil, err
		}
	}

	// Add a call over the given binary
	payload, err := pe.AddCallOver(rawFile)
	if err != nil {
		return nil, err
	}

	// Decide on the architecture, API block, and loader types...
	// we have 3 pre-assembled loaders for public version of amber.
	switch pe.Architecture {
	case 32:
		if pe.SyscallLoader {
			return nil, errors.New("syscall loader only supports 64 bit PE files")
		}
		payload = append(payload, LOADER_32...)
	case 64:
		if pe.SyscallLoader {
			payload = append(payload, SYSCALL_LOADER_64...)
		} else {
			payload = append(payload, LOADER_64...)
		}

	default:
		return nil, ErrUnsupportedArch
	}

	if pe.IatResolver {
		if pe.SyscallLoader {
			return nil, errors.New("cannot use IAT resolver with syscall loader")
		}
		switch pe.Architecture {
		case 32:
			payload = bytes.ReplaceAll(payload, CRC_API_32, IAT_API_32)
		case 64:
			payload = bytes.ReplaceAll(payload, CRC_API_64, IAT_API_64)
		}
	}

	return payload, nil
}

// AddCallOver function adds a call instruction at the beginning of the given payload
// address of the payload will be pushed to the stack and execution will continue after the end of payload
func (pe *PE) AddCallOver(payload []byte) ([]byte, error) {
	// // Perform a short call over the payload
	size := uint32(len(payload))
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, size)
	if err != nil {
		return nil, err
	}
	return append(append([]byte{0xe8}, buf.Bytes()...), payload...), nil
}

func (pe *PE) ScrapePeHeaders() ([]byte, error) {
	rawFile, err := pe.file.Bytes()
	if err != nil {
		return nil, err
	}

	// Scrape MZ magic bytes...
	if rawFile[0] == 'M' &&
		rawFile[1] == 'Z' {
		rawFile[0] = 0x00
		rawFile[1] = 0x00
	} else {
		return nil, ErrInvalidPeHeaders
	}

	// Scrape the DOS stub message...
	if bytes.Contains(rawFile, []byte(PE_DOS_STUB)) {
		return nil, ErrInvalidPeHeaders
	}

	return bytes.Replace(rawFile, []byte(PE_DOS_STUB), make([]byte, len(PE_DOS_STUB)), 1), nil
}
