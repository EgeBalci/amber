package amber

// FixedLoaderX86 contains the 64 bit PE loader for non-relocatable PE files
const FixedLoaderX86 = `
start:                    ;
  cld                     ; Clear direction flags
  pop esi                 ; Get the address of image to esi
  call get_ip             ; Push the current EIP to stack
get_ip:
  sub [esp],esi           ; Subtract &PE from EIP and get image_size
  mov eax,[esi+0x3C]      ; Get the offset of "PE" to eax
  mov ebx,[eax+esi+0x34]  ; Get the image base address to ebx
  mov eax,[eax+esi+0x28]  ; Get the address of entry point to eax
  push eax                ; Save the address of entry to stack
  push ebx                ; Save image base to stack
  push 0x00000000         ; Allocate a DWORD variable inside stack
  push esp                ; lpflOldProtect
  push 0x40               ; PAGE_EXECUTE_READWRITE
  push dword [esp+0x14]   ; dwSize
  push ebx                ; lpAddress
  push 0x80886EF1         ; crc32( "kernel32.dll", "VirtualProtect" )
  call api_call           ; VirtualProtect( ImageBase, image_size, PAGE_EXECUTE_READWRITE, lpflOldProtect)
  pop eax                 ; Fix the stack
  mov eax,[esi+0x3C]      ; Offset to IMAGE_NT_HEADER ("PE")
  mov eax,[eax+esi+0x80]  ; Import table RVA
  add eax,esi             ; Import table memory address (first image import descriptor)
  push eax                ; Save the address of import descriptor to stack
get_modules:
  cmp dword [eax],0x00    ; Check if the import names table RVA is NULL
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
  cmp dword [eax],0x00    ; Check if end of the import names table
  jz all_resolved         ; If yes resolving process is done
  mov eax,[eax]           ; Get the RVA of function hint to eax
  cmp eax,0x80000000      ; Check if the high order bit is set
  js name_resolve         ; If high order bit is not set resolve with INT entry
  sub eax,0x80000000      ; Zero out the high bit
  call GetProcAddress     ; Get the API address with hint
  jmp insert_iat          ; Insert the address of API tı IAT
name_resolve:
  add eax,esi             ; Set the address of function hint
  add eax,0x02            ; Move to function name
  call GetProcAddress     ; Get the function address to eax
insert_iat:
  mov ecx,[esp+4]         ; Move the IAT address to ecx
  mov [ecx],eax           ; Insert the function address to IAT
  add dword [esp],0x04    ; Increase the import names table index
  add dword [esp+4],0x04  ; Increase the IAT index
  mov eax,[esp]           ; Set the address of import names table address to eax
  jmp resolve             ; Loop
all_resolved:
  mov ecx,[esp+4]         ; Move the IAT address to ecx
  mov dword [ecx],0x00    ; Insert a NULL dword
  pop ecx                 ; Deallocate index values
  pop ecx                 ; ...
  pop ecx                 ; Put back the ecx value
  ret                     ; <-
LoadLibraryA:
  push ecx                ; Save ecx to stack
  push edx                ; Save edx to stack
  push eax                ; Push the address of linrary name string
  push 0xE2E6A091         ; ror13( "kernel32.dll", "LoadLibraryA" )
  call api_call           ; LoadLibraryA([esp+4])
  pop edx                 ; Retreive edx
  pop ecx                 ; Retreive ecx
  ret                     ; <-
GetProcAddress:
  push ecx                ; Save ecx to stack
  push edx                ; Save edx to stack
  push eax                ; Push the address of proc name string
  push ebx                ; Push the dll handle
  push 0xA18B0B38         ; ror13( "kernel32.dll", "GetProcAddress" )
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
  jmp PE_start            ; Wipe artifacts from memory and start PE
%s
PE_start:
%s
  mov ecx,wipe             ; Get the number of bytes until wipe label
  call wipe_start          ; Call wipe_start
wipe_start:
  pop eax                  ; Get EIP to EAX
wipe:
  mov byte [eax],0         ; Wipe 1 byte at a time
  dec eax                  ; Decraise EAX
  loop wipe                ; Loop until ECX = 0
  jmp edi                  ; Call the AOE

`
