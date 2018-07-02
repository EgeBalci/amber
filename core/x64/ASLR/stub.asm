; 
; Author: Ege Balcı <ege.balci@protonmail.com> 
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
;#- relocate.asm -------------------------------
;
;
;
;

[BITS 64]
[ORG 0]

	cld						; Clear direction flags
	call Stub				; Call Stub
PE:
	incbin "Mem.map"		; PE file image
	ImageSize: equ $-PE		; Size of the PE image
Stub:
	pop rsi					; Get the address of image to rsi
	call Start				; Call Start
	%include "block_api.asm";
Start:						;
	pop rbp					; Get the address of hook_api to rbp
	mov r9d,dword 0x40 		; PAGE_EXECUTE_READ_WRITE
	mov r8d,dword 0x103000	; MEM_COMMI | MEM_TOP_DOWN | MEM_RESERVE
	mov edx,dword ImageSize	; dwSize
	mov ecx,dword 0x00		; lpAddress
	mov r10d,0xE553A458		; hash( "kernel32.dll", "VirtualAlloc" )
	call rbp				; VirtualAlloc(lpAddress,dwSize,MEM_COMMIT|MEM_TOP_DOWN|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	test rax,rax			; Check success 
	jz OpEnd				; If VirtualAlloc fails don't bother :/	
	sub rsp,0x28			; Clear stack
	mov rdi,rax				; Save the new base address to rdi
	mov rax,[rsi+0x3C]		; Get the offset of "PE" to eax
	mov rbx,[rax+rsi+0x30]	; Get the image base address to rbx
	mov r12,[rax+rsi+0x28]	; Get the address of entry point to r12
	%include "relocate.asm"	; Make image base relocation
	%include "resolve.asm"	; Call the module responsible for building the import address table
	xor rcx,rcx 			; Zero out the ECX
	mov r13,rdi				; Copy the new base value to rbx
	add r13,r12				; Add the address of entry value to new base address
Memcpy:	
	mov al,[rsi] 			; Move 1 byte of PE image to AL register
	mov [rdi],al 			; Move 1 byte of PE image to image base
	inc rsi 				; Increase PE image index
	inc rdi 				; Increase image base index
	inc rcx 				; Decrease loop counter
	cmp rcx,ImageSize 		; Check if ECX is 0
	jnz Memcpy 				; If not loop
CreateThread:
	xor rax,rax				; Zero out the eax
	push rax				; lpThreadId
	push rax				; dwCreationFlags
	mov r9,rax				; lpParameter
	mov r8,r13				; lpStartAddress
 	mov rdx,rax				; dwStackSize
	mov rcx,rax				; lpThreadAttributes
	mov r10d,0x160D6838		; hash( "kernel32.dll","CreateThread" )
	call rbp				; CreateThread( NULL, 0, &threadstart, NULL, 0, NULL );
  	jmp OpEnd				; <-
OpEnd:
	nop						; Chill ;)
	jmp OpEnd				; To infinity and beyond !
