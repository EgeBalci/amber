package amber

const LoaderX64 = `
start:
  pop rsi                         ; Get the address of image to rsi
  call get_ip                     ; Push the current EIP to stack
get_ip:
  cld                             ; Clear direction flags
  sub [rsp],rsi                   ; Subtract the address of pre mapped PE image and get the image_size+8 to ST[0]
  mov rbp,rsp                     ; Copy current stack address to rbp
  and rbp,-0x1000                 ; Create a new shadow stack address
  mov eax,dword [rsi+0x3C]        ; Get the offset of "PE" to eax
  mov rbx,qword [rax+rsi+0x30]    ; Get the image base address to rbx
  mov r12d,dword [rax+rsi+0x28]   ; Get the address of entry point to r12
  mov r9d,0x40                    ; PAGE_EXECUTE_READ_WRITE
  mov r8d,0x00103000              ; MEM_COMMIT | MEM_TOP_DOWN | MEM_RESERVE
  mov rdx,[rsp]                   ; dwSize
  xor rcx,rcx                     ; lpAddress
  xchg rsp,rbp                    ; Swap shadow stack
  mov r10d,0x2C39DFEC             ; crc32("KERNEL32.DLL", "VirtualAlloc")
  call api_call                   ; VirtualAlloc(lpAddress,dwSize,MEM_COMMIT|MEM_TOP_DOWN|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
  xchg rsp,rbp                    ; Swap shadow stack
  mov rdi,rax                     ; Save the new base address to rdi
  xor rax,rax                     ; Zero out the RAX
  xor r8,r8                       ; Zero out the R8
  xor r13,r13                     ; Zero out the R13
  xor r14,r14                     ; Zero out the R14
  mov eax,dword [rsi+0x3C]        ; Offset to IMAGE_NT_HEADER ("PE")
  mov ecx,dword [rax+rsi+0xB4]    ; Base relocation table size
  mov eax,dword [rax+rsi+0xB0]    ; Base relocation table RVA
  add rax,rsi                     ; Base relocation table memory address
  add rcx,rax                     ; End of base relocation table
calc_delta:
  mov rdx,rdi                     ; Move the new base address to rdx
  sub rdx,rbx                     ; Delta value
  mov r13d,dword [rax]            ; Move the reloc RVA to R13D
  mov r14d,dword [rax+4]          ; Move the reloc table size to R14D
  add rax,8                       ; Move to the reloc descriptor
  jmp fix                         ; Start fixing
get_rva:
  cmp rcx,rax                     ; Check if the end of the reloc section
  jle reloc_fin                   ; If yes goto fin
  mov r13d,dword [rax]            ; Move the new reloc RVA
  mov r14d,dword [rax+4]          ; Move the new reloc table size
  add rax,8                       ; Move 8 bytes
fix:
  cmp r14w,8                      ; Check if the end of the reloc block
  jz get_rva                      ; If yes set the next block RVA
  mov r8w,word [rax]              ; Move the reloc desc to r8w
  cmp r8w,0                       ; Check if it is a padding word
  je pass                         ; Pass padding bytes
  and r8w,0x0FFF                  ; Get the last 12 bits
  add r8d,r13d                    ; Add block RVA to desc value
  add r8,rsi                      ; Add the start address of the image
  add [r8],rdx                    ; Add the delta value to calculated absolute address
pass:
  sub r14d,2                      ; Decrease the index
  add rax,2                       ; Move to the next reloc desc.
  xor r8,r8                       ; Zero out r8
  jmp fix                         ; Loop
reloc_fin:                        ; All done !
  xor r14,r14                     ; Zero out r14
  xor r15,r15                     ; Zero out r15
  xor rcx,rcx                     ; Zero out rcx
  mov eax,dword [rsi+0x3C]        ; Offset to IMAGE_NT_HEADER ("PE")
  mov eax,dword [rax+rsi+0x90]    ; Import table RVA
  add rax,rsi                     ; Import table memory address (first image import descriptor)
  push rax                        ; Save import descriptor to stack
get_modules:
  cmp dword [rax],0               ; Check if the import names table RVA is NULL
  jz complete                     ; If yes building process is done
  mov ecx,dword [rax+0x0C]        ; Get RVA of dll name to eax
  add rcx,rsi                     ; Get the dll name address
  call LoadLibraryA               ; Load the library
  mov r13,rax                     ; Move the dll handle to R13
  mov rax,[rsp]                   ; Move the address of current _IMPORT_DESCRIPTOR to eax
  call get_procs                  ; Resolve all windows API function addresses
  add dword [rsp],0x14            ; Move to the next import descriptor
  mov rax,[rsp]                   ; Set the new import descriptor address to eax
  jmp get_modules                 ; Get other modules
get_procs:
  mov r14d,dword [rax+0x10]       ; Save the current import descriptor IAT RVA
  add r14,rsi                     ; Get the IAT memory address
  mov rax,[rax]                   ; Set the import names table RVA to eax
  add rax,rsi                     ; Get the current import descriptor's import names table address
  mov r15,rax                     ; Save &INT to R15
resolve:
  cmp dword [rax],0               ; Check if end of the import names table
  jz all_resolved                 ; If yes resolving process is done
  mov rax,[rax]                   ; Get the RVA of function hint to eax
  btr rax,0x3F                    ; Check if the high order bit is set
  jnc name_resolve                ; If high order bit is not set resolve with INT entry
  shl rax,2                       ; Discard the high bit by shifting
  shr rax,2                       ; Shift back the original value
  call GetProcAddress             ; Get the API address with hint
  jmp insert_iat                  ; Insert the address of API tı IAT
name_resolve:
  add rax,rsi                     ; Set the address of function hint
  add rax,2                       ; Move to function name
  call GetProcAddress             ; Get the function address to eax
insert_iat:
  mov [r14],rax                   ; Insert the function address to IAT
  add r14,8                       ; Increase the IAT index
  add r15,8                       ; Increase the import names table index
  mov rax,r15                     ; Set the address of import names table address to eax
  jmp resolve                     ; Loop
all_resolved:
  mov qword [r14],0               ; Insert a NULL qword
  ret                             ; <-
LoadLibraryA:
  xchg rbp,rsp                    ; Swap shadow stack
  mov r10d,0xE2E6A091             ; crc32("KERNEL32.DLL", "LoadLibraryA")
  call api_call                   ; LoadLibraryA(RCX)
  xchg rbp,rsp                    ; Swap shadow stack
  ret                             ; <-
GetProcAddress:
  xchg rbp,rsp                    ; Swap shadow stack
  mov rcx,r13                     ; Move the module handle to RCX as first parameter
  mov rdx,rax                     ; Move the address of function name string to RDX as second parameter
  mov r10d,0xA18B0B38             ; crc32("KERNEL32.DLL", "GetProcAddress")
  call api_call                   ; GetProcAddress(ebx,[esp+4])
  xchg rbp,rsp                    ; Swap shadow stack
  ret                             ; <-
complete:
  pop rax                         ; Clean out the stack
  pop rcx                         ; Pop the ImageSize into RCX
  push rdi                        ; Save ImageBase to stack
  mov r13,rdi                     ; Copy the new base value to r13
  add r13,r12                     ; Add the address of entry value to new base address
memcpy:
  mov al,[rsi]                    ; Move 1 byte of PE image to AL register
  mov [rdi],al                    ; Move 1 byte of PE image to image base
  mov byte [rsi],0                ; Overwrite copied byte (for less memory footprint)
  inc rsi                         ; Increase PE image index
  inc rdi                         ; Increase image base index
  loop memcpy                     ; Loop until zero
PE_start:
  pop r13                         ; Pop the image base to r13
  or rcx,-1                       ; hProcess
  xor rdx,rdx                     ; lpBaseAddress
  xor r8,r8                       ; hProcess
  xchg rbp,rsp                    ; Swap shadow stack
  mov r10d,0x975B539E             ; crc32("KERNEL32.dll", "FlushInstructionCache")
  call api_call                   ; FlushInstructionCache(0xffffffff,NULL,NULL);
  xchg rbp,rsp                    ; Swap shadow stack
  add r13,r12                     ; Add the address of entry value to image base
  jmp r13                         ; Call the AOE

`

