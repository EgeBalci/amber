[BITS 32]

; memcpy(&dst, &src, size)
memcpy:
	push ebp
	mov ebp, esp	
	push esi
	push edi
	push ecx
	mov edi,[ebp+8]
	mov esi,[ebp+12]
	mov ecx,[ebp+16]
copy_byte:
	rep movsb            ; Copy the CX number of bytes from RSI to RDI
	pop ecx
	pop edi
	pop esi
	mov esp,ebp
	pop ebp
	ret                  ; Return