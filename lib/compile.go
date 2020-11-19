package amber

import (
	pe "amber/debug/pe"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"io/ioutil"
	"os"
)

// CompileStub generates the final stub file with given payload
func (bp *Blueprint) CompileStub(payload []byte) ([]byte, error) {

	tmpStub, err := ioutil.TempFile(os.TempDir(), "amber_")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpStub.Name())

	stub := []byte{}

	switch bp.Architecture {
	case 32:
		stub, err = base64.StdEncoding.DecodeString(STUB32)
	case 64:
		stub, err = base64.StdEncoding.DecodeString(STUB64)
	}

	tmpStub.Write(stub)
	tmpStubPE, err := pe.Open(tmpStub.Name())

	var sizeOfImage, fileAlignment uint32
	var imageBase uint64

	lastSection := tmpStubPE.Sections[tmpStubPE.NumberOfSections-1]

	switch tmpStubPE.FileHeader.Machine {
	case pe.IMAGE_FILE_MACHINE_I386:
		// cast those back to a uint32 before use in 32bit
		hdr := tmpStubPE.OptionalHeader.(*pe.OptionalHeader32)
		imageBase = uint64(hdr.ImageBase) // cast this back to a uint32 before use in 32bit
		sizeOfImage = hdr.SizeOfImage
		fileAlignment = hdr.FileAlignment
		break
	case pe.IMAGE_FILE_MACHINE_AMD64:
		hdr := tmpStubPE.OptionalHeader.(*pe.OptionalHeader64)
		imageBase = hdr.ImageBase
		sizeOfImage = hdr.SizeOfImage
		fileAlignment = hdr.FileAlignment
		break
	}

	lastSectionOffset := bytes.Index(stub, []byte(lastSection.Name))
	newSectionName := "." + randomString(7)
	stub = append(stub[:lastSectionOffset], bytes.Replace(stub[lastSectionOffset:], []byte(lastSection.Name), []byte(newSectionName), 1)...)
	oldBytes := make([]byte, 4)
	newBytes := make([]byte, 4)

	// Change ImageBase
	binary.LittleEndian.PutUint32(oldBytes, uint32(imageBase))
	binary.LittleEndian.PutUint32(newBytes, uint32(bp.ImageBase))
	stub = bytes.Replace(stub, oldBytes, newBytes, 1)
	// Change SizeOfImage
	binary.LittleEndian.PutUint32(oldBytes, uint32(sizeOfImage))
	binary.LittleEndian.PutUint32(newBytes, uint32(bp.SizeOfImage+sizeOfImage))
	stub = bytes.Replace(stub, oldBytes, newBytes, 1)
	// Change SectionVirtualSize
	binary.LittleEndian.PutUint32(oldBytes, uint32(lastSection.VirtualSize))
	binary.LittleEndian.PutUint32(newBytes, align(uint32(len(payload)), fileAlignment, 0))
	stub = append(stub[:lastSectionOffset], bytes.Replace(stub[lastSectionOffset:], oldBytes, newBytes, 1)...)
	// Change SectionRawSize
	binary.LittleEndian.PutUint32(oldBytes, uint32(lastSection.Size))
	binary.LittleEndian.PutUint32(newBytes, align(uint32(len(payload)), fileAlignment, 0))
	stub = append(stub[:lastSectionOffset], bytes.Replace(stub[lastSectionOffset:], oldBytes, newBytes, 3)...)

	oldSectionData, err := lastSection.Data()
	if err != nil {
		return nil, err
	}
	stub = bytes.Replace(stub, oldSectionData, payload, 1)

	return stub, nil
}

func align(size, align, addr uint32) uint32 {
	if 0 == (size % align) {
		return addr + size
	}
	return addr + (size/align+1)*align
}
