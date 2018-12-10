; X64 Reflective Stub
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
	mov r9d,dword 0x40              ; PAGE_EXECUTE_READ_WRITE
	mov r8d,dword 0x103000          ; MEM_COMMI | MEM_TOP_DOWN | MEM_RESERVE
	mov edx,dword image_size        ; dwSize
	mov ecx,dword 0x00              ; lpAddress
	mov r10d,0xE553A458             ; hash( "kernel32.dll", "VirtualAlloc" )
	call rbp                        ; VirtualAlloc(lpAddress,dwSize,MEM_COMMIT|MEM_TOP_DOWN|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	add rsp,0x20                    ; Clear stack
	mov rdi,rax                     ; Save the new base address to rdi

;#- relocate.asm -------------------------------
; RSI = &PE
; RCX = &end_of_base_realocation_table
; RDX = base_relocation_delta
; R8 = Reloc block descriptor
; R13D = Reloc block RVA
; R14D = Reloc block table size
;
	xor rax,rax                     ; Zero out the RAX
	xor r8,r8                       ; Zero out the R8
	xor r13,r13                     ; Zero out the R13
	xor r14,r14                     ; Zero out the R14
	mov eax,dword [rsi+0x3C]        ; Offset to IMAGE_NT_HEADER ("PE")
	mov ecx,dword [rax+rsi+0xB4]    ; Base relocation table size
	mov eax,dword [rax+rsi+0xB0]    ; Base relocation table RVA
	add rax,rsi                     ; Base relocation table memory address
	add rcx,rax                     ; End of base relocation table
calc_delta:
	mov rdx,rdi                     ; Move the new base address to rdx
	sub rdx,rbx                     ; Delta value
	mov r13d,dword [rax]            ; Move the reloc RVA to R13D
	mov r14d,dword [rax+4]          ; Move the reloc table size to R14D
	add rax,0x08                    ; Move to the reloc descriptor
	jmp fix                         ; Start fixing 
get_rva:
	cmp rcx,rax                     ; Check if the end of the reloc section
	jle reloc_fin                   ; If yes goto fin
	mov r13d,dword [rax]            ; Move the new reloc RVA
	mov r14d,dword [rax+4]          ; Move the new reloc table size
	add rax,0x08                    ; Move 8 bytes
fix:
	cmp r14w,0x08                   ; Check if the end of the reloc block
	jz get_rva                      ; If yes set the next block RVA
	mov r8w,word [rax]              ; Move the reloc desc to r8w
	cmp r8w,word 0x00               ; Check if it is a padding word
	je pass
	and r8w,0x0FFF                  ; Get the last 12 bits
	add r8d,r13d                    ; Add block RVA to desc value
	add r8,rsi                      ; Add the start address of the image
	add [r8],rdx                    ; Add the delta value to calculated absolute address
pass:
	sub r14d,0x02                   ; Decrease the index 
	add rax,0x02                    ; Move to the next reloc desc.
	xor r8,r8                       ; Zero out r8
	jmp fix                         ; Loop
reloc_fin:                          ; All done !

;#- resolve.asm --------------------------------
; (RCX/RDX/R8/R9/R10/R11) = function_parameters
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
	mov rax,[rsp]                   ; Set the new import descriptor address to eax
	jmp get_modules
get_procs:
	mov r14d,dword [rax+0x10]       ; Save the current import descriptor IAT RVA
	add r14,rsi                     ; Get the IAT memory address 
	mov rax,[rax]                   ; Set the import names table RVA to eax
	add rax,rsi                     ; Get the current import descriptor's import names table address	
	mov r15,rax                     ; Save &INT to R15
resolve: 
	cmp dword [rax],0x00            ; Check if end of the import names table
	jz all_resolved                 ; If yes resolving process is done
	mov rax,[rax]                   ; Get the RVA of function hint to eax
	cmp eax,0x80000000              ; Check if the high order bit is set
	js name_resolve                 ; If high order bit is not set resolve with INT entry
	sub eax,0x80000000              ; Zero out the high bit
	call GetProcAddress             ; Get the API address with hint
	jmp insert_iat                  ; Insert the address of API tı IAT
name_resolve:
	add rax,rsi                     ; Set the address of function hint
	add rax,0x02                    ; Move to function name
	call GetProcAddress             ; Get the function address to eax
insert_iat: 
	mov [r14],rax                   ; Insert the function address to IAT
	add r14,0x08                    ; Increase the IAT index
	add r15,0x08                    ; Increase the import names table index
	mov rax,r15                     ; Set the address of import names table address to eax
	jmp resolve                     ; Loop
all_resolved:
	mov qword [r14],0x00            ; Insert a NULL dword
	ret                             ; <-
LoadLibraryA:
	push rcx                        ; Save ecx to stack
	mov rcx,rax                     ; Move the address of library name string to RCX
	mov r10d,0x0726774C             ; hash( "kernel32.dll", "LoadLibraryA" )
	call rbp                        ; LoadLibraryA([esp+4])
	add rsp,32                      ; Fix the stack
	pop rcx                         ; Retreive ecx
	ret                             ; <-
GetProcAddress:
	mov rcx,r13                     ; Move the module handle to RCX as first parameter
	mov rdx,rax                     ; Move the address of function name string to RDX as second parameter
	mov r10d,0x7802F749             ; hash( "kernel32.dll", "GetProcAddress" )
	call rbp                        ; GetProcAddress(ebx,[esp+4])
	add rsp,24                      ; Fix the stack
	pop rcx                         ; ...
	ret                             ; <-
complete:
	pop rax                         ; Clean out the stack
;-----------------------------------------------------------------------------------
; All done now copy the image and run
; (RCX/RDX/R8/R9/R10/R11) = function_parameters
; R10D = hash("lib.dll", "function")
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
CreateThread:
	xor rax,rax                     ; Zero out the eax
	push rax                        ; lpThreadId
	push rax                        ; dwCreationFlags
	mov r9,rax                      ; lpParameter
	mov r8,r13                      ; lpstartAddress
 	mov rdx,rax                     ; dwStackSize
	mov rcx,rax                     ; lpThreadAttributes
	mov r10d,0x160D6838             ; hash( "kernel32.dll","CreateThread" )
	call rbp                        ; CreateThread( NULL, 0, &threadstart, NULL, 0, NULL );
  	jmp end                         ; <-
end:
	nop                             ; Chill ;)
	jmp end                         ; Loop !
