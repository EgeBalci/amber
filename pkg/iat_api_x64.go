package amber

// IAT64 https://github.com/EgeBalci/IAT_API
const IAT64 = `

api_call:
	push r9                 ; Save the 4th parameter
	push r8                 ; Save the 3rd parameter
	push rdx                ; Save the 2nd parameter
	push rcx                ; Save the 1st parameter
	push rsi                ; Save RSI
	xor rdx,rdx             ; Zero rdx
 	mov rdx,gs:[rdx+0x60]   ; Get a pointer to the PEB
	mov rdx,[rdx+0x18]      ; Get PEB->Ldr
	mov rdx,[rdx+0x20]      ; Get the first module from the InMemoryOrder module list
	mov rdx,[rdx+0x20]      ; Get this modules base address
	push rdx                ; Save the image base to stack (will use this alot)
	add dx,word [rdx+0x3C]  ; "PE" Header
	mov edx,dword [rdx+0x90]; Import table RVA
	add rdx,[rsp]           ; Address of Import Table
	push rdx                ; Save the &IT to stack (will use this alot)
 	mov rsi,[rsp+8]         ; Move the image base to RSI
	sub rsp,0x10            ; Allocate space for import descriptor counter & hash
	sub rdx,0x14            ; Prepare import descriptor pointer for processing
next_desc:
	add rdx,0x14            ; Get the next import descriptor
	cmp dword [rdx],0       ; Check if import descriptor is valid
	jz not_found            ; If import name array RVA is zero finish parsing
	mov rsi,[rsp+0x10]      ; Move import table address to RSI
	mov si,[rdx+0xC]        ; Get pointer to module name string RVA
	xor rdi,rdi	            ; Clear RDI which will store the hash of the module name
loop_modname:
	xor rax,rax             ; Clear RAX for calculating the hash
	lodsb                   ; Read in the next byte of the name
	cmp al,'a'              ; Some versions of windows use lower case module names
	jl not_lowercase        ;
	sub al,0x20             ; If so normalize to uppercase 
not_lowercase:
	crc32 edi,al            ; Calculate CRC32 of module name
	crc32 edi,ah            ; Feed NULL for unicode effect
	test al,al              ; Check if end of the module name
	jnz loop_modname        ; 
	; We now have the module hash computed
	mov [rsp+8],rdx         ; Save the current position in the module listfor later
	mov [rsp],edi           ; Save the current module hash for later
	; Proceed to itterate the export address table, 
	mov ecx,dword [rdx]     ; Get RVA of import names table
	add rcx,[rsp+0x18]      ; Add the image base and get the address of import names table
	sub rcx,8               ; Go 4 bytes back
get_next_func:              ;
	mov rdi,[rsp]           ; Reset module hash
	add rcx,8               ; 8 byte forward
	cmp dword [rcx],0       ; Check if end of INT 
	jz next_desc            ; If no INT present, process the next import descriptor
	mov esi,dword [rcx]     ; Get the RVA of func name hint
  btr rax,0x3F            ; Check if the high order bit is set
  btr rsi,0x3F            ; Check if the high order bit is set
  jc get_next_func        ; If high order bit is not set resolve with INT entry
	add rsi,[rsp+0x18]      ; Add the image base and get the address of function name hint
	add rsi,2               ; Move 2 bytes forward to asci function name
	; now ecx returns to its regularly scheduled counter duties
	; Computing the module hash + function hash
	; And compare it to the one we want
loop_funcname:
	xor rax,rax             ; Clear RAX
	lodsb                   ; Read in the next byte of the ASCII function name
	crc32 edi,al            ; Calculate CRC32 of the function name
	cmp al,ah               ; Compare AL (the next byte from the name) to AH (null)
	jne loop_funcname       ; If we have not reached the null terminator, continue
	cmp edi,r10d            ; Compare the hash to the one we are searchnig for 
	jnz get_next_func       ; Go compute the next function hash if we have not found it
	; If found, fix up stack, call the function and then value else compute the next one...
	mov eax,dword [rdx+0x10]; Get the RVA of current descriptor's IAT
	mov edx,dword [rdx]     ; Get the import names table RVA of current import descriptor
	add rdx,[rsp+0x18]      ; Get the address of import names table of current import descriptor
	sub rcx,rdx             ; Find the function array index ?
	add rax,[rsp+0x18]      ; Add the image base to current descriptors IAT RVA
	add rax,rcx             ; Add the function index
	; Now clean the stack
	; We now fix up the stack and perform the call to the drsired function...
finish:
	pop r8                  ; Clear off the current modules hash
	pop r8                  ; Clear off the current position in the module list
	pop r8                  ; Clear off the import table address of last module
	pop r8                  ; Clear off the image base address of last module
	pop rsi                 ; Restore RSI
	pop rcx                 ; Restore the 1st parameter
	pop rdx                 ; Restore the 2nd parameter
	pop r8                  ; Restore the 3rd parameter
	pop r9                  ; Restore the 4th parameter
	pop r10                 ; Pop off the return address
	sub rsp,0x20            ; reserve space for the four register params (4 * sizeof(QWORD) = 32)
                            ; It is the callers responsibility to restore RSP if need be (or alloc more space or align RSP).
	push r10                ; Push back the return address
	mov rax,[rax]           ; Get the address of the desired API
	jmp rax                 ; Jump to target function
	; We now automagically return to the correct caller...
not_found:
	add rsp,0x48            ; Clean out the stack
	ret                     ; Return to caller

`
