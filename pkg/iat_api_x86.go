package amber

// IAT32 https://github.com/EgeBalci/IAT_API
const IAT32 = `

api_call:
	pushad                  ; We preserve all the registers for the caller, bar EAX and ECX.
	xor eax,eax             ; Zero EAX (upper 3 bytes will remain zero until function is found)
	mov edx,fs:[eax+0x30]   ; Get a pointer to the PEB
	mov edx,[edx+0x0C]      ; Get PEB->Ldr
	mov edx,[edx+0x14]      ; Get the first module from the InMemoryOrder module list
	mov edx,[edx+0x10]      ; Get this modules base address
	push edx                ; Save the image base to stack (will use this alot)
	add edx,[edx+0x3C]      ; "PE" Header
	mov edx,[edx+0x80]      ; Import table RVA
	add edx,[esp]           ; Address of Import Table
	push edx                ; Save the &IT to stack (will use this alot)	
	mov esi,[esp+4]         ; Move the image base to ESI
	sub esp,0x08            ; Allocate space for import descriptor counter & hash
	sub edx,0x14            ; Prepare the import descriptor pointer for processing
next_desc:
	add edx,0x14            ; Get the next import descriptor
	cmp dword [edx],0x00    ; Check if import descriptor valid
	jz not_found            ; If import name array RVA is zero finish parsing
	mov esi,[esp+0x08]      ; Move the import table address to esi
	mov si,[edx+0x0C]       ; Get pointer to module name string RVA
	xor edi,edi             ; Clear EDI which will store the hash of the module name
loop_modname:               ;
	lodsb                   ; Read in the next byte of the name
	cmp al,'a'              ; Some versions of Windows use lower case module names
	jl not_lowercase        ;
	sub al, 0x20            ; If so normalise to uppercase
not_lowercase:              ;
	crc32 edi,al            ; Calculate CRC32 of module name
	crc32 edi,ah            ; Add NULL for unicode effect
	test al,al              ; Check if we read all
	jnz loop_modname
	; We now have the module hash computed
	mov [esp+4],edx         ; Save the current position in the module list for later
	mov [esp],edi           ; Save the current module hash for later
	; Proceed to iterate the export address table,
	mov ecx,[edx]           ; Get the RVA of import names table 
	add ecx,[esp+0x0C]      ; Add image base and get address of import names table
	sub ecx,0x04            ; Go 4 byte back	
get_next_func:
	mov edi,dword [esp]
	; use ecx as our EAT pointer here so we can take advantage of jecxz.
	add ecx,0x04            ; 4 byte forward
	cmp dword [ecx],0x00    ; Check if end of INT
	jz next_desc            ; If no INT present, process the next import descriptor
	mov esi,[ecx]           ; Get the RVA of func name hint
	cmp esi,0x80000000      ; Check if the high order bit is set 
	jns get_next_func       ; If not, there is no function name string :(
	add esi,[esp+0x0C]      ; Add the image base and get the address of function hint
	add dword esi,0x02      ; Move 2 bytes forward to asci function name
	; now ecx returns to its regularly scheduled counter duties
	; Computing the module hash + function hash
	; And compare it to the one we want
loop_funcname:              ;
	lodsb                   ; Read in the next byte of the ASCII function name
	crc32 edi,al            ; Calculate CRC32 of function name
	test al,al              ; Check if AL == 0
	jnz loop_funcname       ; If we have not reached the null terminator, continue
	cmp edi,[esp+0x34]      ; Compare the hash to the one we are searching for
	jnz get_next_func       ; Go compute the next function hash if we have not found it
	; If found, fix up stack, call the function and then value else compute the next one...
	mov eax,[edx+0x10]      ; Get the RVA of current descriptor's IAT 
	mov edx,[edx]           ; Get the import names table RVA of current import descriptor
	add edx,[esp+0x0C]      ; Get the address of import names table of current import descriptor
	sub ecx,edx             ; Find the function array index ?
	add eax,[esp+0x0C]      ; Add the image base to current descriptors IAT RVA
	add eax,ecx             ; Add the function index
	; Now we clean the stack	
	; We now fix up the stack and perform the call to the desired function...
finish:
	mov [esp+0x2C],eax      ; Overwrite the old EAX value with the desired api address for the upcoming popad
	add esp,0x10            ; Deallocate saved module hash, import descriptor address and import table address
	popad                   ; Restore all of the callers registers, bar EAX, ECX and EDX which are clobbered
	pop ebx                 ; Pop off the origional return address our caller will have pushed
	pop edx                 ; Pop off the hash value our caller will have pushed
	push ebx                ; Push back the return address
	mov eax,[eax]           ; Get the address of the desired API
	jmp eax                 ; Jump to target function
	; We now automagically return to the correct caller...
not_found:
	add esp,0x0F            ; Fix the stack
	popad                   ; Restore all registers
	ret                     ; Return

`
