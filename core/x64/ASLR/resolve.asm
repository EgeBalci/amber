; 
; Author: Ege Balcı <ege.balci@protonmail.com> 
; Version: 1.0
;
;#- stub.asm ------------------------------------ 
; (RCX/RDX/R8/R9/R10/R11) = function_parameters
; R10D = hash("lib.dll", "function")
; RSI = &PE
; RBP = &block_api.asm
; RBX = pe_image_base
; RDI = new_image_base
; R12 = pe_address_of_entry
;
;#- relocate.asm ------------------------------- 
; RCX = &end_of_base_realocation_table
; RDX = base_relocation_delta
; R8 = 
;
;#- resolve.asm --------------------------------
; STACK[0] = &_IMPORT_DESCRIPTOR
; R13 = Module HANDLE
; R14 = &IAT
; R15 = &INT
;
;

[BITS 64]
[ORG 0]

	mov rax,[rsi+0x3C]		; Offset to IMAGE_NT_HEADER ("PE")
	mov rax,[rax+rsi+0x90] 	; Import table RVA
	add rax,rsi				; Import table memory address (first image import descriptor)
	push rax				; Save import descriptor to stack
GetDLLs:
	cmp dword [rax],0x00	; Check if the import names table RVA is NULL
	jz Complete				; If yes building process is done
	mov rax,[rax+0x0C]		; Get RVA of dll name to eax
	add rax,rsi				; Get the dll name address		
	call LoadLibraryA		; Load the library
	mov r13,rax 			; Move the dll handle to R13
	mov rax,[rsp]			; Move the address of current _IMPORT_DESCRIPTOR to eax 
	call GetProcs			; Resolve all windows API function addresses
	add dword [rsp],0x14	; Move to the next import descriptor
	mov rax,[rsp]			; Set the new import descriptor address to eax
	jmp GetDLLs
;-----------------------------------------------------------------------------------
GetProcs:
	mov r14,dword [rax+0x10]	; Save the current import descriptor IAT RVA
	add r14,rsi					; Get the IAT memory address 
	mov rax,[rax]				; Set the import names table RVA to eax
	add rax,rsi					; Get the current import descriptor's import names table address	
	mov r15,rax					; Save &INT to R15
Resolve: 
	cmp dword [rax],0x00 		; Check if end of the import names table
	jz AllResolved				; If yes resolving process is done
	mov rax,[rax]				; Get the RVA of function hint to eax
	cmp rax,0x8000000000000000	; Check if the high order bit is set
	js NameResolve				; If high order bit is not set resolve with INT entry
	sub rax,0x800000000000000	; Zero out the high bit
	call GetProcAddress			; Get the API address with hint
	jmp InsertIAT				; Insert the address of API tı IAT
NameResolve:
	add rax,rsi					; Set the address of function hint
	add dword rax,0x02			; Move to function name
	call GetProcAddress			; Get the function address to eax
InsertIAT:
	mov rcx,r14					; Move the IAT address to ecx 
	mov [rcx],rax				; Insert the function address to IAT
	add r14,0x04				; Increase the import names table index
	add r15,0x04				; Increase the IAT index
	mov rax,r15					; Set the address of import names table address to eax
	jmp Resolve					; Loop
AllResolved:
	mov rcx,[rsp+4]         ; Move the IAT address to ecx 
	mov dword [rcx],0x00	; Insert a NULL dword
	add rsp,0x08			; Deallocate index values
	pop rcx					; Put back the ecx value
	ret						; <-
;-----------------------------------------------------------------------------------
LoadLibraryA:
	push rcx				; Save ecx to stack
	mov rcx,rax 			; Move the address of library name string to RCX
	mov r10d,0x0726774C 	; hash( "kernel32.dll", "LoadLibraryA" )
	call rbp 				; LoadLibraryA([esp+4])
	sub rsp,0x28			; Clear the stack
	pop rcx					; Retreive ecx
	ret 					; <-
;-----------------------------------------------------------------------------------
GetProcAddress:
	mov rcx,r13 			; Move the module handle to RCX as first parameter
	mov rdx,rax				; Save edx to stack
	mov r10d,0x7802F749		; hash( "kernel32.dll", "GetProcAddress" )
	call rbp				; GetProcAddress(ebx,[esp+4])
	sub rsp,0x28			; Retrieve ecx
	ret 					; <-
;-----------------------------------------------------------------------------------
Complete:
	pop rax					; Clean out the stack
