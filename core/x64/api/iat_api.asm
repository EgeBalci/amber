;-----------------------------------------------------------------------------;
; Author: Ege Balcı (ege.balci[at]invictuseurope[dot]com)
; Compatible: Windows 7, 2003
; Architecture: x64
; Size: 200 bytes
;-----------------------------------------------------------------------------;

[BITS 64]

; Windows x64 calling convention:
; http://msdn.microsoft.com/en-us/library/9b372w95.aspx

; Input: The hash of the API to call in r10d and all its parameters (rcx/rdx/r8/r9/any stack params)
; Output: The return value from the API call will be in RAX.
; Clobbers: RAX, RCX, RDX, R8, R9, R10, R11
; Un-Clobbered: RBX, RSI, RDI, RBP, R12, R13, R14, R15.
;               RSP will be off by -40 hence the 'add rsp, 40' after each call to this function
; Note: This function assumes the direction flag has allready been cleared via a CLD instruction.
; Note: This function is unable to call forwarded exports.

%define ROTATION 13		; Rotation value for ROR hash

api_call:
  	push r9                 ; Save the 4th parameter
  	push r8                 ; Save the 3rd parameter
  	push rdx                ; Save the 2nd parameter
  	push rcx                ; Save the 1st parameter
  	push rsi                ; Save RSI
  	xor rdx,rdx            	; Zero rdx
  	mov rdx,[gs:rdx+96]    	; Get a pointer to the PEB
  	mov rdx,[rdx+24]       	; Get PEB->Ldr
  	mov rdx,[rdx+32]       	; Get the first module from the InMemoryOrder module list
  	mov rdx,[rdx+32]	   	; Get this modules base address
  	push rdx				; Save the image base to stack (will use this alot)
  	add dx,word [rdx+60]    ; "PE" Header
  	mov edx,dword [rdx+144]	; Import table RVA
 	add rdx,[rsp]			; Address of Import Table
	push rdx				; Save the &IT to stack (will use this alot)
  	mov rsi,[rsp+8]			; Move the image base to RSI
	sub rsp,16				; Allocate space for import descriptor counter & hash
	sub rdx,20				; Prepare import descriptor pointer for processing
next_desc:
	add rdx,20				; Get the next import descriptor
	cmp dword [rdx],0		; Check if import descriptor is valid
	jz not_found			; If import name array RVA is zero finish parsing
	mov rsi,[rsp+16]		; Move import table address to RSI
	mov si,[rdx+12]			; Get pointer to module name string RVA
	xor rdi,rdi				; Clear RDI which will store the hash of the module name
	xor rax,rax				; Clear RAX for calculating the hash
loop_modname:
	lodsb					; Read in the next byte of the name
	cmp al,'a'				; Some versions of windows use lower case module names
	jl not_lowercase		;
	sub al,32				; If so normalize to uppercase 
not_lowercase:
  	ror edi, ROTATION       ; Rotate right our hash value
  	add edi, eax            ; Add the next byte of the name
	ror edi,ROTATION		; In order to calculate the same hash values as Stephen Fewer's hash API we need to rotate one more and add a null byte.
	test al,al				; Check if we read all
	jnz loop_modname		; 
  	; We now have the module hash computed
	mov [rsp+8],rdx			; Save the current position in the module listfor later
	mov [rsp],edi			; Save the current module hash for later
  	; Proceed to itterate the export address table, 
  	mov ecx,dword [rdx]     ; Get RVA of import names table
  	add rcx,[rsp+24]  		; Add the image base and get the address of import names table
	sub rcx,8				; Go 4 bytes back
get_next_func:             	;
	add rcx,8				; 8 byte forward
	cmp dword [rcx],0		; Check if end of INT
	jz next_desc			; If no INT present, process the next import descriptor
	mov esi,dword [rcx]		; Get the RVA of func name hint
	cmp esi,0x80000000		; Check if the high order bit is set
	jns get_next_func		; If not, there is no function name string :(
	add rsi,[rsp+24]		; Add the image base and get the address of function name hint
	add rsi,2				; Move 2 bytes forward to asci function name
	; now ecx returns to its regularly scheduled counter duties
	; Computing the module hash + function hash
	xor rdi,rdi
	xor rax,rax
	; And compare it to the one we want
loop_funcname:
	lodsb                   ; Read in the next byte of the ASCII function name
  	ror edi,ROTATION        ; Rotate right our hash value
  	add edi,eax             ; Add the next byte of the name
  	cmp al,ah               ; Compare AL (the next byte from the name) to AH (null)
  	jne loop_funcname       ; If we have not reached the null terminator, continue
  	add edi,[rsp]          	; Add the current module hash to the function hash
  	cmp edi,r10d      		; Compare the hash to the one we are searchnig for 
  	jnz get_next_func       ; Go compute the next function hash if we have not found it
  ; If found, fix up stack, call the function and then value else compute the next one...
	mov eax,dword [rdx+16]	; Get the RVA of current descriptor's IAT
	mov edx,dword [rdx]		; Get the import names table RVA of current import descriptor
	add rdx,[rsp+24]		; Get the address of import names table of current import descriptor
	sub rcx,rdx				; Find the function array index ?
	add rax,[rsp+24]		; Add the image base to current descriptors IAT RVA
	add rax,rcx				; Add the function index
	; Now clean the stack
  	; We now fix up the stack and perform the call to the drsired function...
finish:
  	pop r8                   ; Clear off the current modules hash
 	pop r8                   ; Clear off the current position in the module list
	pop r8					 ; Clear off the import table address of last module
	pop r8 					 ; Clear off the image base address of last module
  	pop rsi                  ; Restore RSI
  	pop rcx                  ; Restore the 1st parameter
  	pop rdx                  ; Restore the 2nd parameter
  	pop r8                   ; Restore the 3rd parameter
  	pop r9                   ; Restore the 4th parameter
  	pop rdi                  ; pop off the return address
  	sub rsp, 32              ; reserve space for the four register params (4 * sizeof(QWORD) = 32)
   		                     ; It is the callers responsibility to restore RSP if need be (or alloc more space or align RSP).
	mov rax,[rax]			 ; Get the address of the desired API
	call rax				 ; Call the API
	push rdi				 ; Push back the return address
	ret						 ; Finito !
  ; We now automagically return to the correct caller...
not_found:
	add rsp,72			 	 ; Clean out the stack
	ret						 ; Return to caller
