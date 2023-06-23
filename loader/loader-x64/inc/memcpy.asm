[BITS 64]

; memcpy(&dst, &src, size)
; RCX = &dst
; RDX = &src
; R8  = size
memcpy:
	push rsi
	push rdi
	mov rdi,rcx
	mov rsi,rdx
	mov rcx,r8
copy_byte:
	rep movsb            ; Copy the CX number of bytes from RSI to RDI
	pop rdi              ; Restore RDI
	pop rsi              ; Restore RSI
	ret                  ; Return