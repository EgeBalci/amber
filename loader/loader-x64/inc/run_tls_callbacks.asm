[BITS 64]


run_tls_callbacks:
	push   rsi
	push   rbx
	sub    rsp,0x28
	mov    rsi,rcx
	movsxd rax,DWORD [rcx+0x3c]
	mov    eax,DWORD [rcx+rax*1+0xd0]
	mov    edx,0x0
	test   eax,eax
	je     loc_1400033ad 
	mov    eax,eax
	mov    rbx,QWORD [rcx+rax*1+0x18]
	mov    edx,0x1
	test   rbx,rbx
	jne    loc_1400033ca 
loc_1400033ad:
	mov    eax,edx
	add    rsp,0x28
	pop    rbx
	pop    rsi
	ret
loc_1400033b6:
	mov    r8d,0x0
	mov    edx,0x1
	mov    rcx,rsi
	call   rax
	add    rbx,0x8
loc_1400033ca:
	mov    rax,QWORD [rbx]
	test   rax,rax
	jne    loc_1400033b6 
	mov    edx,0x1
	jmp    loc_1400033ad 