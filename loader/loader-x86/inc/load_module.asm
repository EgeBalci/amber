[BITS 32]


load_module:
loc_40a08a:	push   edi
loc_40a08b:	push   ebx
loc_40a08c:	sub    esp,0x234
loc_40a092:	mov    ebx,DWORD [esp+0x240]
loc_40a099:	mov    DWORD [esp+0x22c],0x0
loc_40a0a4:	lea    edi,[esp+0x1c]
loc_40a0a8:	mov    ecx,0x82
loc_40a0ad:	mov    eax,0x0
loc_40a0b2:	rep stosd
loc_40a0b4:	cmp    BYTE [ebx],0x0
loc_40a0b7:	je     loc_40a14c 
loc_40a0bd:	mov    edx,0x0
loc_40a0c2:	mov    eax,edx
loc_40a0c4:	add    edx,0x1
loc_40a0c7:	cmp    BYTE [ebx+edx*1],0x0
loc_40a0cb:	jne    loc_40a0c2 
loc_40a0cd:	add    edx,edx
loc_40a0cf:	mov    WORD [esp+0x224],dx
loc_40a0d7:	add    edx,0x2
loc_40a0da:	mov    WORD [esp+0x226],dx
loc_40a0e2:	lea    edx,[esp+0x1c]
loc_40a0e6:	mov    DWORD [esp+0x228],edx
loc_40a0ed:	test   eax,eax
loc_40a0ef:	js     loc_40a103 
loc_40a0f1:	movsx  dx,BYTE [ebx+eax*1]
loc_40a0f6:	mov    WORD [esp+eax*2+0x1c],dx
loc_40a0fb:	sub    eax,0x1
loc_40a0fe:	cmp    eax,0xffffffff
loc_40a101:	jne    loc_40a0f1 
loc_40a103:	lea    eax,[esp+0x22c]
loc_40a10a:	mov    DWORD [esp+0xc],eax
loc_40a10e:	lea    eax,[esp+0x224]
loc_40a115:	mov    DWORD [esp+0x8],eax
loc_40a119:	mov    DWORD [esp+0x4],0x0
loc_40a121:	mov    DWORD [esp],0x0
loc_111111: push   0xB4EBB9A4
loc_222222: call   api_call
loc_xxxxxx: add    esp,4
loc_40a128:	call   eax
loc_40a12e:	sub    esp,0x10
loc_40a131:	test   eax,eax
loc_40a133:	js     loc_40a145 
loc_40a135:	mov    eax,DWORD [esp+0x22c]
loc_40a13c:	add    esp,0x234
loc_40a142:	pop    ebx
loc_40a143:	pop    edi
loc_40a144:	ret
loc_40a145:	mov    eax,0x0
loc_40a14a:	jmp    loc_40a13c 
loc_40a14c:	mov    WORD [esp+0x224],0x0
loc_40a156:	mov    WORD [esp+0x226],0x2
loc_40a160:	lea    eax,[esp+0x1c]
loc_40a164:	mov    DWORD [esp+0x228],eax
loc_40a16b:	jmp    loc_40a103 