const LoaderX86 = `
start:
  cld                     ; Clear direction flags
  pop esi                 ; Get the address of image to esi
  call get_ip             ; Push the current EIP to stack
get_ip:
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
  mov [edx],al            ; Move 1 byte of PE image to image base
  inc esi                 ; Increase PE image index
  inc edx                 ; Increase image base index
  loop memcpy             ; Loop until ECX = 0
PE_Start:
  jmp edi                ; Call PE AOE

`

const FixedLoaderX64 = `
start:
	pop rsi                         ; Get the address of image to rsi
  call get_ip                     ; Push the current EIP to stack
get_ip:
	sub [rsp],rsi                   ; Subtract the address of pre mapped PE image and get the image_size to R11
	mov rbp,rsp                     ; Copy current stack address to rbp
	and rbp,-0x1000                 ; Create a new shadow stack address
	mov eax,dword [rsi+0x3C]        ; Get the offset of "PE" to eax
	mov rbx,qword [rax+rsi+0x30]    ; Get the image base address to rbx
	mov r12d,dword [rax+rsi+0x28]   ; Get the address of entry point to r12
	push rax                        ; Allocate 8 bytes for lpflOldProtect
	mov r9,rsp                      ; lpflOldProtect
	mov r8d,dword 0x40              ; PAGE_EXECUTE_READWRITE
	mov rdx,qword [rsp+8]           ; dwSize
	mov rcx,rbx                     ; lpAddress
	mov r10d,0x80886EF1             ; crc32( "kernel32.dll", "VirtualProtect" )
	xchg rsp,rbp                    ; Swap shadow stack
	call api_call                   ; VirtualProtect( image_base, image_size, PAGE_EXECUTE_READWRITE, lpflOldProtect)
	xchg rsp,rbp                    ; Swap shadow stack
	xor rax,rax                     ; Zero EAX
	xor r14,r14                     ; Zero R14
	xor r15,r15                     ; Zero R15
	mov eax,dword [rsi+0x3C]        ; Offset to IMAGE_NT_HEADER ("PE")
	mov eax,dword [rax+rsi+0x90]    ; Import table RVA
	add rax,rsi                     ; Import table memory address (first image import descriptor)
	push rax                        ; Save import descriptor to stack
get_modules:
	cmp dword [rax],0               ; Check if the import names table RVA is NULL
	jz complete                     ; If yes building process is done
	mov ecx,dword [rax+0x0C]        ; Get RVA of dll name to eax
	add rcx,rsi                     ; Get the dll name address
	call LoadLibraryA               ; Load the library
	mov r13,rax                     ; Move the dll handle to R13
	mov rax,[rsp]                   ; Move the address of current _IMPORT_DESCRIPTOR to eax 
	call get_procs                  ; Resolve all windows API function addresses
	add dword [rsp],0x14            ; Move to the next import descriptor
	mov rax,[rsp]                   ; Set the new import descriptor address to RAX
	jmp get_modules
get_procs:
	mov r14d,dword [rax+0x10]       ; Save the current import descriptor IAT RVA to R14D
	add r14,rsi                     ; Get the IAT memory address 
	mov rax,[rax]                   ; Set the import names table RVA to RAX
	add rax,rsi                     ; Get the current import descriptor's import names table address	
	mov r15,rax                     ; Save &INT to R15
resolve: 
	cmp dword [rax],0               ; Check if end of the import names table
	jz all_resolved                 ; If yes resolving stage is done
	mov rax,[rax]                   ; Get the RVA of function hint to eax
	btr rax,0x3F                    ; Check if the high order bit is set
	jnc name_resolve                ; If high order bit is not set resolve with INT entry
	shl rax,2                       ; Discard the high bit by shifting
	shr rax,2                       ; Shift back the original value
	call GetProcAddress             ; Get API address with hint
	jmp insert_iat                  ; Insert the address of API to IAT
name_resolve:
	add rax,rsi                     ; Set the address of function hint
	add rax,2                       ; Move to function name
	call GetProcAddress             ; Get the function address to eax
insert_iat: 
	mov [r14],rax                   ; Insert the function address to IAT
	add r14,8                       ; Increase the IAT index
	add r15,8                       ; Increase the import names table index
	mov rax,r15                     ; Set the address of import names table address to RAX
	jmp resolve                     ; Loop
all_resolved:
	mov qword [r14],0               ; Insert a NULL qword
	ret                             ; <-
LoadLibraryA:
	xchg rbp,rsp                    ; Swap shadow stack
	mov r10d,0xE2E6A091             ; hash( "kernel32.dll", "LoadLibraryA" )
	call api_call                   ; LoadLibraryA(RCX)
	xchg rbp,rsp                    ; Swap shadow stack
	ret                             ; <-
GetProcAddress:
	xchg rbp,rsp                    ; Swap shadow stack
	mov rcx,r13                     ; Move the module handle to RCX as first parameter
	mov rdx,rax                     ; Move the address of function name string to RDX as second parameter
	mov r10d,0xA18B0B38             ; hash( "kernel32.dll", "GetProcAddress" )
	call api_call                   ; GetProcAddress(RCX,RDX)
	xchg rbp,rsp                    ; Swap shadow stack
	ret                             ; <-
complete:
	pop rax                         ; Clean out the stack
	pop rax                         ; ...
	pop rcx                         ; Pop the image_size to RCX
	push rbx                        ; Push the new base adress to stack
	add [rsp],r12                   ; Add the address of entry value to new base address
memcpy:	
	mov al,[rsi]                    ; Move 1 byte of PE image to AL register
	mov [rbx],al                    ; Move 1 byte of PE image to image base
	inc rsi                         ; Increase PE image index
	inc rbx                         ; Increase image base index
	loop memcpy                     ; Loop until zero
	ret                             ; <-

`

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
  push byte 0x40          ; PAGE_EXECUTE_READWRITE
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
  mov [edx],al            ; Move 1 byte of PE image to image base
  inc esi                 ; Increase PE image index
  inc edx                 ; Increase image base index
  loop memcpy             ; Loop until ECX = 0
