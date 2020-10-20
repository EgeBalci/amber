package amber

/*
import (
	"bytes"
	"math/rand"

	"github.com/EgeBalci/debug/pe"
)

// CompileStub generates the final stub file with given payload
func (bp *Blueprint) CompileStub(payload []byte) ([]byte, error) {
	var (
		peFile *pe.File
		err    error
	)

	switch bp.architecture {
	case 32:
		peFile, err = pe.NewFile(bytes.NewReader(PE_STUB32))
	case 64:
		peFile, err = pe.NewFile(bytes.NewReader(PE_STUB64))
	}

	if err != nil {
		return nil, err
	}

	var entryPoint, sectionAlignment, fileAlignment, scAddr uint32
	var imageBase uint64

	lastSection := peFile.Sections[peFile.NumberOfSections-1]

	switch hdr := (peFile.OptionalHeader).(type) {
	case *pe.OptionalHeader32:
		imageBase = uint64(hdr.ImageBase) // cast this back to a uint32 before use in 32bit
		entryPoint = hdr.AddressOfEntryPoint
		sectionAlignment = hdr.SectionAlignment
		fileAlignment = hdr.FileAlignment
		scAddr = align(lastSection.Size, fileAlignment, lastSection.Offset) //PointerToRawData
		//shellcode = api.ApplySuffixJmpIntel32(shellcodeBytes, scAddr, entryPoint+uint32(imageBase), binary.LittleEndian)
		break
	case *pe.OptionalHeader64:
		imageBase = hdr.ImageBase
		entryPoint = hdr.AddressOfEntryPoint
		sectionAlignment = hdr.SectionAlignment
		fileAlignment = hdr.FileAlignment
		scAddr = align(lastSection.Size, fileAlignment, lastSection.Offset) //PointerToRawData
		//shellcode = api.ApplySuffixJmpIntel32(shellcodeBytes, scAddr, entryPoint+uint32(imageBase), binary.LittleEndian)
		break
	}

	newsection := new(pe.Section)
	newsection.Name = "." + randomString(rand.Intn(5)+1)
	o := []byte(newsection.Name)
	newsection.OriginalName = [8]byte{o[0], o[1], o[2], o[3], o[4], o[5], 0, 0}
	newsection.VirtualSize = uint32(len(payload))
	newsection.VirtualAddress = align(lastSection.VirtualSize, sectionAlignment, lastSection.VirtualAddress)
	newsection.Size = align(uint32(len(payload)), fileAlignment, 0)                //SizeOfRawData
	newsection.Offset = align(lastSection.Size, fileAlignment, lastSection.Offset) //PointerToRawData
	newsection.Characteristics = pe.IMAGE_SCN_CNT_CODE | pe.IMAGE_SCN_MEM_EXECUTE | pe.IMAGE_SCN_MEM_READ

	peFile.InsertionAddr = scAddr
	peFile.InsertionBytes = payload

	switch hdr := (peFile.OptionalHeader).(type) {
	case *pe.OptionalHeader32:
		v := newsection.VirtualSize
		if v == 0 {
			v = newsection.Size // SizeOfRawData
		}
		hdr.SizeOfImage = align(v, sectionAlignment, newsection.VirtualAddress)
		//hdr.AddressOfEntryPoint = bp.addressOfEntryPoint
		hdr.CheckSum = 0
		// disable ASLR
		//hdr.DllCharacteristics ^= pe.IMAGE_DLLCHARACTERISTICS_DYNAMIC_BASE
		//hdr.DataDirectory[5].VirtualAddress = 0
		//hdr.DataDirectory[5].Size = 0
		//peFile.FileHeader.Characteristics |= pe.IMAGE_FILE_RELOCS_STRIPPED
		//disable DEP
		//hdr.DllCharacteristics ^= pe.IMAGE_DLLCHARACTERISTICS_NX_COMPAT
		// zero out cert table offset and size
		//hdr.DataDirectory[4].VirtualAddress = 0
		//hdr.DataDirectory[4].Size = 0
		break
	case *pe.OptionalHeader64:
		v := newsection.VirtualSize
		if v == 0 {
			v = newsection.Size // SizeOfRawData
		}
		hdr.SizeOfImage = align(v, sectionAlignment, newsection.VirtualAddress)
		//hdr.AddressOfEntryPoint = bp.addressOfEntryPoint
		hdr.CheckSum = 0
		// disable ASLR
		//hdr.DllCharacteristics ^= pe.IMAGE_DLLCHARACTERISTICS_DYNAMIC_BASE
		//hdr.DataDirectory[5].VirtualAddress = 0
		//hdr.DataDirectory[5].Size = 0
		//peFile.FileHeader.Characteristics |= pe.IMAGE_FILE_RELOCS_STRIPPED
		//disable DEP
		//hdr.DllCharacteristics ^= pe.IMAGE_DLLCHARACTERISTICS_NX_COMPAT
		// zero out cert table offset and size
		//hdr.DataDirectory[4].VirtualAddress = 0
		//hdr.DataDirectory[4].Size = 0
		break
	}

	peFile.FileHeader.NumberOfSections++
	peFile.Sections = append(peFile.Sections, newsection)

	return peFile.Bytes()
}

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

func align(size, align, addr uint32) uint32 {
	if 0 == (size % align) {
		return addr + size
	}
	return addr + (size/align+1)*align
}

*/
