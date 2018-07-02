; 
; Author: Ege Balcı <ege.balci@protonmail.com> 
; Version: 1.0

[BITS 32]
[ORG 0]
	
	xor edx,edx				; Zero out the edx	
Relocate:
	mov eax,[esi+0x3C]      ; Offset to IMAGE_NT_HEADER ("PE")
	mov ecx,[eax+esi+0xA4]	; Base relocation table size
	mov eax,[eax+esi+0xA0]  ; Base relocation table RVA
	add eax,esi             ; Base relocation table memory address
	add ecx,eax				; End of base relocation table
CalcDelta:
	mov edi,[esp]			; Move the new base address to edi
	sub edi,ebx				; Delta value
	push dword [eax]		; Reloc RVA
	push dword [eax+4]		; Reloc table size
	add eax,0x08			; Move to the reloc descriptor
	jmp Fix					; Start fixing 
GetRVA:
	cmp ecx,eax				; Check if the end of the reloc section ?
	jle RelocFin			; If yes goto fin
	add esp,0x08			; Deallocate old reloc RVA and reloc table size variables
	push dword [eax]		; Push new reloc RVA
	push dword [eax+4]		; Push new reloc table size
	add eax,0x08			; Move 8 bytes
Fix:
	cmp word [esp],0x08		; Check if the end of the reloc block
	jz GetRVA				; If yes set the next block RVA
	mov dx,word [eax]		; Move the reloc desc to dx
	cmp dx,word 0x00		; Check if it is a padding word
	je Pass
	and dx,0x0FFF			; Get the last 12 bits
	add edx,[esp+4]			; Add block RVA to desc value
	add edx,esi				; Add the start address of the image
	add dword [edx],edi		; Add the delta value to calculated absolute address
Pass:
	sub dword [esp],0x02	; Decrease the index 
	add eax,0x02			; Move to the next reloc desc.
	xor edx,edx				; Zero out edx
	jmp Fix					; Loop
RelocFin:
	add esp,0x08			; Deallocate all vars