PE_Start:
  jmp edi                ; Call PE AOE

`

const CRC_API_64 = `

api_call:
  push r9                  ; Save the 4th parameter
  push r8                  ; Save the 3rd parameter
  push rdx                 ; Save the 2nd parameter
  push rcx                 ; Save the 1st parameter
  push rsi                 ; Save RSI
  xor rdx, rdx             ; Zero rdx
  mov rdx, gs:[rdx+0x60]   ; Get a pointer to the PEB
  mov rdx, [rdx+0x18]      ; Get PEB->Ldr
  mov rdx, [rdx+0x20]      ; Get the first module from the InMemoryOrder module list
next_mod:                  ;
  mov rsi, [rdx+0x50]      ; Get pointer to modules name (unicode string)
  movzx rcx, word [rdx+0x4A]; Set rcx to the length we want to check 
  xor r9, r9               ; Clear r9 which will store the hash of the module name
loop_modname:              ;
  xor rax, rax             ; Clear rax
  lodsb                    ; Read in the next byte of the name
  cmp al, 'a'              ; Some versions of Windows use lower case module names
  jl not_lowercase         ;
  sub al, 0x20             ; If so normalise to uppercase
not_lowercase:             ;
  crc32 r9d,al             ; Calculate CRC3
  loop loop_modname        ; Loop untill we have read enough
  ; We now have the module hash computed
  push rdx                 ; Save the current position in the module list for later
  push r9                  ; Save the current module hash for later
  ; Proceed to itterate the export address table, 
  mov rdx, [rdx+0x20]      ; Get this modules base address
  mov eax, dword [rdx+0x3C]; Get PE header
  add rax, rdx             ; Add the modules base address
  cmp word [rax+0x18],0x020B ; is this module actually a PE64 executable? 
  ; this test case covers when running on wow64 but in a native x64 context via nativex64.asm and 
  ; their may be a PE32 module present in the PEB's module list, (typicaly the main module).
  ; as we are using the win64 PEB ([gs:96]) we wont see the wow64 modules present in the win32 PEB ([fs:48])
  jne get_next_mod1        ; if not, proceed to the next module
  mov eax, dword [rax+0x88] ; Get export tables RVA
  test rax, rax            ; Test if no export address table is present
  jz get_next_mod1         ; If no EAT present, process the next module
  add rax, rdx             ; Add the modules base address
  push rax                 ; Save the current modules EAT
  mov ecx, dword [rax+0x18]; Get the number of function names  
  mov r8d, dword [rax+0x20]; Get the rva of the function names
  add r8, rdx              ; Add the modules base address
  ; Computing the module hash + function hash
