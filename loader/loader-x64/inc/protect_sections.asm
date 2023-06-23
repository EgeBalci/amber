[BITS 64]


protect_sections:
	push   r14
	push   r13
	push   r12
	push   rbp
	push   rdi
	push   rsi
	push   rbx
	sub    rsp,0x50
	mov    rbp,rcx
	movsxd rdi,DWORD [rcx+0x3c]
	add    rdi,rcx
	movzx  eax,WORD [rdi+0x14]
	lea    rbx,[rdi+rax*1+0x18]
	mov    QWORD [rsp+0x48],0x0
	cmp    WORD [rdi+0x6],0x0
	je     loc_140002e87 
	mov    esi,0x0
	mov    r12d,0x0
	lea    r14,[rsp+0x40]
	lea    r13,[rsp+0x48]
	jmp    loc_140002f50 
loc_140002e38:
	mov    ecx,0x1
loc_140002e3d:
	and    eax,0x60000000
	mov    r8d,0x1
	cmp    eax,0x60000000
	mov    r9d,0x20
	mov    eax,0x80
	cmovne r9d,eax
	jmp    loc_140002eed 
loc_140002e61:
	mov    r9d,0x20
	jmp    loc_140002efd 
loc_140002e6c:
	mov    eax,0x0
	jmp    loc_140002e78 
loc_140002e73:
	mov    eax,0x1
loc_140002e78:
	add    rsp,0x50
	pop    rbx
	pop    rsi
	pop    rdi
	pop    rbp
	pop    r12
	pop    r13
	pop    r14
	ret
loc_140002e87:
	mov    eax,0x1
	jmp    loc_140002e78 
loc_140002e8e:
	mov    ecx,0x1
	mov    r8d,r12d
	mov    r9d,0x10
	jmp    loc_140002edd 
loc_140002e9e:
	mov    ecx,r12d
	test   eax,0x20000000
	je     loc_140002eda 
	mov    ecx,0x0
	test   eax,eax
	js     loc_140002e3d 
	mov    ecx,eax
	shr    ecx,0x1f
	mov    r8d,ecx
	mov    ecx,r12d
	mov    r9d,0x10
	jmp    loc_140002edd 
loc_140002ec4:
	test   eax,0x20000000
	jne    loc_140002e38 
	mov    ecx,0x1
	mov    r9d,0x4
loc_140002eda:
	mov    r8d,r12d
loc_140002edd:
	and    eax,0x60000000
	cmp    eax,0x60000000
	je     loc_140002e61 
loc_140002eed:
	test   cl,cl
	je     loc_140002efd 
	test   r8b,r8b
	mov    eax,0x40
	cmovne r9d,eax
loc_140002efd:
	mov    eax,DWORD [rdx+0xc]
	add    rax,rbp
	mov    QWORD [rsp+0x48],rax
	mov    eax,DWORD [rdx+0x10]
	mov    QWORD [rsp+0x40],rax
	mov    DWORD [rsp+0x3c],0x0
	lea    rax,[rsp+0x3c]
	mov    QWORD [rsp+0x20],rax
	mov    r8,r14
	mov    rdx,r13
	mov    rcx,0xffffffffffffffff
	mov    r10, 0x6EDE4D41
	call   api_call
	call   rax        ; <NtProtectVirtualMemory>
	test   eax,eax
	js     loc_140002e6c 
loc_140002f3d:
	add    esi,0x1
	add    rbx,0x28
	movzx  eax,WORD [rdi+0x6]
	cmp    eax,esi
	jle    loc_140002e73 
loc_140002f50:
	mov    rdx,rbx
	mov    eax,DWORD [rbx+0x24]
	test   eax,eax
	je     loc_140002f3d 
	mov    r9d,eax
	sar    r9d,0x1f
	and    r9d,0xffffffc8
	add    r9d,0x40
	test   eax,0x40000000
	je     loc_140002e9e 
	test   eax,eax
	js     loc_140002ec4 
	test   eax,0x20000000
	jne    loc_140002e8e 
	mov    ecx,0x1
	mov    r9d,0x2
	jmp    loc_140002eda 