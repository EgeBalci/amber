[BITS 64]


get_module_by_crc:
	push   rdi
	push   rsi
	push   rbx
	sub    rsp,0x20
	mov    esi,ecx
	mov    rax,QWORD gs:0x60
	mov    rax,QWORD [rax+0x18]
	lea    rdi,[rax+0x20]
	mov    rbx,QWORD [rax+0x20]
	cmp    rdi,rbx
	je     loc_140102e89 
loc_140102e5d:
	movzx  edx,WORD [rbx+0x48]
	mov    rcx,QWORD [rbx+0x50]
	call   calc_crc
	cmp    eax,esi
	je     loc_140102e7d 
	mov    rbx,QWORD [rbx]
	cmp    rdi,rbx
	jne    loc_140102e5d 
	mov    eax,0x0
	jmp    loc_140102e81 
loc_140102e7d:
	mov    rax,QWORD [rbx+0x20]
loc_140102e81:
	add    rsp,0x20
	pop    rbx
	pop    rsi
	pop    rdi
	ret
loc_140102e89:
	mov    eax,0x0
	jmp    loc_140102e81 
