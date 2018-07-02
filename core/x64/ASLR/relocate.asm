; 
; Author: Ege Balcı <ege.balci@protonmail.com> 
; Version: 1.0
;
;#- stub.asm ------------------------------------ 
; (RCX/RDX/R8/R9/R10/R11) = function_parameters
; R10D = hash("lib.dll", "function")
; RSI = &PE
; RBP = &block_api.asm
; RBX = pe_image_base
; RDI = new_image_base
; R12 = pe_address_of_entry
;
;#- relocate.asm -------------------------------
; RCX = &end_of_base_realocation_table
; RDX = base_relocation_delta
; R8 = 
;
;

[BITS 64]
[ORG 0]
	
	xor r8,r8				; Zero out the r8	
	mov rax,[rsi+0x3C]      ; Offset to IMAGE_NT_HEADER ("PE")
	mov rcx,[rax+rsi+0xB4]	; Base relocation table size
	mov rax,[rax+rsi+0xB0]  ; Base relocation table RVA
	add rax,rsi             ; Base relocation table memory address
	add rcx,rax				; End of base relocation table
CalcDelta:
	mov rdx,rdi				; Move the new base address to edi
	sub rdx,rbx				; Delta value
	push dword [rax]		; Reloc RVA
	push dword [rax+4]		; Reloc table size
	add rax,0x08			; Move to the reloc descriptor
	jmp Fix					; Start fixing 
GetRVA:
	cmp rcx,rax				; Check if the end of the reloc section ?
	jle RelocFin			; If yes goto fin
	add rsp,0x08			; Deallocate old reloc RVA and reloc table size variables
	push dword [rax]		; Push new reloc RVA
	push dword [rax+4]		; Push new reloc table size
	add rax,0x08			; Move 8 bytes
Fix:
	cmp word [rsp],0x08		; Check if the end of the reloc block
	jz GetRVA				; If yes set the next block RVA
	mov r8w,word [rax]		; Move the reloc desc to dx
	cmp r8w,word 0x00		; Check if it is a padding word
	je Pass
	and r8w,0x0FFF			; Get the last 12 bits
	add r8,[rsp+4]			; Add block RVA to desc value
	add r8,rsi				; Add the start address of the image
	add dword [r8],rdi		; Add the delta value to calculated absolute address
Pass:
	sub dword [rsp],0x02	; Decrease the index 
	add rax,0x02			; Move to the next reloc desc.
	xor r8,r8				; Zero out edx
	jmp Fix					; Loop
RelocFin:
	add rsp,0x08			; Deallocate all vars



