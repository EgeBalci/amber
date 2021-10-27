package amber

import (
	"errors"
	"fmt"
	"strings"
)

// AssemblePayload generates the binary stub bla bla...
func (bp *Blueprint) AssemblePayload() ([]byte, error) {
	peMap, err := bp.file.CreateMemoryMapping()
	if err != nil {
		return nil, err
	}

	err = bp.file.PerformIntegrityChecks(peMap)
	if err != nil && !bp.IgnoreIntegrity {
		return nil, err
	}

	payload, err := bp.AddCallOver(peMap)
	if err != nil {
		return nil, err
	}

	stub, err := bp.GetLoaderAssembly()
	if err != nil {
		return nil, err
	}

	stub = strings.ReplaceAll(stub, ";", "\n;")
	stub = strings.ReplaceAll(stub, "$+", "")
	stubBin, ok := bp.Assemble(stub)
	if !ok {
		bp.printFaultingLine(stub)
		return nil, errors.New("core loader stub assembly failed")
	}

	return append(payload, stubBin...), nil
}

// GetLoaderAssembly returns the corresponding PE loader assembly code
// based on the given blueprint strunct
func (bp *Blueprint) GetLoaderAssembly() (string, error) {
	var stub, dllPrologue string

	if bp.IsDLL {
		switch bp.Architecture {
		case 64:
			dllPrologue = `		
			mov rcx,r13                     ; hinstDLL
			mov rdx,0x01                    ; fdwReason
			xor r8,r8                       ; lpReserved
	
			`
		case 32:
			dllPrologue = `
			push edi                ; AOE
			sub [esp],eax           ; hinstDLL
			push 0x01               ; fdwReason
			push 0x00               ; lpReserved

			`
		default:
			return "", errors.New("invalid architecture selected")
		}
	} else {
		dllPrologue = ""
	}

	api, err := bp.GetAPIResolverBlockAssembly()
	if err != nil {
		return "", err
	}

	switch bp.Architecture {
	case 32:

		if bp.HasRelocData {
			stub = LoaderX86
		} else {
			stub = FixedLoaderX86
		}
	case 64:

		if bp.HasRelocData {
			stub = LoaderX64
		} else {
			stub = FixedLoaderX64
		}
	default:
		return "", errors.New("invalid architecture selected")
	}

	return fmt.Sprintf(stub, api, dllPrologue), nil

}

// GetAPIResolverBlockAssembly returns the corresponding API resolver block assembly code
// based on the given blueprint strunct
func (bp *Blueprint) GetAPIResolverBlockAssembly() (string, error) {
	switch bp.Architecture {
	case 32:

		if bp.IAT {
			return IAT32, nil
		}
		return CRC32, nil

	case 64:

		if bp.IAT {
			return IAT64, nil
		}
		return CRC64, nil
	default:
		return "", errors.New("invalid architecture selected")
	}
}

// AddCallOver function adds a call instruction over the end of the given payload
// address of the payload will be pushed to the stack and execution will continiou after the end of payload
func (bp *Blueprint) AddCallOver(payload []byte) ([]byte, error) {

	// Perform a shport call over the payload
	call := fmt.Sprintf("call 0x%x", len(payload)+5)
	callBin, ok := bp.Assemble(call)
	if !ok {
		return nil, errors.New("call-over assembly failed")
	}
	payload = append(callBin, payload...)

	return payload, nil
}

func (bp *Blueprint) printFaultingLine(code string) {

	for i, line := range strings.Split(string(code), "\n") {
		if _, ok := bp.Assemble(line); !ok && !strings.Contains(line, "call") {
			fmt.Printf("%d: %s\n", i, line)
		}
	}

}