get_next_func:             ;
  jrcxz get_next_mod       ; When we reach the start of the EAT (we search backwards), process the next module
  mov r9, [rsp+8]          ; Reset the current module hash
  dec rcx                  ; Decrement the function name counter
  mov esi, dword [r8+rcx*4]; Get rva of next module name
  add rsi, rdx             ; Add the modules base address
  ; And compare it to the one we want
loop_funcname:             ;
  xor rax, rax             ; Clear rax
  lodsb                    ; Read in the next byte of the ASCII function name
  crc32 r9d,al             ; Calculate CRC32
  cmp al, ah               ; Compare AL (the next byte from the name) to AH (null)
  jne loop_funcname        ; If we have not reached the null terminator, continue
  cmp r9d, r10d            ; Compare the hash to the one we are searchnig for 
  jnz get_next_func        ; Go compute the next function hash if we have not found it
  ; If found, fix up stack, call the function and then value else compute the next one...
  pop rax                  ; Restore the current modules EAT
  mov r8d, dword [rax+0x24]; Get the ordinal table rva      
  add r8, rdx              ; Add the modules base address
  mov cx, [r8+2*rcx]       ; Get the desired functions ordinal
  mov r8d, dword [rax+0x1C]; Get the function addresses table rva  
  add r8, rdx              ; Add the modules base address
  mov eax, dword [r8+4*rcx]; Get the desired functions RVA
  add rax, rdx             ; Add the modules base address to get the functions actual VA
  ; We now fix up the stack and perform the call to the drsired function...
