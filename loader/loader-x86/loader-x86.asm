;#===========================================#
;# x86 Reflective Loader                     #
;# Author: Ege Balcı <egebalci@pm.me>        #
;# Version: 3.0                              #
;#===========================================#

[BITS 32]

%define e_lfanew 0x3C
%define _AddressOfEntry 0x28

	call start                        ; Start by calling over the PE image
	incbin "putty.exe"                ; PE image
start:
	pop esi                           ; Get the PE address into esi
	push ebp                          ; Save EBP
	mov ebp,esp                       ; Create a stack frame
	push esi                          ; Push the PE address as first parameter
	call map_image                    ; Perform PE image mapping
	pop esi                           ; Pop out the PE address
	push eax                          ; Push new image baes to stack
	call relocate_image               ; Perform image relocation
	call resolve_imports              ; Resolve image imports & create IAT table
	call protect_sections             ; Apply proper section memory protections
	call run_tls_callbacks            ; Try to execute TLS callbacks. May fail... ¯\_(ツ)_/¯ 
	pop edi                           ; Get the new image base value into edi
	mov eax,[edi+e_lfanew]            ; Get the file header offset
	mov eax,[edi+eax+_AddressOfEntry] ; Get the AddressOfEntry into eax
	add eax,edi                       ; Add the AOE onto new image base
	jmp wipe                          ; Wipe memory artifacts
	; ------------------------ FUNCTIONS ------------------------------------
	%include "./inc/memcpy.asm"
	%include "./inc/calc_crc.asm"
	%include "./inc/map_image.asm"
	%include "./inc/load_module.asm"
	%include "./inc/relocate_image.asm"
	%include "./inc/resolve_imports.asm"
	%include "./inc/get_proc_by_crc.asm"
	%include "./inc/get_module_by_crc.asm"
	%include "./inc/protect_sections.asm"
	%include "./inc/run_tls_callbacks.asm"
	%include "../crc32_api/crc32_api_x86.asm"
	;------------------------ FUNCTIONS -------------------------------------
wipe:
	wipe_len_delta	equ     wipe_end-wipe
	call $+5                          ; Get current EIP to stack
	pop ecx                           ; Pop currect EIP to RCX
	sub ecx,esi                       ; Calculate the size of the PE file
	add ecx,wipe_len_delta            ; Add the size of wipe code 
	mov edi,esi                       ; Move the PE address to RDI
	sub edi,0x5                       ; Go back 5 bytes for wiping the initial call as well
wipe_end:
	rep stosb                         ; Store AL into RDI and decrement RDI until RCX = 0
	; -------------------- SWITCH TO PE ----------------------------
	cld                               ; Clear direction flags
	mov esp, ebp                      ; Restore stack frame
	pop ebp                           ; Restore RBP
	jmp eax                           ; Jmp to the PE->AOE