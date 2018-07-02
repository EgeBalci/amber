; 
; Author: Ege Balcı <ege.balci@protonmail.com> 
; Version: 1.0

[BITS 32]
[ORG 0]
	cld
	call Stub				; ...
PE:
	incbin "Mem.map"		; PE file image
	ImageSize: equ $-PE		; Size of the PE image
Stub:
	pop esi					; Get the address of image to esi
	call IAT_API			;
	%include "iat_api.asm"	;
IAT_API:					;
	pop ebp					; Get the address of hook_api to ebp
	push dword 0x40 		; PAGE_EXECUTE_READ_WRITE
	push dword 0x103000		; MEM_COMMI | MEM_TOP_DOWN | MEM_RESERVE
	push dword ImageSize	; dwSize
	push dword 0x00			; lpAddress
	push 0xE553A458			; hash( "kernel32.dll", "VirtualAlloc" )
	call ebp				; VirtualAlloc(lpAddress,dwSize,MEM_COMMIT|MEM_TOP_DOWN|MEM_RESERVE, PAGE_EXECUTE_READWRITE)

	test eax,eax			; Check success 
	jz OpEnd				; If VirtualAlloc fails don't bother :/	
	push eax				; Save the new base address to stack
	call GetAOE				; Get the AOE and image base 	
	%include "relocate.asm"	; Make image base relocation
	%include "BuildImportTable.asm"	; Call the module responsible for building the import address table
	push 0x00000000			; Push NULL byte string terminator
	push 0x6c6c642e			; "lld."
	push 0x32336c65			; "23le"
	push 0x6e72656b			; "nrek"
	push esp				; Push the address of "kernel32.dll" string
	push 0x0726774C			; hash( "kernel32.dll","LoadLibraryA" )
	call ebp				; LoadLibraryA("kernel32.dll")
	push 0x00000000			; Push NULL byte string terminator
	push 0x64616572			; "daer"
	push 0x68546574			; "hTet"
	push 0x61657243			; "aerC"
	push esp				; Push the address of "CreateThread" string
	push eax				; Push the kernel32.dll handle
	push 0x7802F749			; hash( "kernel32.dll","GetProcAddress" )
	call ebp				; GetProcAddress(HANDLE,"CreateThread")
	add esp,0x20			; Clean the stack
	mov [esp+4],eax			; Save the address of CreateThread API to stack
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
	pop ebp					; Pop the address of CreateThread API to EBP
	xor eax,eax				; Zero out the eax
	push eax				; lpThreadId
	push eax				; dwCreationFlags
 	push eax				; lpParameter
  	push ebx				; lpStartAddress
  	push eax				; dwStackSize
  	push eax				; lpThreadAttributes
	call ebp				; CreateThread( NULL, 0, &threadstart, NULL, 0, NULL );
  	jmp OpEnd				; <-
GetAOE:
	mov eax,[esi+0x3C]		; Get the offset of "PE" to eax
	mov ebx,[eax+esi+0x34]	; Get the image base address to ebx
	mov eax,[eax+esi+0x28]	; Get the address of entry point to eax
	ret						; <-
OpEnd:
	nop						; Chill ;)
	jmp OpEnd				; To infinity and beyond !
