;-----------------------------------------------------------------------------;
; Authors: Michael Schierl, Ege Balcı
; Version: 2.0 (02 December 2017)
;-----------------------------------------------------------------------------;
[BITS 64]

; Input: EBP - Data to decode
;        ESI - Key
; ECX - Data size
; EDI - Scratch place for S-box
; Direction flag has to be cleared
; Output: None. Data is decoded in place.
; Clobbers: EAX, EBX, ECX, EDX, EBP

  cld
  call start
Payload:
  incbin "payload.enc"
PSize: equ $-Payload
Key:
	incbin "payload.key"
KSize: equ $-Key
; Initialize S-box
start:
	pop rbp                  ; Pop out the address of payload to ebp
  lea rsi,[rbp+PSize]      ; Load the address of key to esi
	mov rcx,PSize            ; Move the size of the amber payload to ecx
	mov rdi,rsp              ; Set the address of stack as scratch box  
	xor rax, rax             ; Start with 0
init:
	stosb                    ; Store next SBox byte S[i] = i
	inc al                   ; Increase byte to write (EDI is increased automatically)
	jnz init                 ; Loop until we wrap around
	sub rdi, 0x100           ; Restore EDI
	; Permute S-box according to key
	xor rbx, rbx             ; Clear EBX (EAX is already cleared)
permute:
  add bl,[rdi+rax]         ; BL += S[AL] + KEY[AL % sizeof(Key)]
  mov rdx,rax              ; 
  and dl,KSize-1           ; dl & sizeof(Key)
  add bl,[rsi+rdx]         ; Move next byte of key to bl 
  mov dl,[rdi+rax]         ; swap S[AL] and S[BL]
  xchg dl,[rdi+rbx]        ; ..
  mov [rdi+rax], dl        ; ..
  inc al                   ; AL += 1 until we wrap around
  jnz permute              ;
  ; Decryption loop
  xor rbx, rbx             ; Clear EBX (EAX is already cleared)
decrypt:
  inc al                   ; AL += 1
  add bl,[rdi+rax]         ; BL += S[AL]
  mov dl,[rdi+rax]         ; swap S[AL] and S[BL]
  xchg dl,[rdi+rbx]        ;
  mov [rdi+rax], dl        ;
  add dl,[rdi+rbx]         ; DL = S[AL]+S[BL]
  mov dl,[rdi+rdx]         ; DL = S[DL]
  xor [rbp],dl             ; [EBP] ^= DL
  inc rbp                  ; Advance data pointer
  dec rcx                  ; Reduce counter
  jnz decrypt              ; Until finished
  jmp Payload
