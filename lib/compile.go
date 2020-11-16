package amber

import (
	pe "amber/debug/pe"
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
)

// CompileStub generates the final stub file with given payload
func (bp *Blueprint) CompileStub(payload []byte) ([]byte, error) {
	var (
		peFile *pe.File
		err    error
	)

	tmpStub, err := ioutil.TempFile(os.TempDir(), "amber_")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpStub.Name())

	switch bp.Architecture {
	case 32:
		tmpStub.Write(PE_STUB32)
	case 64:
		tmpStub.Write(PE_STUB64)
	}

	tmpStubPE, err := pe.Open(tmpStub.Name())

	var sizeOfImage, fileAlignment uint32
	var imageBase uint64

	lastSection := tmpStubPE.Sections[peFile.NumberOfSections-1]

	switch hdr := (tmpStubPE.OptionalHeader).(type) {
	case *pe.OptionalHeader32:
		imageBase = uint64(hdr.ImageBase) // cast this back to a uint32 before use in 32bit
		sizeOfImage = hdr.SizeOfImage
		fileAlignment = hdr.FileAlignment
	case *pe.OptionalHeader64:
		imageBase = hdr.ImageBase
		sizeOfImage = hdr.SizeOfImage
		fileAlignment = hdr.FileAlignment
	}

	newSectionName := "." + randomString(7)
	final := bytes.Replace(PE_STUB32, []byte(lastSection.Name), []byte(newSectionName), 1)
	oldBytes := make([]byte, 4)
	newBytes := make([]byte, 4)

	switch bp.Architecture {
	case 32:
		// Change ImageBase
		binary.LittleEndian.PutUint32(oldBytes, uint32(imageBase))
		binary.LittleEndian.PutUint32(newBytes, uint32(bp.ImageBase))
		final = bytes.Replace(PE_STUB32, oldBytes, newBytes, 1)
		// Change SizeOfImage
		binary.LittleEndian.PutUint32(oldBytes, uint32(sizeOfImage))
		binary.LittleEndian.PutUint32(newBytes, uint32(bp.SizeOfImage))
		final = bytes.Replace(PE_STUB32, oldBytes, newBytes, 1)
		// Change SectionSize
		binary.LittleEndian.PutUint32(oldBytes, uint32(lastSection.Size))
		binary.LittleEndian.PutUint32(newBytes, uint32(align(uint32(len(payload)), fileAlignment, 0)))
		final = bytes.Replace(PE_STUB32, oldBytes, newBytes, 1)
		// Change SectionVirtualSize
		binary.LittleEndian.PutUint32(oldBytes, uint32(lastSection.VirtualSize))
		binary.LittleEndian.PutUint32(newBytes, uint32(len(payload)))
		final = bytes.Replace(PE_STUB32, oldBytes, newBytes, 1)
	case 64:
		// Change ImageBase
		binary.LittleEndian.PutUint32(oldBytes, uint32(imageBase))
		binary.LittleEndian.PutUint32(newBytes, uint32(bp.ImageBase))
		final = bytes.Replace(PE_STUB32, oldBytes, newBytes, 1)
		// Change SizeOfImage
		binary.LittleEndian.PutUint32(oldBytes, uint32(sizeOfImage))
		binary.LittleEndian.PutUint32(newBytes, uint32(bp.SizeOfImage))
		final = bytes.Replace(PE_STUB32, oldBytes, newBytes, 1)
		// Change SectionSize
		binary.LittleEndian.PutUint32(oldBytes, uint32(lastSection.Size))
		binary.LittleEndian.PutUint32(newBytes, uint32(align(uint32(len(payload)), fileAlignment, 0)))
		final = bytes.Replace(PE_STUB32, oldBytes, newBytes, 1)
		// Change SectionVirtualSize
		binary.LittleEndian.PutUint32(oldBytes, uint32(lastSection.VirtualSize))
		binary.LittleEndian.PutUint32(newBytes, uint32(len(payload)))
		final = bytes.Replace(PE_STUB32, oldBytes, newBytes, 1)
	}

	return final, nil
}

