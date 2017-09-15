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

%include "../../IAT.asm"

	call Stub				; ...
PE:
	incbin "Mem.map"		; PE file image
	ImageSize: equ $-PE		; Size of the PE image
Stub:
	pop esi					; Get the address of image to esi
	push dword 0x40 		; PAGE_EXECUTE_READ_WRITE
	push dword 0x103000		; MEM_COMMI | MEM_TOP_DOWN | MEM_RESERVE
	push dword ImageSize	; dwSize
	push dword 0x00			; lpAddress
	call [VA]				; VirtualAlloc(lpAddress,dwSize,MEM_COMMIT|MEM_TOP_DOWN|MEM_RESERVE, PAGE_EXECUTE_READWRITE)

	test eax,eax			; Check success 
	jz OpEnd				; If VirtualAlloc fails don't bother :/	
	push eax				; Save the new base address to stack
	call GetAOE				; Get the AOE and image base 	
	%include "Relocate.asm"	; Make image base relocation
	%include "BuildImportTable.asm"	; Call the module responsible for building the import address table
	xor ecx,ecx 			; Zero out the ECX
	call GetAOE				; Get image base and AOE
	mov ebx,[esp]			; Copy the address of new base to ebx
	add [esp],eax			; Add the AOE value to new base
Memcpy:	
	mov al,[esi] 			; Move 1 byte of PE image to AL register
	mov [ebx],al 			; Move 1 byte of PE image to image base
	inc esi 				; Increase PE image index
	inc ebx 				; Increase image base index
	inc ecx 				; Decrease loop counter
	cmp ecx,ImageSize 		; Check if ECX is 0
	jnz Memcpy 				; If not loop
	mov dword eax,[esp]		; Copy the AOEP to eax
CreateThread:
	pop ebx					; Pop back the AOE to ebx
	xor eax,eax				; Zero out the eax
	push eax				; lpThreadId
	push eax				; dwCreationFlags
 	push eax				; lpParameter
  	push ebx				; lpStartAddress
  	push eax				; dwStackSize
  	push eax				; lpThreadAttributes
  	call [CT] 				; CreateThread( NULL, 0, &threadstart, NULL, 0, NULL );
  	jmp OpEnd				; <-
GetAOE:
	mov eax,[esi+0x3C]		; Get the offset of "PE" to eax
	mov ebx,[eax+esi+0x34]	; Get the image base address to ebx
	mov eax,[eax+esi+0x28]	; Get the address of entry point to eax
	ret						; <-
OpEnd:
	nop						; Chill ;)
	jmp OpEnd				; To infinity and beyond !
