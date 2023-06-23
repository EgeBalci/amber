[BITS 64]


load_module:
	push   rdi
	sub    rsp,0x250
	mov    r8,rcx
	mov    QWORD [rsp+0x248],0x0
	mov    DWORD [rsp+0x234],0x0
	lea    rdi,[rsp+0x20]
	mov    ecx,0x41
	mov    eax,0x0
	rep stosq
	cmp    BYTE [r8],0x0
	je     loc_140003873 
	mov    edx,0x1
loc_1400037eb:
	mov    rax,rdx
	add    rdx,0x1
	cmp    BYTE [r8+rdx*1-0x1],0x0
	jne    loc_1400037eb 
	lea    edx,[rax+rax*1]
	mov    WORD [rsp+0x230],dx
	add    edx,0x2
	mov    WORD [rsp+0x232],dx
	lea    rdx,[rsp+0x20]
	mov    QWORD [rsp+0x238],rdx
	sub    eax,0x1
	js     loc_140003837 
	cdqe
loc_140003824:
	movsx  dx,BYTE [r8+rax*1]
	mov    WORD [rsp+rax*2+0x20],dx
	sub    rax,0x1
	test   eax,eax
	jns    loc_140003824 
loc_140003837:
	lea    r9,[rsp+0x248]
	lea    r8,[rsp+0x230]
	mov    edx,0x0
	mov    ecx,0x0
	mov    r10, 0xB4EBB9A4
	call   api_call
	call   rax        ; <LdrLoadDll> 
	test   eax,eax
	js     loc_14000386c 
	mov    rax,QWORD [rsp+0x248]
loc_140003863:
	add    rsp,0x250
	pop    rdi
	ret
loc_14000386c:
	mov    eax,0x0
	jmp    loc_140003863 
loc_140003873:
	mov    WORD [rsp+0x230],0x0
	mov    WORD [rsp+0x232],0x2
	lea    rax,[rsp+0x20]
	mov    QWORD [rsp+0x238],rax
	jmp    loc_140003837