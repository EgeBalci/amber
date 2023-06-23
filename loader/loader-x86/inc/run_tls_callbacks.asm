[BITS 32]


run_tls_callbacks:
loc_406cdf:	push   esi
loc_406ce0:	push   ebx
loc_406ce1:	sub    esp,0x14
loc_406ce4:	mov    esi,DWORD [esp+0x20]
loc_406ce8:	mov    eax,DWORD [esi+0x3c]
loc_406ceb:	mov    edx,DWORD [esi+eax*1+0xc0]
loc_406cf2:	mov    eax,0x0
loc_406cf7:	test   edx,edx
loc_406cf9:	je     loc_406d08 
loc_406cfb:	mov    ebx,DWORD [esi+edx*1+0xc]
loc_406cff:	mov    eax,0x1
loc_406d04:	test   ebx,ebx
loc_406d06:	jne    loc_406d29 
loc_406d08:	add    esp,0x14
loc_406d0b:	pop    ebx
loc_406d0c:	pop    esi
loc_406d0d:	ret
loc_406d0e:	mov    DWORD [esp+0x8],0x0
loc_406d16:	mov    DWORD [esp+0x4],0x1
loc_406d1e:	mov    DWORD [esp],esi
loc_406d21:	call   eax
loc_406d23:	sub    esp,0xc
loc_406d26:	add    ebx,0x4
loc_406d29:	mov    eax,DWORD [ebx]
loc_406d2b:	test   eax,eax
loc_406d2d:	jne    loc_406d0e 
loc_406d2f:	mov    eax,0x1
loc_406d34:	jmp    loc_406d08 