finish:
  pop r8                   ; Clear off the current modules hash
  pop r8                   ; Clear off the current position in the module list
  pop rsi                  ; Restore RSI
  pop rcx                  ; Restore the 1st parameter
  pop rdx                  ; Restore the 2nd parameter
  pop r8                   ; Restore the 3rd parameter
  pop r9                   ; Restore the 4th parameter
  pop r10                  ; Pop off the return address
  sub rsp, 0x20            ; Reserve space for the four register params (4 * sizeof(QWORD) = 32)
                           ; It is the callers responsibility to restore RSP if need be (or alloc more space or align RSP).
  push r10                 ; Push back the return address
  jmp rax                  ; Jump to required function
  ; We now automagically return to the correct caller...
get_next_mod:              ;
  pop rax                  ; Pop off the current (now the previous) modules EAT
get_next_mod1:             ;
  pop r9                   ; Pop off the current (now the previous) modules hash
  pop rdx                  ; Restore our position in the module list
  mov rdx, [rdx]           ; Get the next module
  jmp next_mod             ; Process this module

`

const CRC_API_32 = `

api_call:
  pushad                 ; We preserve all the registers for the caller, bar EAX and ECX.
  mov ebp, esp           ; Create a new stack frame
  xor eax, eax           ; Zero EAX (upper 3 bytes will remain zero until function is found)
  mov edx, fs:[eax+0x30] ; Get a pointer to the PEB
  mov edx, [edx+0xC]     ; Get PEB->Ldr
  mov edx, [edx+0x14]    ; Get the first module from the InMemoryOrder module list
next_mod:                ;
  mov esi, [edx+0x28]    ; Get pointer to modules name (unicode string)
  movzx ecx, word [edx+0x26] ; Set ECX to the length we want to check
  xor edi, edi           ; Clear EDI which will store the hash of the module name
loop_modname:            ;
  lodsb                  ; Read in the next byte of the name
  cmp al, 'a'            ; Some versions of Windows use lower case module names
  jl not_lowercase       ;
  sub al, 0x20           ; If so normalise to uppercase
not_lowercase:           ;
  crc32 edi,al           ; Calculate CRC32 value
  loop loop_modname      ; Loop until we have read enough

  ; We now have the module hash computed
  push edx               ; Save the current position in the module list for later
  push edi               ; Save the current module hash for later
  ; Proceed to iterate the export address table,
  mov edx, [edx+0x10]    ; Get this modules base address
  mov ecx, [edx+0x3C]    ; Get PE header

  ; use ecx as our EAT pointer here so we can take advantage of jecxz.
  mov ecx, [ecx+edx+0x78] ; Get the EAT from the PE header
  jecxz get_next_mod1    ; If no EAT present, process the next module
  add ecx, edx           ; Add the modules base address
  push ecx               ; Save the current modules EAT
  mov ebx, [ecx+0x20]    ; Get the rva of the function names
  add ebx, edx           ; Add the modules base address
  mov ecx, [ecx+0x18]    ; Get the number of function names
  ; now ecx returns to its regularly scheduled counter duties

  ; Computing the module hash + function hash
