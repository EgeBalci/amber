; X64 Reflective Stub Fixed (No relocation)
; Author: Ege Balcı <egebalci@pm.me> 
; Version: 1.0
;
;#- stub.asm ----------------------------------- 
; (RCX/RDX/R8/R9/R10/R11) = function_parameters
; R10D = hash("lib.dll", "function")
; RSI = &PE
; RBP = &block_api.asm
; RBX = pe_image_base
; RDI = new_image_base
; R12 = pe_address_of_entry
;

[BITS 64]
[ORG 0]

	cld                             ; Clear direction flags
	call stub                       ; Call stub
PE:
	incbin "Mem.map"                ; PE file image
	image_size: equ $-PE            ; Size of the PE image
stub:
	pop rsi                         ; Get the address of image to rsi
	call start                      ; Call start
	%include "api.asm"              ;
start:                              ;
	pop rbp                         ; Get the address of hook_api to rbp
	mov eax,dword [rsi+0x3C]        ; Get the offset of "PE" to eax
	mov rbx,qword [rax+rsi+0x30]    ; Get the image base address to rbx
	mov r12d,dword [rax+rsi+0x28]   ; Get the address of entry point to r12
	push rax                        ; Allocate 8 bytes for lpflOldProtect
	mov r9,rsp                      ; lpflOldProtect
	mov r8d,dword 0x40              ; PAGE_EXECUTE_READWRITE
	mov edx,dword image_size        ; dwSize
	mov rcx,rbx                     ; lpAddress
	mov r10d,0xC38AE110             ; hash( "kernel32.dll", "VirtualProtect" )
	call rbp                        ; VirtualProtect( image_base, image_size, PAGE_EXECUTE_READWRITE, lpflOldProtect)
	pop rdi                         ; Clear stack
	pop rdi                         ; ...
	mov rdi,rax                     ; Save the new base address to rdi

;#- resolve.asm --------------------------------
; STACK[0] = &_IMPORT_DESCRIPTOR
; R13 = Module HANDLE
; R14 = &IAT
; R15 = &INT
;
	xor r14,r14
	xor r15,r15
	mov eax,dword [rsi+0x3C]        ; Offset to IMAGE_NT_HEADER ("PE")
	mov eax,dword [rax+rsi+0x90]    ; Import table RVA
	add rax,rsi                     ; Import table memory address (first image import descriptor)
	push rax                        ; Save import descriptor to stack
get_modules:
	cmp dword [rax],0               ; Check if the import names table RVA is NULL
	jz complete                     ; If yes building process is done
	mov eax,dword [rax+0x0C]        ; Get RVA of dll name to eax
	add rax,rsi                     ; Get the dll name address
	call LoadLibraryA               ; Load the library
	mov r13,rax                     ; Move the dll handle to R13
	mov rax,[rsp]                   ; Move the address of current _IMPORT_DESCRIPTOR to eax 
	call get_procs                  ; Resolve all windows API function addresses
	add dword [rsp],0x14            ; Move to the next import descriptor
	mov rax,[rsp]                   ; Set the new import descriptor address to RAX
	jmp get_modules
get_procs:
	mov r14d,dword [rax+0x10]       ; Save the current import descriptor IAT RVA to R14D
	add r14,rsi                     ; Get the IAT memory address 
	mov rax,[rax]                   ; Set the import names table RVA to RAX
	add rax,rsi                     ; Get the current import descriptor's import names table address	
	mov r15,rax                     ; Save &INT to R15
resolve: 
	cmp dword [rax],0x00            ; Check if end of the import names table
	jz all_resolved                 ; If yes resolving stage is done
	mov rax,[rax]                   ; Get the RVA of function hint to eax
	cmp eax,0x80000000              ; Check if the high order bit is set
	js name_resolve                 ; If high order bit is not set resolve with INT entry
	sub eax,0x80000000              ; Zero out the high bit
	call GetProcAddress             ; Get API address with hint
	jmp insert_iat                  ; Insert the address of API to IAT
name_resolve:
	add rax,rsi                     ; Set the address of function hint
	add rax,0x02                    ; Move to function name
	call GetProcAddress             ; Get the function address to eax
insert_iat: 
	mov [r14],rax                   ; Insert the function address to IAT
	add r14,0x08                    ; Increase the IAT index
	add r15,0x08                    ; Increase the import names table index
	mov rax,r15                     ; Set the address of import names table address to RAX
	jmp resolve                     ; Loop
all_resolved:
	mov qword [r14],0x00            ; Insert a NULL dword
	ret                             ; <-
LoadLibraryA:
	push rcx                        ; Save ecx to stack
	mov rcx,rax                     ; Move the address of library name string to RCX
	mov r10d,0x0726774C             ; hash( "kernel32.dll", "LoadLibraryA" )
	call rbp                        ; LoadLibraryA(RCX)
	add rsp,32                      ; Fix the stack
	pop rcx                         ; Retreive ecx
	ret                             ; <-
GetProcAddress:
	mov rcx,r13                     ; Move the module handle to RCX as first parameter
	mov rdx,rax                     ; Move the address of function name string to RDX as second parameter
	mov r10d,0x7802F749             ; hash( "kernel32.dll", "GetProcAddress" )
	call rbp                        ; GetProcAddress(RCX,RDX)
	add rsp,24                      ; Fix the stack
	pop rdx                         ; ...
	ret                             ; <-
complete:
	pop rax                         ; Clean out the stack
;-----------------------------------------------------------------------------------
; All done now copy the image and run
; RSI = &PE
; RBP = &block_api.asm
; RBX = pe_image_base
; RDI = new_image_base
; R12 = pe_address_of_entry
; R13 = (RDI+R12)
;
	xor rcx,rcx                     ; Zero out the ECX
	mov rcx,image_size              ; Move the image size to RCX
	mov r13,rdi                     ; Copy the new base value to rbx
	add r13,r12                     ; Add the address of entry value to new base address
memcpy:	
	mov al,[rsi]                    ; Move 1 byte of PE image to AL register
	mov [rdi],al                    ; Move 1 byte of PE image to image base
	inc rsi                         ; Increase PE image index
	inc rdi                         ; Increase image base index
	loop memcpy                     ; Loop until zero
	push rdi                        ; Push the new image base to stack
	add [rsp],r12                   ; Add the address of entry
	ret                             ; <-