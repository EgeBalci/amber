[BITS 64]


resolve_imports:
	push   r12
	push   rbp
	push   rdi
	push   rsi
	push   rbx
	sub    rsp,0x20
	mov    rbp,rcx
	movsxd rax,DWORD [rcx+0x3c]
	mov    eax,DWORD [rcx+rax*1+0x90]
	mov    edx,0x0
	test   eax,eax
	je     loc_140002966 
	mov    eax,eax
	lea    r12,[rcx+rax*1]
	mov    ecx,DWORD [r12+0xc]
	test   ecx,ecx
	jne    loc_14000292f 
	mov    edx,0x1
	jmp    loc_140002966 
loc_1400028cf:
	mov    edx,0x0
	mov    rcx,rdi
	call   get_proc_by_crc
	test   rax,rax
	je     loc_1400028e4 
	mov    QWORD [rsi],rax
loc_1400028e4:
	add    rbx,0x8
	add    rsi,0x8
	mov    r8,QWORD [rbx]
	test   r8,r8
	je     loc_140002922 
loc_1400028f4:
	test   r8,r8
	js     loc_1400028cf 
	lea    rcx,[rbp+r8*1+0x2]
	mov    edx,0x0
	call   calc_crc
	mov    edx,eax
	mov    r8d,0xffffffff
	mov    rcx,rdi
	call   get_proc_by_crc
	test   rax,rax
	je     loc_1400028e4 
	mov    QWORD [rsi],rax
loc_140002920:
	jmp    loc_1400028e4 
loc_140002922:
	add    r12,0x14
	mov    ecx,DWORD [r12+0xc]
	test   ecx,ecx
	je     loc_14000295a 
loc_14000292f:
	mov    ecx,ecx
	add    rcx,rbp
	call   load_module
	mov    rdi,rax
	test   rax,rax
	je     loc_140002961 
	mov    ebx,DWORD [r12]
	add    rbx,rbp
	mov    esi,DWORD [r12+0x10]
	add    rsi,rbp
	mov    r8,QWORD [rbx]
	test   r8,r8
	jne    loc_1400028f4 
	jmp    loc_140002922 
loc_14000295a:
	mov    edx,0x1
	jmp    loc_140002966 
loc_140002961:
	mov    edx,0x0
loc_140002966:
	mov    eax,edx
	add    rsp,0x20
	pop    rbx
	pop    rsi
	pop    rdi
	pop    rbp
	pop    r12
	ret