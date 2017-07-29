; This block requires following values inside the specified registers...
;
; #########################################
; #  EBX -> Original image base           #
; #  ESI -> Address of PE image		  #
; #  EBP -> &hash_api			  #  				
; #########################################
; 
; Author: Ege Balcı <egebalci[at]protonmail[dot]com> 
; Version: 1.0
;
;######## FUNCTION USAGE #######
;	
;	LoadLibraryA(string dllName); [eax]
;	GetProcs(HANDLE dllHandle, &_IMPORT_DESCRIPTOR); [ebx] [eax]
;	GetProcAddress(HANDLE dllHandle, string funcName)
;	InserAddress();
;

[BITS 32]
[ORG 0]

BuildImportTable:
	mov eax,[esi+0x3C]	; Offset to IMAGE_NT_HEADER ("PE")
	mov eax,[eax+esi+0x80] 	; Import table RVA
	add eax,esi		; Import table memory address (first image import descriptor)
	push eax		; Save import descriptor to stack
GetDLLs:
	cmp dword [eax],0x00	; Check if the import names table RVA is NULL
	jz Complete		; If yes building process is done
	mov eax,[eax+0x0C]	; Get RVA of dll name to eax
	add eax,esi		; Get the dll name address		
	call LoadLibraryA	; Load the library
	mov ebx,eax 		; Move the dll handle to ebx
	mov eax,[esp]		; Move the address of current _IMPORT_DESCRIPTOR to eax 
	call GetProcs		; Resolve all windows API function addresses
	add dword [esp],0x14	; Move to the next import descriptor
	mov eax,[esp]		; Set the new import descriptor address to eax
	jmp GetDLLs
;-----------------------------------------------------------------------------------
GetProcs:
	push ecx 		; Save ecx to stack
	push dword [eax+0x10]	; Save the current import descriptor IAT RVA
	add [esp],esi		; Get the IAT memory address 
	mov eax,[eax]		; Set the import names table RVA to eax
	add eax,esi		; Get the current import descriptor's import names table address	
	push eax		; Save it to stack
Resolve: 
	cmp dword [eax],0x00 	; Check if end of the import names table
	jz AllResolved		; If yes resolving process is done
	mov eax,[eax]		; Get the RVA of function hint to eax
	cmp eax,0x80000000	; Check if the high order bit is set
	js NameResolve		; If high order bit is not set resolve with INT entry
	sub eax,0x80000000	; Zero out the high bit
	call GetProcAddress	; Get the API address with hint
	jmp InsertIAT		; Insert the address of API tı IAT
NameResolve:
	add eax,esi		; Set the address of function hint
	add dword eax,0x02	; Move to function name
	call GetProcAddress	; Get the function address to eax
InsertIAT:
	mov ecx,[esp+4]		; Move the IAT address to ecx 
	mov [ecx],eax		; Insert the function address to IAT
	add dword [esp],0x04	; Increase the import names table index
	add dword [esp+4],0x04	; Increase the IAT index
	mov eax,[esp]		; Set the address of import names table address to eax
	jmp Resolve		; Loop
AllResolved:
	mov ecx,[esp+4]         ; Move the IAT address to ecx 
	mov dword [ecx],0x00	; Insert a NULL dword
	add esp,0x08		; Deallocate index values
	pop ecx			; Put back the ecx value
	ret			; <-
;-----------------------------------------------------------------------------------
LoadLibraryA:
	push ecx		; Save ecx to stack
	push edx		; Save edx to stack
	push eax 		; Push the address of linrary name string
	call [LLA] 		; LoadLibraryA([esp+4])
	pop edx			; Retreive edx
	pop ecx			; Retreive ecx
	ret 			; <-
;-----------------------------------------------------------------------------------
GetProcAddress:
	push ecx 		; Save ecx to stack
	push edx		; Save edx to stack
	push eax		; Push the address of proc name string
	push ebx 		; Push the dll handle
	call [GPA]		; GetProcAddress(ebx,[esp+4])
	pop edx			; Retrieve edx
	pop ecx			; Retrieve ecx
	ret 			; <-
;-----------------------------------------------------------------------------------
Complete:
	pop eax			; Clean out the stack
