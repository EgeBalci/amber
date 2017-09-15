; 
; Author: Ege Balcı <egebalci[at]protonmail[dot]com> 
; Version: 1.0

[BITS 32]
[ORG 0]

%include "../../IAT.asm"

GetAOE:
	mov eax,[esi+0x3C]			; Get the offset of "PE" to eax
	mov ebx,[eax+esi+0x34]		; Get the image base address to ebx
	mov eax,[eax+esi+0x28]		; Get the address of entry point to eax
	ret							; <-
Start:
	pop ebp
	call Stub
PE:
	incbin "Mem.map"			; PE file image
	ImageSize: equ $-PE
Stub:
	pop esi						; Get the address of image to esi
	call GetAOE	
	push 0x00000000 			; Allocate a DWORD variable inside stack
	push esp					; lpflOldProtect
	push byte 0x40				; PAGE_EXECUTE_READWRITE
	push ImageSize				; dwSize
	push ebx					; lpAddress
	call [VP]					; VirtualProtect( ImageBase, ImageSize, PAGE_EXECUTE_READWRITE, lpflOldProtect)
	test eax,eax				; Check success 
	jz Fail						; If VirtualProtect fails don't bother :/
	%include "BuildImportTable.asm"	; Call the module responsible for building the import address table
	xor ecx,ecx 				; Zero out the ECX
	call GetAOE					; Get image base and AOE
	push ebx					; Store the image base to stack
	add [esp],eax				; Add the AOE value
Memcpy:	
	mov al,[esi] 				; Move 1 byte of PE image to AL register
	mov [ebx],al 				; Move 1 byte of PE image to image base
	inc esi 					; Increase PE image index
	inc ebx 					; Increase image base index
	inc ecx 					; Decrease loop counter
	cmp ecx,ImageSize 			; Check if ECX is 0
	jnz Memcpy 					; If not loop
	mov dword eax,[esp]			; Copy the AOEP to eax
	ret							; Return to the AOEP
Fail:
	ret							; VirtualProtect failed :(
