package amber

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	pe "github.com/EgeBalci/debug/pe"
)

// CompileStub generates the final stub file with given payload
func (bp *Blueprint) CompileStub(payload []byte) ([]byte, error) {

	tmpStub, err := ioutil.TempFile(os.TempDir(), "amber_")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpStub.Name())
	stub := []byte{}

	if bp.CustomStub != nil {
		stub = bp.CustomStub
	} else {
		switch bp.Architecture {
		case 32:
			stub, err = base64.StdEncoding.DecodeString(STUB32)
		case 64:
			stub, err = base64.StdEncoding.DecodeString(STUB64)
		}
	}

	if bp.IAT && !bp.checkRequiredIATFuncs(stub) {
		return nil, errors.New("selected stub does not support IAT resolver API")
	}

	tmpStub.Write(stub)
	tmpStubPE, err := pe.Open(tmpStub.Name())
	if err != nil {
		return nil, err
	}

	var sizeOfImage, fileAlignment uint32
	lastSection := tmpStubPE.Sections[tmpStubPE.NumberOfSections-1]

	switch tmpStubPE.FileHeader.Machine {
	case pe.IMAGE_FILE_MACHINE_I386:
		hdr := tmpStubPE.OptionalHeader.(*pe.OptionalHeader32)
		sizeOfImage = hdr.SizeOfImage
		fileAlignment = hdr.FileAlignment
		break
	case pe.IMAGE_FILE_MACHINE_AMD64:
		hdr := tmpStubPE.OptionalHeader.(*pe.OptionalHeader64)
		sizeOfImage = hdr.SizeOfImage
		fileAlignment = hdr.FileAlignment
		break
	}

	// Edit the pre-compiled EXE stub header values by replacing raw bytes

	// Change Subsystem
	stub, err = setSubsystem(stub, bp.Subsystem)
	if err != nil {
		return nil, err
	}

	// Change SizeOfImage
	stub, err = setSizeOfImage(stub, bp.SizeOfImage+sizeOfImage)
	if err != nil {
		return nil, err
	}

	// Randomize last section name
	lastSectionOffset := bytes.Index(stub, []byte(lastSection.Name))
	newSectionName := "." + randomString(7) // Randomize the section name
	stub = append(stub[:lastSectionOffset], bytes.Replace(stub[lastSectionOffset:], []byte(lastSection.Name), []byte(newSectionName), 1)...)
	oldBytes := make([]byte, 4)
	newBytes := make([]byte, 4)

	// Change SectionVirtualSize
	binary.LittleEndian.PutUint32(oldBytes, uint32(lastSection.VirtualSize))
	binary.LittleEndian.PutUint32(newBytes, uint32(len(payload)))
	stub = append(stub[:lastSectionOffset], bytes.Replace(stub[lastSectionOffset:], oldBytes, newBytes, 1)...)
	// Change SectionRawSize
	binary.LittleEndian.PutUint32(oldBytes, uint32(lastSection.Size))
	binary.LittleEndian.PutUint32(newBytes, align(uint32(len(payload)), fileAlignment))
	stub = append(stub[:lastSectionOffset], bytes.Replace(stub[lastSectionOffset:], oldBytes, newBytes, 1)...)

	// Align section size
	eofTamper := make([]byte, align(uint32(len(payload)), fileAlignment)-uint32(len(payload)))
	for i := range eofTamper {
		eofTamper[i] = 0x00
	}
	payload = append(payload, eofTamper...)

	// Replace the section data with the reflective payload
	oldSectionData, err := lastSection.Data()
	if err != nil {
		return nil, err
	}
	stub = bytes.Replace(stub, oldSectionData, payload, 1)

	return stub, nil
}

func align(size, align uint32) uint32 {
	if 0 == (size % align) {
		return size
	}
	return size + (align - (size % align))
}

func (bp *Blueprint) checkRequiredIATFuncs(stub []byte) bool {

	if !strings.Contains(string(stub), "LoadLibrary") || !strings.Contains(string(stub), "GetProcAddress") {
		return false
	}

	if bp.HasRelocData {
		if !strings.Contains(string(stub), "VirtualAlloc") {
			return false
		}
	} else {
		if !strings.Contains(string(stub), "VirtualProtect") {
			return false
		}
	}

	if bp.Architecture == 64 {
		if !strings.Contains(string(stub), "FlushInstructionCache") {
			return false
		}
	}

	return true
}

// Change the SubSystem header value of given raw PE file with newSubSystem
func setSubsystem(peFile []byte, newSubSystem uint16) ([]byte, error) {
	ntHeaderOffset, err := getNtHeaderOffset(peFile)
	if err != nil {
		return nil, err
	}
	subSystemBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(subSystemBytes, newSubSystem)
	return append(append(peFile[:ntHeaderOffset+0x5c], subSystemBytes...), peFile[ntHeaderOffset+0x5e:]...), nil
}

// Change the SizeOfImage header value of given raw PE file with newSizeOfImage
func setSizeOfImage(peFile []byte, newSizeOfImage uint32) ([]byte, error) {
	ntHeaderOffset, err := getNtHeaderOffset(peFile)
	if err != nil {
		return nil, err
	}
	sizeOfImageBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sizeOfImageBytes, newSizeOfImage)
	return append(append(peFile[:ntHeaderOffset+0x50], sizeOfImageBytes...), peFile[ntHeaderOffset+0x54:]...), nil
}

// Change the ImageBase header value of given raw PE file with newImageBase
func setImageBase(peFile []byte, newImageBase uint64) ([]byte, error) {
	r := bytes.NewReader(peFile)
	f := new(pe.File)
	sr := io.NewSectionReader(r, 0, 1<<63-1)

	ntHeaderOffset, err := getNtHeaderOffset(peFile)
	if err != nil {
		return nil, err
	}

	sr.Seek(int64(ntHeaderOffset+4), 0)
	if err := binary.Read(sr, binary.LittleEndian, &f.FileHeader); err != nil {
		return nil, err
	}
	switch f.FileHeader.Machine {
	case pe.IMAGE_FILE_MACHINE_I386:
		imageBaseBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(imageBaseBytes, uint32(newImageBase))
		return append(append(peFile[:ntHeaderOffset+0x30], imageBaseBytes...), peFile[ntHeaderOffset+0x34:]...), nil
	case pe.IMAGE_FILE_MACHINE_AMD64:
		imageBaseBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(imageBaseBytes, newImageBase)
		return append(append(peFile[:ntHeaderOffset+0x30], imageBaseBytes...), peFile[ntHeaderOffset+0x38:]...), nil
	default:
		return nil, fmt.Errorf("unsupported COFF file header machine value of 0x%x", f.FileHeader.Machine)
	}
}

// Get the NT_IMAGE_HEADER offset of the given raw PE file
func getNtHeaderOffset(peFile []byte) (uint32, error) {
	ntHeaderOffset := binary.LittleEndian.Uint32(peFile[0x3c:0x40])
	if !(peFile[ntHeaderOffset] == 'P' && peFile[ntHeaderOffset+1] == 'E' && peFile[ntHeaderOffset+2] == 0 && peFile[ntHeaderOffset+3] == 0) {
		return 0, fmt.Errorf("invalid PE COFF file signature of %v", ntHeaderOffset)
	}
	return ntHeaderOffset, nil
}