get_next_func:           ;
  jecxz get_next_mod     ; When we reach the start of the EAT (we search backwards), process the next module
  mov edi, [ebp-8]       ; Reset the current module hash
  dec ecx                ; Decrement the function name counter
  mov esi, [ebx+ecx*4]   ; Get rva of next module name
  add esi, edx           ; Add the modules base address
  ; And compare it to the one we want
loop_funcname:           ;
  lodsb                  ; Read in the next byte of the ASCII function name
  crc32 edi,al           ; Calculate CRC32
  cmp al, ah             ; Compare AL (the next byte from the name) to AH (null)
  jne loop_funcname      ; If we have not reached the null terminator, continue
  cmp edi, [ebp+0x24]    ; Compare the hash to the one we are searching for
  jnz get_next_func      ; Go compute the next function hash if we have not found it

  ; If found, fix up stack, call the function and then value else compute the next one...
  pop eax                ; Restore the current modules EAT
  mov ebx, [eax+0x24]    ; Get the ordinal table rva
  add ebx, edx           ; Add the modules base address
  mov cx, [ebx+2*ecx]    ; Get the desired functions ordinal
  mov ebx, [eax+0x1C]    ; Get the function addresses table rva
  add ebx, edx           ; Add the modules base address
  mov eax, [ebx+4*ecx]   ; Get the desired functions RVA
  add eax, edx           ; Add the modules base address to get the functions actual VA
  ; We now fix up the stack and perform the call to the desired function...
finish:
  mov [esp+0x24], eax    ; Overwrite the old EAX value with the desired api address for the upcoming popad
  pop ebx                ; Clear off the current modules hash
  pop ebx                ; Clear off the current position in the module list
  popad                  ; Restore all of the callers registers, bar EAX, ECX and EDX which are clobbered
  pop ecx                ; Pop off the origional return address our caller will have pushed
  pop edx                ; Pop off the hash value our caller will have pushed
  push ecx               ; Puh back the return address
  jmp eax                ; Properly call the required function for EAF bypass
  ret
  ; We now automagically return to the correct caller...

get_next_mod:            ;
  pop edi                ; Pop off the current (now the previous) modules EAT
get_next_mod1:           ;
  pop edi                ; Pop off the current (now the previous) modules hash
  pop edx                ; Restore our position in the module list
  mov edx, [edx]         ; Get the next module
  jmp next_mod     ; Process this module

`

const IAT_API_32 = `

api_call:
	pushad                  ; We preserve all the registers for the caller, bar EAX and ECX.
	xor eax,eax             ; Zero EAX (upper 3 bytes will remain zero until function is found)
	mov edx,fs:[eax+0x30]   ; Get a pointer to the PEB
	mov edx,[edx+0x0C]      ; Get PEB->Ldr
	mov edx,[edx+0x14]      ; Get the first module from the InMemoryOrder module list
	mov edx,[edx+0x10]      ; Get this modules base address
	push edx                ; Save the image base to stack (will use this alot)
	add edx,[edx+0x3C]      ; "PE" Header
	mov edx,[edx+0x80]      ; Import table RVA
	add edx,[esp]           ; Address of Import Table
	push edx                ; Save the &IT to stack (will use this alot)	
	mov esi,[esp+4]         ; Move the image base to ESI
	sub esp,0x08            ; Allocate space for import descriptor counter & hash
	sub edx,0x14            ; Prepare the import descriptor pointer for processing
next_desc:
	add edx,0x14            ; Get the next import descriptor
	cmp dword [edx],0x00    ; Check if import descriptor valid
	jz not_found            ; If import name array RVA is zero finish parsing
	mov esi,[esp+0x08]      ; Move the import table address to esi
	mov si,[edx+0x0C]       ; Get pointer to module name string RVA
	xor edi,edi             ; Clear EDI which will store the hash of the module name
