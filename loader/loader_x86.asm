;#===========================================#
;# x86 Reflective Loader                     #
;# Author: Ege Balcı <egebalci@pm.me>        #
;# Version: 2.0                              #
;#===========================================#

[BITS 32]


  call start              ; Get the address of pre-mapped PE image to stack
  incbin "pemap.bin"      ; Pre-mapped PE image
start:
  cld                     ; Clear direction flags
  pop esi                 ; Get the address of image to esi
  call $+5                ; Push the current EIP to stack
  sub [esp],esi           ; Subtract &PE from EIP and get image_size
  mov eax,[esi+0x3C]      ; Get the offset of "PE" to eax
  mov ebx,[eax+esi+0x34]  ; Get the image base address to ebx
  mov eax,[eax+esi+0x28]  ; Get the address of entry point to eax
  push eax                ; Save the address of entry to stack
  push 0x40               ; PAGE_EXECUTE_READ_WRITE
  push 0x103000           ; MEM_COMMIT | MEM_TOP_DOWN | MEM_RESERVE
  push dword [esp+0xC]    ; dwSize
  push 0                  ; lpAddress
  push 0x2C39DFEC         ; crc32("KERNEL32.DLL", "VirtualAlloc")
  call api_call           ; VirtualAlloc(lpAddress,dwSize,MEM_COMMIT|MEM_TOP_DOWN|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
  push eax                ; Save the new image base to stack
  xor edx,edx             ; Zero out the edx
relocate:
  mov eax,[esi+0x3C]      ; Offset to IMAGE_NT_HEADER ("PE")
  mov ecx,[eax+esi+0xA4]  ; Base relocation table size
  mov eax,[eax+esi+0xA0]  ; Base relocation table RVA
  add eax,esi             ; Base relocation table memory address
  add ecx,eax             ; End of base relocation table
calc_delta:
  mov edi,[esp]           ; Move the new base address to EDI
  sub edi,ebx             ; Delta value
  push dword [eax]        ; Reloc RVA
  push dword [eax+4]      ; Reloc table size
  add eax,8               ; Move to the reloc descriptor
  jmp fix                 ; Start fixing
get_rva:
  cmp ecx,eax             ; Check if the end of the reloc section ?
  jle reloc_fin           ; If yes goto fin
  add esp,8               ; Deallocate old reloc RVA and reloc table size variables
  push dword [eax]        ; Push new reloc RVA
  push dword [eax+4]      ; Push new reloc table size
  add eax,8               ; Move 8 bytes
fix:
  cmp word [esp],8        ; Check if the end of the reloc block
  jz get_rva              ; If yes set the next block RVA
  mov dx,word [eax]       ; Move the reloc desc to dx
  cmp dx,0                ; Check if it is a padding word
  je pass
  and dx,0x0FFF           ; Get the last 12 bits
  add edx,[esp+4]         ; Add block RVA to desc value
  add edx,esi             ; Add the start address of the image
  add dword [edx],edi     ; Add the delta value to calculated absolute address
pass:
  sub dword [esp],2       ; Decrease the index
  add eax,2               ; Move to the next reloc desc.
  xor edx,edx             ; Zero out edx
  jmp fix                 ; Loop
reloc_fin:
  pop eax                 ; Deallocate all vars
  pop eax                 ; ...
  mov eax,[esi+0x3C]      ; Offset to IMAGE_NT_HEADER ("PE")
  mov eax,[eax+esi+0x80]  ; Import table RVA
  add eax,esi             ; Import table memory address (first image import descriptor)
  push eax                ; Save the address of import descriptor to stack
get_modules:
  cmp dword [eax],0       ; Check if the import names table RVA is NULL
  jz complete             ; If yes building process is done
  mov eax,[eax+0x0C]      ; Get RVA of dll name to eax
  add eax,esi             ; Get the dll name address
  call LoadLibraryA       ; Load the library
  mov ebx,eax             ; Move the dll handle to ebx
  mov eax,[esp]           ; Move the address of current _IMPORT_DESCRIPTOR to eax
  call get_procs          ; Resolve all windows API function addresses
  add dword [esp],0x14    ; Move to the next import descriptor
  mov eax,[esp]           ; Set the new import descriptor address to eax
  jmp get_modules
get_procs:
  push ecx                ; Save ecx to stack
  push dword [eax+0x10]   ; Save the current import descriptor IAT RVA
  add [esp],esi           ; Get the IAT memory address
  mov eax,[eax]           ; Set the import names table RVA to eax
  add eax,esi             ; Get the current import descriptor's import names table address
  push eax                ; Save it to stack
resolve:
  cmp dword [eax],0       ; Check if end of the import names table
  jz all_resolved         ; If yes resolving process is done
  mov eax,[eax]           ; Get the RVA of function hint to eax
  cmp eax,0x80000000      ; Check if the high order bit is set
  js name_resolve         ; If high order bit is not set resolve with INT entry
  sub eax,0x80000000      ; Zero out the high bit
  call GetProcAddress     ; Get the API address with hint
  jmp insert_iat          ; Insert the address of API tı IAT
name_resolve:
  add eax,esi             ; Set the address of function hint
  add eax,2               ; Move to function name
  call GetProcAddress     ; Get the function address to eax
insert_iat:
  mov ecx,[esp+4]         ; Move the IAT address to ecx
  mov [ecx],eax           ; Insert the function address to IAT
  add dword [esp],4       ; Increase the import names table index
  add dword [esp+4],4     ; Increase the IAT index
  mov eax,[esp]           ; Set the address of import names table address to eax
  jmp resolve             ; Loop
all_resolved:
  mov ecx,[esp+4]         ; Move the IAT address to ecx
  mov dword [ecx],0       ; Insert a NULL dword
  pop ecx                 ; Deallocate index values
  pop ecx                 ; ...
  pop ecx                 ; Put back the ecx value
  ret                     ; <-
LoadLibraryA:
  push ecx                ; Save ecx to stack
  push edx                ; Save edx to stack
  push eax                ; Push the address of linrary name string
  push 0xE2E6A091         ; crc32( "kernel32.dll", "LoadLibraryA" )
  call api_call           ; LoadLibraryA([esp+4])
  pop edx                 ; Retreive edx
  pop ecx                 ; Retreive ecx
  ret                     ; <-
GetProcAddress:
  push ecx                ; Save ecx to stack
  push edx                ; Save edx to stack
  push eax                ; Push the address of proc name string
  push ebx                ; Push the dll handle
  push 0xA18B0B38         ; crc32( "kernel32.dll", "GetProcAddress" )
  call api_call           ; GetProcAddress(ebx,[esp+4])
  pop edx                 ; Retrieve edx
  pop ecx                 ; Retrieve ecx
  ret                     ; <-
complete:
  pop eax                 ; Clean out the stack
  pop edi                 ; ..
  mov edx,edi             ; Copy the address of new base to EDX
  pop eax                 ; Pop the address_of_entry to EAX
  add edi,eax             ; Add the address of entry to new image base
  pop ecx                 ; Pop the image_size to ECX
memcpy:
  mov al,[esi]            ; Move 1 byte of PE image to AL register
  mov byte [edx],al       ; Move 1 byte of PE image to image base
  mov byte [esi],0        ; Overwrite copied byte (for less memory footprint) 
  inc esi                 ; Increase PE image index
  inc edx                 ; Increase image base index
  loop memcpy             ; Loop until ECX = 0
  jmp PE_start

; ========== API ==========
%include "CRC32_API/x86_crc32_api.asm"

PE_start:
  mov ecx,wipe                    ; Get the number of bytes until wipe label
  call wipe_start                 ; Call wipe_start
wipe_start:
  pop eax                         ; Get EIP to EAX
wipe:
  mov byte [eax],0                ; Wipe 1 byte at a time
  dec eax                         ; Decraise EAX
  loop wipe                       ; Loop until ECX = 0
  jmp edi                         ; Call the AOE