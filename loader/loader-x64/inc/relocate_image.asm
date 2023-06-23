[BITS 64]


relocate_image:
	mov    r9,rcx
	movsxd rdx,DWORD [rcx+0x3c]
	add    rdx,rcx
	mov    eax,DWORD [rdx+0xb0]
	mov    ecx,0x0
	test   eax,eax
	je     loc_14000261a 
	mov    eax,eax
	lea    rcx,[r9+rax*1]
	mov    r10,r9
	sub    r10,QWORD [rdx+0x30]
	cmp    DWORD [rcx],0x0
	jne    loc_140002601 
	mov    ecx,0x1
	jmp    loc_14000261a 
loc_1400025bd:
	mov    edx,DWORD [rcx]
	movzx  r8d,WORD [rax]
	and    r8d,0xfff
	add    rdx,r8
	add    QWORD [r9+rdx*1],r10
loc_1400025d1:
	add    rax,0x2
	mov    edx,DWORD [rcx+0x4]
	add    rdx,rcx
	cmp    rax,rdx
	je     loc_1400025f9 
loc_1400025e0:
	movzx  edx,BYTE [rax+0x1]
	mov    r8d,edx
	and    r8d,0xfffffff0
	cmp    r8b,0xa0
	je     loc_1400025bd 
	cmp    dl,0xf
	jbe    loc_1400025d1 
	jmp    loc_1400025d1 
loc_1400025f9:
	mov    rcx,rax
loc_1400025fc:
	cmp    DWORD [rcx],0x0
	je     loc_140002615 
loc_140002601:
	lea    rax,[rcx+0x8]
	mov    edx,DWORD [rcx+0x4]
	add    rdx,rcx
	cmp    rax,rdx
	jne    loc_1400025e0 
	mov    rcx,rdx
	jmp    loc_1400025fc 
loc_140002615:
	mov    ecx,0x1
loc_14000261a:
	mov    eax,ecx
	ret