/*
// func setSizeOfImage(fileName string, newSizeOfImage uint64) {
// 	verbose("Aligning SizeOfImage...", "*")
// 	rawFile, err := ioutil.ReadFile(fileName)
// 	parseErr(err)
// 	oFileHeader := int(rawFile[0x3c]) // Offset to EXE header
// 	if target.arch == "x64" {
// 		oSizeOfImage := int(oFileHeader + 0x50) // Offset to SizeOfImage
// 		SizeOfImage := binary.LittleEndian.Uint64(rawFile[oSizeOfImage : oSizeOfImage+8])
// 		verbose("SizeOfImage "+fmt.Sprintf("0x%016X -> ", SizeOfImage)+fmt.Sprintf("0x%016X", newSizeOfImage), "*")
// 		binary.LittleEndian.PutUint64(rawFile[oSizeOfImage:oSizeOfImage+8], newSizeOfImage)
// 	} else {
// 		oSizeOfImage := int(oFileHeader + 0x50) // Offset to SizeOfImage
// 		SizeOfImage := binary.LittleEndian.Uint32(rawFile[oSizeOfImage : oSizeOfImage+8])
// 		verbose("SizeOfImage "+fmt.Sprintf("0x%08X -> ", SizeOfImage)+fmt.Sprintf("0x%08X", newSizeOfImage), "*")
// 		binary.LittleEndian.PutUint32(rawFile[oSizeOfImage:oSizeOfImage+4], uint32(newSizeOfImage))
// 	}
// 	err = ioutil.WriteFile(target.FileName, rawFile, os.ModeDir)
// 	parseErr(err)
// }

// func setSizeOfInitializedData(fileName string, newSizeOfInitializedData uint64) {
// 	verbose("Aligning SizeOfImage...", "*")
// 	rawFile, err := ioutil.ReadFile(fileName)
// 	parseErr(err)
// 	oFileHeader := int(rawFile[0x3c]) // Offset to EXE header
// 	if target.arch == "x64" {
// 		oSizeOfInitializedData := int(oFileHeader + 0x20) // Offset to SizeOfInitializedData
// 		SizeOfInitializedData := binary.LittleEndian.Uint64(rawFile[oSizeOfInitializedData : oSizeOfInitializedData+8])
// 		verbose("SizeOfInitializedData "+fmt.Sprintf("0x%016X -> ", SizeOfInitializedData)+fmt.Sprintf("0x%016X", newSizeOfInitializedData), "*")
// 		binary.LittleEndian.PutUint64(rawFile[oSizeOfInitializedData:oSizeOfInitializedData+8], newSizeOfInitializedData)
// 	} else {
// 		oSizeOfInitializedData := int(oFileHeader + 0x20) // Offset to SizeOfImage
// 		SizeOfInitializedData := binary.LittleEndian.Uint32(rawFile[oSizeOfInitializedData : oSizeOfInitializedData+8])
// 		verbose("SizeOfInitializedData "+fmt.Sprintf("0x%08X -> ", SizeOfInitializedData)+fmt.Sprintf("0x%08X", newSizeOfInitializedData), "*")
// 		binary.LittleEndian.PutUint32(rawFile[oSizeOfInitializedData:oSizeOfInitializedData+4], uint32(newSizeOfInitializedData))
// 	}
// 	err = ioutil.WriteFile(target.FileName, rawFile, os.ModeDir)
// 	parseErr(err)
// }

// func setImageBase(fileName string, newImageBase uint64) {
// 	file, err := pe.Open(fileName)
// 	parseErr(err)
// 	opt := mape.ConvertOptionalHeader(file)
// 	if newImageBase != opt.ImageBase {
// 		verbose("Aligning image base...", "*")
// 		rawFile, err := ioutil.ReadFile(fileName)
// 		parseErr(err)
// 		oFileHeader := int(rawFile[0x3c]) // Offset to EXE header
// 		if target.arch == "x64" {
// 			oImageBase := int(oFileHeader + 0x30)
// 			imageBase := binary.LittleEndian.Uint64(rawFile[oImageBase : oImageBase+8])
// 			verbose("ImageBase "+fmt.Sprintf("0x%016X -> ", imageBase)+fmt.Sprintf("0x%016X", newImageBase), "*")
// 			binary.LittleEndian.PutUint64(rawFile[oImageBase:oImageBase+8], newImageBase)
// 		} else {
// 			oImageBase := int(oFileHeader + 0x34)
// 			imageBase := binary.LittleEndian.Uint32(rawFile[oImageBase : oImageBase+4])
// 			verbose("ImageBase "+fmt.Sprintf("0x%08X -> ", imageBase)+fmt.Sprintf("0x%08X", newImageBase), "*")
// 			binary.LittleEndian.PutUint32(rawFile[oImageBase:oImageBase+4], uint32(newImageBase))
// 		}
// 		err = ioutil.WriteFile(target.FileName, rawFile, os.ModeDir)
// 		parseErr(err)
// 	}
// }

*/

func align(size, align, addr uint32) uint32 {
	if 0 == (size % align) {
		return addr + size
	}
	return addr + (size/align+1)*align
}
