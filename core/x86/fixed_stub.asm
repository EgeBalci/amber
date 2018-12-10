; x86 Reflective Stub Fixed (No relocation)
; Author: Ege Balcı <egebalci@pm.me> 
; Version: 1.0
;
;#- stub.asm ------------------------------
; ESI = &PE
; RBP = &api.asm
; EBX = image_base
; STACK[0] = address_of_entry
; STACK[1] = image_base

[BITS 32]
[ORG 0]
	cld
	call stub                   ; ...
PE:
	incbin "Mem.map"            ; PE file image
	image_size: equ $-PE        ; Size of the PE image
stub:
	pop esi                     ; Get the address of image to esi
	call start                  ;
	%include "api.asm";
start:                          ;
	pop ebp                     ; Get the address of api to ebp
	mov eax,[esi+0x3C]          ; Get the offset of "PE" to eax
	mov ebx,[eax+esi+0x34]      ; Get the image base address to ebx
	mov eax,[eax+esi+0x28]      ; Get the address of entry point to eax
	push eax                    ; Save the adress of entry to stack
	push ebx                    ; Save the image base to stack
	push 0x00000000             ; Allocate a DWORD variable inside stack
	push esp                    ; lpflOldProtect
	push byte 0x40              ; PAGE_EXECUTE_READWRITE
	push image_size             ; dwSize
	push ebx                    ; lpAddress
	push 0xC38AE110             ; hash( "kernel32.dll", "VirtualProtect" )
	call ebp                    ; VirtualProtect( ImageBase, image_size, PAGE_EXECUTE_READWRITE, lpflOldProtect)
	test eax,eax                ; Check success 
	jz fail                     ; If VirtualProtect fails we are FUCKED !
	pop eax                     ; Fix the stack
;#- resolve.asm ------------------------------
; ESI = &PE
; EBX = HANDLE(module)
; ECX = &IAT
; STACK[0] = address_of_entry
; STACK[1] = image_base
; STACK[2] = &_IMAGE_IMPORT_DESCRIPTOR
; STACK[3] = _IMAGE_IMPORT_DESCRIPTOR->IAT (RVA)
;
	mov eax,[esi+0x3C]         ; Offset to IMAGE_NT_HEADER ("PE")
	mov eax,[eax+esi+0x80]     ; Import table RVA
	add eax,esi                ; Import table memory address (first image import descriptor)
	push eax                   ; Save the address of import descriptor to stack
get_modules:
	cmp dword [eax],0x00       ; Check if the import names table RVA is NULL
	jz complete                ; If yes building process is done
	mov eax,[eax+0x0C]         ; Get RVA of dll name to eax
	add eax,esi                ; Get the dll name address       
	call LoadLibraryA          ; Load the library
	mov ebx,eax                ; Move the dll handle to ebx
	mov eax,[esp]              ; Move the address of current _IMPORT_DESCRIPTOR to eax 
	call get_procs             ; Resolve all windows API function addresses
	add dword [esp],0x14       ; Move to the next import descriptor
	mov eax,[esp]              ; Set the new import descriptor address to eax
	jmp get_modules
get_procs:
	push ecx                   ; Save ecx to stack
	push dword [eax+0x10]      ; Save the current import descriptor IAT RVA
	add [esp],esi              ; Get the IAT memory address 
	mov eax,[eax]              ; Set the import names table RVA to eax
	add eax,esi                ; Get the current import descriptor's import names table address
	push eax                   ; Save it to stack
resolve: 
	cmp dword [eax],0x00       ; Check if end of the import names table
	jz all_resolved            ; If yes resolving process is done
	mov eax,[eax]              ; Get the RVA of function hint to eax
	cmp eax,0x80000000         ; Check if the high order bit is set
	js name_resolve            ; If high order bit is not set resolve with INT entry
	sub eax,0x80000000         ; Zero out the high bit
	call GetProcAddress        ; Get the API address with hint
	jmp insert_iat             ; Insert the address of API tı IAT
name_resolve:
	add eax,esi                ; Set the address of function hint
	add dword eax,0x02         ; Move to function name
	call GetProcAddress        ; Get the function address to eax
insert_iat:
	mov ecx,[esp+4]            ; Move the IAT address to ecx 
	mov [ecx],eax              ; Insert the function address to IAT
	add dword [esp],0x04       ; Increase the import names table index
	add dword [esp+4],0x04     ; Increase the IAT index
	mov eax,[esp]              ; Set the address of import names table address to eax
	jmp resolve                ; Loop
all_resolved:
	mov ecx,[esp+4]            ; Move the IAT address to ecx 
	mov dword [ecx],0x00       ; Insert a NULL dword
	pop ecx                    ; Deallocate index values
	pop ecx                    ; ...
	pop ecx                    ; Put back the ecx value
	ret                        ; <-
LoadLibraryA:
	push ecx                   ; Save ecx to stack
	push edx                   ; Save edx to stack
	push eax                   ; Push the address of linrary name string
	push 0x0726774C            ; hash( "kernel32.dll", "LoadLibraryA" )
	call ebp                   ; LoadLibraryA([esp+4])
	pop edx                    ; Retreive edx
	pop ecx                    ; Retreive ecx
	ret                        ; <-
GetProcAddress:
	push ecx                   ; Save ecx to stack
	push edx                   ; Save edx to stack
	push eax                   ; Push the address of proc name string
	push ebx                   ; Push the dll handle
	push 0x7802F749            ; hash( "kernel32.dll", "GetProcAddress" )
	call ebp                   ; GetProcAddress(ebx,[esp+4])
	pop edx                    ; Retrieve edx
	pop ecx                    ; Retrieve ecx
	ret                        ; <-
complete:
	pop eax                    ; Clean out the stack

;----------------------------------------------------------------------
; All done, now copy the image to new base and start a new thread
; ESI = &PE
; STACK[0] = address_of_entry 
; STACK[1] = image_base 
;
	xor ecx,ecx               ; Zero out the ECX
	pop ebx                   ; Pop the image base to EBX
	add [esp],ebx             ; Add image base to address of entry
	mov ecx,image_size        ; Move the image size to ECX
memcpy:
	mov al,[esi]              ; Move 1 byte of PE image to AL register
	mov [ebx],al              ; Move 1 byte of PE image to image base
	inc esi                   ; Increase PE image index
	inc ebx                   ; Increase image base index
	loop memcpy               ; Loop until ECX = 0
	ret                       ; Return to the AOEP
fail:
	ret                       ; <-
