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

	var stub string

	switch bp.Architecture {
	case 32:

		if bp.HasRelocData {
			stub = STUB32
		} else {
			stub = FIXED_STUB32
		}

		if bp.IAT {
			stub += IAT_API_32
		} else {
			stub += CRC_API_32
		}

	case 64:

		if bp.HasRelocData {
			stub = STUB64
		} else {
			stub = FIXED_STUB64
		}

		if bp.IAT {
			stub += IAT_API_64
		} else {
			stub += CRC_API_64
		}

	}

	stub = strings.ReplaceAll(stub, ";", "\n;")
	stub = strings.ReplaceAll(stub, "$+", "")
	stub = strings.ReplaceAll(stub, "fs:", "fs+")
	stubBin, ok := bp.Assemble(stub)
	if !ok {
		bp.printFaultingLine(stub)
		return nil, errors.New("core loader stub assembly failed")
	}

	return append(payload, stubBin...), nil
}

// AddCallOver function adds a call instruction over the end of the given payload
// address of the payload will be pushed to the stack and execution will continiou after the end of payload
func (bp Blueprint) AddCallOver(payload []byte) ([]byte, error) {

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
