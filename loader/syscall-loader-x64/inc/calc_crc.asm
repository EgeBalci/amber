[BITS 64]


calc_crc:
	test   dx,dx
	je     loc_1400039c9 
	mov    r8,rcx
	movzx  edx,dx
	lea    eax,[rdx-0x1]
	lea    rdx,[rcx+rax*1+0x1]
	mov    eax,0x0
loc_1400039b8:
	crc32  eax,BYTE [r8]
	add    r8,0x1
	cmp    r8,rdx
	jne    loc_1400039b8 
	jmp    loc_1400039ea 
loc_1400039c9:
	movzx  edx,BYTE [rcx]
	test   dl,dl
	je     loc_1400039eb 
	add    rcx,0x1
	mov    eax,0x0
loc_1400039d9:
	crc32  eax,dl
	add    rcx,0x1
	movzx  edx,BYTE [rcx-0x1]
	test   dl,dl
	jne    loc_1400039d9 
loc_1400039ea:
	ret
loc_1400039eb:
	mov    eax,0x0
	jmp    loc_1400039ea