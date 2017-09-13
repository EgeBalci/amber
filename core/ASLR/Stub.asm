;
;
;
; #########################
; #  ESI -> &PE           #
; #  EBP -> &hash_api     #
; #########################
;
; 
; Author: Ege Balcı <egebalci[at]protonmail[dot]com> 
; Version: 1.0

[BITS 32]
[ORG 0]

;%define VirtualAlloc

	pushad						; Save all registers to stack
	pushfd						; Save all flags to stack
	cld							; Clear out direction flags
	call Start					; Start OP
	%include "HASH-API.asm"		; hash_api
GetAOE:
	mov eax,[esi+0x3C]			; Get the offset of "PE" to eax
	mov ebx,[eax+esi+0x34]		; Get the image base address to ebx
	mov eax,[eax+esi+0x28]		; Get the address of entry point to eax
	ret							; <-
Start:
	pop ebp						; Pop the address of hash_api to ebp
	call Stub					; ...
PE:
	incbin "Mem.map"			; PE file image
	ImageSize: equ $-PE			; Size of the PE image
Stub:
	pop esi						; Get the address of image to esi
	push dword 0x40 			; PAGE_EXECUTE_READ_WRITE
	push dword 0x103000			; MEM_COMMI | MEM_TOP_DOWN | MEM_RESERVE
	push dword ImageSize		; dwSize
	push dword 0x00				; lpAddress
	push 0xE553A458				; hash( "kernel32.dll", "VirtualAlloc" )
	call ebp					; VirtualAlloc(lpAddress,ImageSize,MEM_COMMIT|MEM_TOP_DOWN|MEM_RESERVE,PAGE_EXECUTE_READWRITE)

	test eax,eax				; Check success 
	jz OpEnd					; If VirtualAlloc fails don't bother :/	
	push eax					; Save the new base address to stack
	call GetAOE					; Get the AOE and image base 	
	%include "Relocate.asm"		; Make image base relocation
	%include "BuildImportTable.asm"	; Call the module responsible for building the import address table
	xor ecx,ecx 				; Zero out the ECX
	call GetAOE					; Get image base and AOE
	mov ebx,[esp]				; Copy the address of new base to ebx
	add [esp],eax				; Add the AOE value to new base
Memcpy:	
	mov al,[esi] 				; Move 1 byte of PE image to AL register
	mov [ebx],al 				; Move 1 byte of PE image to image base
	inc esi 					; Increase PE image index
	inc ebx 					; Increase image base index
	inc ecx 					; Decrease loop counter
	cmp ecx,ImageSize 			; Check if ECX is 0
	jnz Memcpy 					; If not loop
	mov dword eax,[esp]			; Copy the AOEP to eax
CreateThread:
	pop ebx						; Pop back the AOE to ebx
	xor eax,eax					; Zero out the eax
	push eax					; lpThreadId
	push eax					; dwCreationFlags
 	push eax					; lpParameter
  	push ebx					; lpStartAddress
  	push eax					; dwStackSize
  	push eax					; lpThreadAttributes
  	push 0x160D6838 			; hash( "kernel32.dll", "CreateThread" )
  	call ebp 					; CreateThread( NULL, 0, &threadstart, NULL, 0, NULL );
  	jmp OpEnd					; <-
OpEnd:
	popfd						; Put back all saved flags
	popad						; Put back all saved registers
	ret							; Continue the execution
