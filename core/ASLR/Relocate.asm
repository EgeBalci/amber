; This block requires following values inside the specified registers...
;
; ############################################
; #  Stack[0] -> New base address            #
; #  EBX      -> Deault image base          #
; ############################################
; 
; Author: Ege Balcı <egebalci[at]protonmail[dot]com> 
; Version: 1.0

[BITS 32]
[ORG 0]
	
	xor edx,edx		; Zero out the edx	
Relocate:
	mov eax,[esi+0x3C]      ; Offset to IMAGE_NT_HEADER ("PE")
	mov ecx,[eax+esi+0xA4]	; Base relocation table size
	mov eax,[eax+esi+0xA0]  ; Base relocation table RVA
	add eax,esi             ; Base relocation table memory address
	add ecx,eax		; End of base relocation table
CalcDelta:
	mov edi,[esp]		; Move the new base address to edi
	sub edi,ebx		; Delta value
	push dword [eax]	; Reloc RVA
	add eax,0x08		; Move to the reloc descriptor
	jmp Fix			; Start fixing 
GetRVA:
	add esp,0x04		; Deallocate old reloc RVA
	push dword [eax]	; Push new reloc RVA
	add eax,0x08		; Move 8 bytes
Fix:
	cmp ecx,eax		; Check if the end of the reloc section ?
	jle RelocFin		; If yes goto fin
	cmp word [eax],0x00	; Check if the end of the reloc block
	jz BlockEnd		; If yes set the next block RVA
	mov dx,word [eax]	; Move the reloc desc to dl
	and dx,0x0FFF		; Get the last 12 bits
	add edx,[esp]		; Add block RVA to desc value
	add edx,esi		; Add the start address of the image
	add dword [edx],edi	; Add the delta value to calculated absolute address 
	add eax,0x02		; Move to the next reloc desc.
	xor edx,edx		; Zero out edx
	jmp Fix			; Loop
BlockEnd:
	add eax,0x02		; Move 2 byte further
	jmp GetRVA		; ...
RelocFin:
	add esp,0x04		; Deallocate all vars



