[BITS 64]


get_proc_by_crc:
	push   r15
	push   r14
	push   r13
	push   r12
	push   rbp
	push   rdi
	push   rsi
	push   rbx
	sub    rsp,0x258
	mov    rbx,rcx
	mov    r13d,edx
	mov    ebp,r8d
	movsxd rax,DWORD [rcx+0x3c]
	add    rax,rcx
	mov    esi,DWORD [rax+0x88]
	add    rsi,rcx
	mov    eax,DWORD [rax+0x8c]
	mov    DWORD [rsp+0x2c],eax
	mov    r12d,DWORD [rsi+0x20]
	mov    r14d,DWORD [rsi+0x1c]
	mov    r15d,DWORD [rsi+0x24]
	mov    eax,DWORD [rsi+0x18]
	test   eax,eax
	je     loc_140003764 
	mov    eax,eax
	mov    QWORD [rsp+0x20],rax
	mov    edi,0x0
	add    r12,rcx
loc_1400035dc:
	mov    ecx,DWORD [r12+rdi*4]
	add    rcx,rbx
	mov    edx,0x0
	call   calc_crc
	cmp    ebp,edi
	je     loc_14000360b 
	cmp    eax,r13d
	je     loc_14000360b 
	add    rdi,0x1
	cmp    QWORD [rsp+0x20],rdi
	jne    loc_1400035dc 
	mov    eax,0x0
	jmp    loc_140003738 
loc_14000360b:
	lea    rax,[rbx+rdi*2]
	movzx  eax,WORD [rax+r15*1]
	lea    rax,[rbx+rax*4]
	mov    eax,DWORD [rax+r14*1]
	add    rbx,rax
	cmp    rbx,rsi
	jb     loc_140003735 
	mov    eax,DWORD [rsp+0x2c]
	add    rsi,rax
	cmp    rbx,rsi
	jae    loc_140003735 
	mov    QWORD [rsp+0x30],0x0
	mov    QWORD [rsp+0x38],0x0
	lea    rdi,[rsp+0x40]
	mov    eax,0x0
	mov    ecx,0x1e
	rep stosq
	mov    DWORD [rdi],0x0
	mov    QWORD [rsp+0x140],0x0
	mov    QWORD [rsp+0x148],0x0
	lea    rdi,[rsp+0x150]
	mov    ecx,0x1e
	rep stosq
	mov    DWORD [rdi],0x0
	cmp    BYTE [rbx],0x2e
	je     loc_14000374c 
	mov    eax,0x1
loc_14000369e:
	mov    r8,rax
	add    rax,0x1
	cmp    BYTE [rbx+rax*1-0x1],0x2e
	jne    loc_14000369e 
	mov    esi,r8d
loc_1400036af:
	lea    rcx,[rsp+0x30]
	mov    rdx,rbx
	call   memcpy
	lea    ecx,[rsi+0x1]
	movsxd rcx,ecx
	add    rcx,rbx
	cmp    BYTE [rcx],0x0
	je     loc_14000375c 
	mov    eax,0x1
	movsxd rdx,esi
	add    rdx,rbx
loc_1400036d9:
	mov    r8,rax
	add    rax,0x1
	cmp    BYTE [rdx+rax*1],0x0
	jne    loc_1400036d9 
loc_1400036e6:
	lea    rax,[rsp+0x140]
	mov    rdx,rcx
	mov    rcx,rax
	call   memcpy
	lea    rcx,[rsp+0x30]
	call   load_module
	mov    rbx,rax
	mov    eax,0x0
	test   rbx,rbx
	je     loc_140003738 
	lea    rcx,[rsp+0x140]
	mov    edx,0x0
	call   calc_crc
	mov    edx,eax
	mov    r8d,0xffffffff
	mov    rcx,rbx
	call   get_proc_by_crc
	mov    rbx,rax
loc_140003735:
	mov    rax,rbx
loc_140003738:
	add    rsp,0x258
	pop    rbx
	pop    rsi
	pop    rdi
	pop    rbp
	pop    r12
	pop    r13
	pop    r14
	pop    r15
	ret
loc_14000374c:
	mov    esi,0x0
	mov    r8d,0x0
	jmp    loc_1400036af 
loc_14000375c:
	mov    r8d,0x0
	jmp    loc_1400036e6 
loc_140003764:
	mov    eax,0x0
	jmp    loc_140003738 