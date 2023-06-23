[BITS 32]


map_image:
loc_401f00:	push   ebp
loc_401f01:	push   edi
loc_401f02:	push   esi
loc_401f03:	push   ebx
loc_401f04:	sub    esp,0x3c
loc_401f07:	mov    ebp,DWORD [esp+0x50]
loc_401f0b:	mov    edi,ebp
loc_401f0d:	add    edi,DWORD [ebp+0x3c]
loc_401f10:	mov    eax,0x0
loc_401f15:	cmp    DWORD [edi],0x4550
loc_401f1b:	jne    loc_401fcf 
loc_401f21:	mov    DWORD [esp+0x2c],0x0
loc_401f29:	mov    eax,DWORD [edi+0x50]
loc_401f2c:	mov    DWORD [esp+0x28],eax
loc_401f30:	mov    DWORD [esp+0x14],0x4
loc_401f38:	mov    DWORD [esp+0x10],0x103000
loc_401f40:	lea    eax,[esp+0x28]
loc_401f44:	mov    DWORD [esp+0xc],eax
loc_401f48:	mov    DWORD [esp+0x8],0x0
loc_401f50:	lea    eax,[esp+0x2c]
loc_401f54:	mov    DWORD [esp+0x4],eax
loc_401f58:	mov    DWORD [esp],0xffffffff
loc_333333: push   0x99CE7C55
loc_444444: call   api_call
loc_yyyyyy: add    esp,4
loc_401f5f:	call   eax
loc_401f65:	sub    esp,0x18
loc_401f68:	mov    edx,eax
loc_401f6a:	mov    eax,0x0
loc_401f6f:	test   edx,edx
loc_401f71:	js     loc_401fcf 
loc_401f73:	mov    eax,DWORD [edi+0x54]
loc_401f76:	mov    DWORD [esp+0x8],eax
loc_401f7a:	mov    DWORD [esp+0x4],ebp
loc_401f7e:	mov    eax,DWORD [esp+0x2c]
loc_401f82:	mov    DWORD [esp],eax
loc_401f85:	call   memcpy
loc_401f8a:	movzx  eax,WORD [edi+0x14]
loc_401f8e:	lea    ebx,[edi+eax*1+0x18]
loc_401f92:	cmp    WORD [edi+0x6],0x0
loc_401f97:	je     loc_401fcb 
loc_401f99:	mov    esi,0x0
loc_401f9e:	mov    eax,DWORD [ebx+0xc]
loc_401fa1:	add    eax,DWORD [esp+0x2c]
loc_401fa5:	mov    edx,ebp
loc_401fa7:	add    edx,DWORD [ebx+0x14]
loc_401faa:	mov    ecx,DWORD [ebx+0x10]
loc_401fad:	mov    DWORD [esp+0x8],ecx
loc_401fb1:	mov    DWORD [esp+0x4],edx
loc_401fb5:	mov    DWORD [esp],eax
loc_401fb8:	call   memcpy
loc_401fbd:	add    esi,0x1
loc_401fc0:	add    ebx,0x28
loc_401fc3:	movzx  eax,WORD [edi+0x6]
loc_401fc7:	cmp    eax,esi
loc_401fc9:	jg     loc_401f9e 
loc_401fcb:	mov    eax,DWORD [esp+0x2c]
loc_401fcf:	add    esp,0x3c
loc_401fd2:	pop    ebx
loc_401fd3:	pop    esi
loc_401fd4:	pop    edi
loc_401fd5:	pop    ebp
loc_401fd6:	ret
