;----------------------------------------------------------------------------- ;
; Authors: Michael Schierl, Ege Balcı
; Version: 1.0 (01 December 2018)
;----------------------------------------------------------------------------- ;
[BITS 64]

; Input: EBP - Data to decode
;        ESI - key
; ECX - Data size
; EDI - Scratch place for S-box
; Direction flag has to be cleared
; Output: None. Data is decoded in place.
; Clobbers: EAX, EBX, ECX, EDX, EBP

  cld
  call start
payload:
  incbin "payload.enc"
payload_size: equ $-payload
key:
	incbin "payload.key"
key_size: equ $-key
; Initialize S-box
start:
	pop rbp                   ; Pop out the address of payload to ebp
	lea rsi,[rbp+payload_size]; Load the address of key to esi
	mov rcx,payload_size      ; Move the size of the amber payload to ecx
	mov rdi,rsp               ; Set the address of stack as scratch box  
	xor rax, rax              ; Start with 0
init:
	stosb                     ; Store next SBox byte S[i] = i
	inc al                    ; Increase byte to write (EDI is increased automatically)
	jnz init                  ; Loop until we wrap around
	sub rdi, 0x100            ; Restore EDI
	; Permute S-box according to key
	xor rbx, rbx              ; Clear EBX (EAX is already cleared)
permute:
	add bl,[rdi+rax]          ; BL += S[AL] + KEY[AL % sizeof(key)]
	mov rdx,rax               ; 
	and dl,key_size-1         ; dl & sizeof(key)
	add bl,[rsi+rdx]          ; Move next byte of key to bl 
	mov dl,[rdi+rax]          ; swap S[AL] and S[BL]
	xchg dl,[rdi+rbx]         ; ..
	mov [rdi+rax], dl         ; ..
	inc al                    ; AL += 1 until we wrap around
	jnz permute               ;
	; Decryption loop
	xor rbx, rbx              ; Clear EBX (EAX is already cleared)
decrypt:
	inc al                    ; AL += 1
	add bl,[rdi+rax]          ; BL += S[AL]
	mov dl,[rdi+rax]          ; swap S[AL] and S[BL]
	xchg dl,[rdi+rbx]         ;
	mov [rdi+rax], dl         ;
	add dl,[rdi+rbx]          ; DL = S[AL]+S[BL]
	mov dl,[rdi+rdx]          ; DL = S[DL]
	xor [rbp],dl              ; [EBP] ^= DL
	inc rbp                   ; Advance data pointer
	dec rcx                   ; Reduce counter
	jnz decrypt               ; Until finished
	jmp payload