loop_modname:               ;
	lodsb                   ; Read in the next byte of the name
	cmp al,'a'              ; Some versions of Windows use lower case module names
	jl not_lowercase        ;
	sub al, 0x20            ; If so normalise to uppercase
not_lowercase:              ;
	crc32 edi,al            ; Calculate CRC32 of module name
	crc32 edi,ah            ; Add NULL for unicode effect
	test al,al              ; Check if we read all
	jnz loop_modname
	; We now have the module hash computed
	mov [esp+4],edx         ; Save the current position in the module list for later
	mov [esp],edi           ; Save the current module hash for later
	; Proceed to iterate the export address table,
	mov ecx,[edx]           ; Get the RVA of import names table 
	add ecx,[esp+0x0C]      ; Add image base and get address of import names table
	sub ecx,0x04            ; Go 4 byte back	
get_next_func:
	mov edi,dword [esp]
	; use ecx as our EAT pointer here so we can take advantage of jecxz.
	add ecx,0x04            ; 4 byte forward
	cmp dword [ecx],0x00    ; Check if end of INT
	jz next_desc            ; If no INT present, process the next import descriptor
	mov esi,[ecx]           ; Get the RVA of func name hint
	cmp esi,0x80000000      ; Check if the high order bit is set
	jns get_next_func       ; If not, there is no function name string :(
	add esi,[esp+0x0C]      ; Add the image base and get the address of function hint
	add dword esi,0x02      ; Move 2 bytes forward to asci function name
	; now ecx returns to its regularly scheduled counter duties
	; Computing the module hash + function hash
	; And compare it to the one we want
loop_funcname:              ;
	lodsb                   ; Read in the next byte of the ASCII function name
	crc32 edi,al            ; Calculate CRC32 of function name
	test al,al              ; Check if AL == 0
	jnz loop_funcname       ; If we have not reached the null terminator, continue
	cmp edi,[esp+0x34]      ; Compare the hash to the one we are searching for
	jnz get_next_func       ; Go compute the next function hash if we have not found it
	; If found, fix up stack, call the function and then value else compute the next one...
	mov eax,[edx+0x10]      ; Get the RVA of current descriptor's IAT 
	mov edx,[edx]           ; Get the import names table RVA of current import descriptor
	add edx,[esp+0x0C]      ; Get the address of import names table of current import descriptor
	sub ecx,edx             ; Find the function array index ?
	add eax,[esp+0x0C]      ; Add the image base to current descriptors IAT RVA
	add eax,ecx             ; Add the function index
	; Now we clean the stack	
	; We now fix up the stack and perform the call to the desired function...
finish:
	mov [esp+0x2C],eax      ; Overwrite the old EAX value with the desired api address for the upcoming popad
	add esp,0x10            ; Deallocate saved module hash, import descriptor address and import table address
	popad                   ; Restore all of the callers registers, bar EAX, ECX and EDX which are clobbered
	pop ebx                 ; Pop off the origional return address our caller will have pushed
	pop edx                 ; Pop off the hash value our caller will have pushed
	push ebx                ; Push back the return address
	mov eax,[eax]           ; Get the address of the desired API
	jmp eax                 ; Jump to target function
	; We now automagically return to the correct caller...
not_found:
	add esp,0x0F            ; Fix the stack
	popad                   ; Restore all registers
	ret                     ; Return

`

const IAT_API_64 = `

api_call:
	push r9                 ; Save the 4th parameter
	push r8                 ; Save the 3rd parameter
	push rdx                ; Save the 2nd parameter
	push rcx                ; Save the 1st parameter
	push rsi                ; Save RSI
	xor rdx,rdx             ; Zero rdx
 	mov rdx,gs:[rdx+0x60]   ; Get a pointer to the PEB
	mov rdx,[rdx+0x18]      ; Get PEB->Ldr
	mov rdx,[rdx+0x20]      ; Get the first module from the InMemoryOrder module list
	mov rdx,[rdx+0x20]      ; Get this modules base address
	push rdx                ; Save the image base to stack (will use this alot)
	add dx,word [rdx+0x3C]  ; "PE" Header
	mov edx,dword [rdx+0x90]; Import table RVA
	add rdx,[rsp]           ; Address of Import Table
	push rdx                ; Save the &IT to stack (will use this alot)
 	mov rsi,[rsp+8]         ; Move the image base to RSI
	sub rsp,0x10            ; Allocate space for import descriptor counter & hash
	sub rdx,0x14            ; Prepare import descriptor pointer for processing
next_desc:
	add rdx,0x14            ; Get the next import descriptor
	cmp dword [rdx],0       ; Check if import descriptor is valid
	jz not_found            ; If import name array RVA is zero finish parsing
	mov rsi,[rsp+0x10]      ; Move import table address to RSI
	mov si,[rdx+0xC]        ; Get pointer to module name string RVA
	xor rdi,rdi	            ; Clear RDI which will store the hash of the module name
loop_modname:
	xor rax,rax             ; Clear RAX for calculating the hash
	lodsb                   ; Read in the next byte of the name
	cmp al,'a'              ; Some versions of windows use lower case module names
	jl not_lowercase        ;
	sub al,0x20             ; If so normalize to uppercase 
not_lowercase:
	crc32 edi,al            ; Calculate CRC32 of module name
	crc32 edi,ah            ; Feed NULL for unicode effect
	test al,al              ; Check if end of the module name
	jnz loop_modname        ; 
	; We now have the module hash computed
	mov [rsp+8],rdx         ; Save the current position in the module listfor later
	mov [rsp],edi           ; Save the current module hash for later
	; Proceed to itterate the export address table, 
	mov ecx,dword [rdx]     ; Get RVA of import names table
	add rcx,[rsp+0x18]      ; Add the image base and get the address of import names table
	sub rcx,8               ; Go 4 bytes back
get_next_func:              ;
	mov rdi,[rsp]           ; Reset module hash
	add rcx,8               ; 8 byte forward
	cmp dword [rcx],0       ; Check if end of INT 
	jz next_desc            ; If no INT present, process the next import descriptor
	mov esi,dword [rcx]     ; Get the RVA of func name hint
	cmp esi,0x80000000      ; Check if the high order bit is set
	jns get_next_func       ; If not, there is no function name string :(
	add rsi,[rsp+0x18]      ; Add the image base and get the address of function name hint
	add rsi,2               ; Move 2 bytes forward to asci function name
	; now ecx returns to its regularly scheduled counter duties
	; Computing the module hash + function hash
	; And compare it to the one we want
loop_funcname:
	xor rax,rax             ; Clear RAX
	lodsb                   ; Read in the next byte of the ASCII function name
	crc32 edi,al            ; Calculate CRC32 of the function name
	cmp al,ah               ; Compare AL (the next byte from the name) to AH (null)
	jne loop_funcname       ; If we have not reached the null terminator, continue
	cmp edi,r10d            ; Compare the hash to the one we are searchnig for 
	jnz get_next_func       ; Go compute the next function hash if we have not found it
	; If found, fix up stack, call the function and then value else compute the next one...
	mov eax,dword [rdx+0x10]; Get the RVA of current descriptor's IAT
	mov edx,dword [rdx]     ; Get the import names table RVA of current import descriptor
	add rdx,[rsp+0x18]      ; Get the address of import names table of current import descriptor
	sub rcx,rdx             ; Find the function array index ?
	add rax,[rsp+0x18]      ; Add the image base to current descriptors IAT RVA
	add rax,rcx             ; Add the function index
	; Now clean the stack
	; We now fix up the stack and perform the call to the drsired function...
finish:
	pop r8                  ; Clear off the current modules hash
	pop r8                  ; Clear off the current position in the module list
	pop r8                  ; Clear off the import table address of last module
	pop r8                  ; Clear off the image base address of last module
	pop rsi                 ; Restore RSI
	pop rcx                 ; Restore the 1st parameter
	pop rdx                 ; Restore the 2nd parameter
	pop r8                  ; Restore the 3rd parameter
	pop r9                  ; Restore the 4th parameter
	pop r10                 ; Pop off the return address
	sub rsp,0x20            ; reserve space for the four register params (4 * sizeof(QWORD) = 32)
                            ; It is the callers responsibility to restore RSP if need be (or alloc more space or align RSP).
	push r10                ; Push back the return address
	mov rax,[rax]           ; Get the address of the desired API
	jmp rax                 ; Jump to target function
	; We now automagically return to the correct caller...
not_found:
	add rsp,0x48            ; Clean out the stack
	ret                     ; Return to caller

`
