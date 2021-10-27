package amber

// CRC64 https://github.com/EgeBalci/CRC32_API
const CRC64 = `

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
