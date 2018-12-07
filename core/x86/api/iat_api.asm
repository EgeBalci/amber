;-----------------------------------------------------------------------------;
; Author: Ege Balcı <ege.balci[at]invictuseurope[dot]com>
; Compatible: Windows 10/8.1/8/7/2008/Vista/2003/XP/2000/NT4
; Version: 1.0 (25 January 2018)
; Size: 177 bytes
;-----------------------------------------------------------------------------;

; This block locates addresses from import address table with given ror(13) hash value.
; Design is inpired from Stephen Fewer's hash api.

[BITS 32]

; Input: The hash of the API to call and all its parameters must be pushed onto stack.
; Output: The return value from the API call will be in EAX.
; Clobbers: EAX, EBX, ECX and EDX (NOT !! the normal stdcall calling convention because EBX is clobbered)
; Un-Clobbered: ESI, EDI, ESP and EBP can be expected to remain un-clobbered.
; Note: This function assumes the direction flag has allready been cleared via a CLD instruction.
; Note: This function is unable to call forwarded exports.

%define ROTATION 0x0D		; Rotation value for ROR hash

set_essentials:
  	pushad                 	; We preserve all the registers for the caller, bar EAX and ECX.
  	xor eax,eax           	; Zero EAX (upper 3 bytes will remain zero until function is found)
  	mov edx,[fs:eax+0x30] 	; Get a pointer to the PEB
  	mov edx,[edx+0x0C]		; Get PEB->Ldr
	mov edx,[edx+0x14]		; Get the first module from the InMemoryOrder module list
	mov edx,[edx+0x10]		; Get this modules base address
	push edx				; Save the image base to stack (will use this alot)
  	add edx,[edx+0x3C]     	; "PE" Header
	mov edx,[edx+0x80]		; Import table RVA
	add edx,[esp]			; Address of Import Table
	push edx				; Save the &IT to stack (will use this alot)	
	mov esi,[esp+4]			; Move the image base to ESI
	sub esp,0x08			; Allocate space for import descriptor counter & hash
	sub edx,0x14			; Prepare the import descriptor pointer for processing
next_desc:
	add edx,0x14			; Get the next import descriptor
	cmp dword [edx],0x00	; Check if import descriptor valid
	jz not_found			; If import name array RVA is zero finish parsing
	mov esi,[esp+0x08]		; Move the import table address to esi
	mov si,[edx+0x0C]     	; Get pointer to module name string RVA
	xor edi,edi           	; Clear EDI which will store the hash of the module name
loop_modname:            	;
	lodsb                  	; Read in the next byte of the name
	cmp al, 'a'            	; Some versions of Windows use lower case module names
	jl not_lowercase       	;
	sub al, 0x20           	; If so normalise to uppercase
not_lowercase:           	;
  	ror edi,ROTATION        ; Rotate right our hash value
 	add edi,eax           	; Add the next byte of the name
	ror edi,ROTATION		; In order to calculate the same hash values as Stephen Fewer's hash API we need to rotate one more and add a null byte.
  	test al,al				; Check if we read all
	jnz loop_modname
	; We now have the module hash computed
	mov [esp+4],edx	    	; Save the current position in the module list for later
	mov [esp],edi      		; Save the current module hash for later
  	; Proceed to iterate the export address table,
	mov ecx,[edx]      		; Get the RVA of import names table 
	add ecx,[esp+0x0C]      ; Add image base and get address of import names table
	sub ecx,0x04			; Go 4 byte back	
get_next_func:
  	; use ecx as our EAT pointer here so we can take advantage of jecxz.
  	add ecx,0x04			; 4 byte forward
  	cmp dword [ecx],0x00	; Check if end of INT
	jz next_desc    		; If no INT present, process the next import descriptor
  	mov esi,[ecx]           ; Get the RVA of func name hint
  	cmp esi,0x80000000      ; Check if the high order bit is set
	jns get_next_func		; If not, there is no function name string :(
	add esi,[esp+0x0C]		; Add the image base and get the address of function hint
	add dword esi,0x02		; Move 2 bytes forward to asci function name
  	; now ecx returns to its regularly scheduled counter duties
  	; Computing the module hash + function hash
  	xor edi,edi           	; Clear EDI which will store the hash of the function name
  	; And compare it to the one we want
loop_funcname:           	;
  	lodsb                  	; Read in the next byte of the ASCII function name
  	ror edi,ROTATION        ; Rotate right our hash value
  	add edi,eax           	; Add the next byte of the name
  	cmp al,ah             	; Compare AL (the next byte from the name) to AH (null)
  	jne loop_funcname      	; If we have not reached the null terminator, continue
  	add edi,[esp]       	; Add the current module hash to the function hash
  	cmp edi,[esp+0x34]      ; Compare the hash to the one we are searching for
  	jnz get_next_func      	; Go compute the next function hash if we have not found it
  	; If found, fix up stack, call the function and then value else compute the next one...
	mov eax,[edx+0x10]		; Get the RVA of current descriptor's IAT 
	mov edx,[edx]			; Get the import names table RVA of current import descriptor
	add edx,[esp+0x0C]		; Get the address of import names table of current import descriptor
	sub ecx,edx				; Find the function array index ?
	add eax,[esp+0x0C]		; Add the image base to current descriptors IAT RVA
	add eax,ecx				; Add the function index
	; Now we clean the stack	
  	
; We now fix up the stack and perform the call to the desired function...
finish:
  	mov [esp+0x2C],eax      ; Overwrite the old EAX value with the desired api address for the upcoming popad
	add esp,0x10			; Deallocate saved module hash, import descriptor address and import table address
  	popad                  	; Restore all of the callers registers, bar EAX, ECX and EDX which are clobbered
  	pop ebx                	; Pop off the origional return address our caller will have pushed
  	pop edx                	; Pop off the hash value our caller will have pushed
	mov eax,[eax]			; Get the address of the desired API
	call eax				; Call API
  	push ebx               	; Push back the return value
	ret						; 
  	; We now automagically return to the correct caller...
not_found:
	add esp,0x0F			; Fix the stack
	popad					; Restore all registers
	ret						; Return
	; (API is not found)
