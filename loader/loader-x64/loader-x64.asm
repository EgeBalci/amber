;#==============================================#
;# X64 Reflective Loader                        #
;# Author: Ege Balcı <egebalci@pm.me>           #
;# Version: 3.0                                 #
;#==============================================#
;
[BITS 64]

%define e_lfanew 0x3C
%define _AddressOfEntry 0x28

  call start                                 ; Get the address of PE image to stack
  incbin "putty.exe"                         ; PE file.
start:
	pop rsi                                  ; Get the address of PE to RSI
	push rbp                                 ; Save RBP
	mov rbp,rsp                              ; Create a stack frame
	mov rcx,rsi                              ; Move the image address as first parameter
	call map_image                           ; Perform PE image mapping
	mov rdi, rax                             ; Get the address of mapped PE image into RDI
	mov rcx, rdi                             ; Move a copy of the mapped image address into RCX as first parameter
	call resolve_imports                     ; Resolve image imports
	mov rcx, rdi                             ; Set the mapped image address as first parameter
	call relocate_image                      ; Perform image base relocation
	mov rcx, rdi                             ; Set the mapped image address as first parameter
	call protect_sections                    ; Apply proper section memory protections
	mov rcx, rdi                             ; Set the mapped image address as first parameter
	call run_tls_callbacks                   ; Try to execute TLS callbacks. May fail... ¯\_(ツ)_/¯
	xor rax, rax                             ; Clear out RAX
	mov eax, DWORD [rdi+e_lfanew]            ; Get the file header offset
	mov eax, DWORD [rdi+rax+_AddressOfEntry] ; Get the AddressOfEntry into EAX
	add rax,rdi                              ; Add the AOE onto new image base
	jmp wipe                                 ; Start wiping memory artifacts...
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
	%include "../crc32_api/crc32_api_x64.asm"
	;------------------------ FUNCTIONS -------------------------------------
wipe:
	wipe_len_delta	equ     wipe_end-wipe
	call $+5                          ; Get current EIP to stack
	pop rcx                           ; Pop currect EIP to RCX
	sub rcx,rsi                       ; Calculate the size of the PE file
	add rcx,wipe_len_delta            ; Add the size of wipe code 
	mov rdi,rsi                       ; Move the PE address to RDI
	sub rdi,0x5                       ; Go back 5 bytes for wiping the initial call as well
wipe_end:
	rep stosb                         ; Store AL into RDI and decrement RDI until RCX = 0
	; -------------------- SWITCH TO PE ----------------------------
	cld                               ; Clear direction flags
	mov rsp, rbp                      ; Restore stack frame
	pop rbp                           ; Restore RBP
	jmp rax                           ; Jmp to the PE->AOE