[BITS 32]


calc_crc:
loc_40b22f:	mov    eax,DWORD [esp+0x4]
loc_40b233:	mov    edx,DWORD [esp+0x8]
loc_40b237:	test   dx,dx
loc_40b23a:	je     loc_40b256 
loc_40b23c:	mov    ecx,eax
loc_40b23e:	movzx  edx,dx
loc_40b241:	add    eax,edx
loc_40b243:	mov    edx,0x0
loc_40b248:	crc32  edx,BYTE [ecx]
loc_40b24d:	add    ecx,0x1
loc_40b250:	cmp    ecx,eax
loc_40b252:	jne    loc_40b248 
loc_40b254:	jmp    loc_40b275 
loc_40b256:	movzx  ecx,BYTE [eax]
loc_40b259:	test   cl,cl
loc_40b25b:	je     loc_40b278 
loc_40b25d:	add    eax,0x1
loc_40b260:	mov    edx,0x0
loc_40b265:	crc32  edx,cl
loc_40b26a:	add    eax,0x1
loc_40b26d:	movzx  ecx,BYTE [eax-0x1]
loc_40b271:	test   cl,cl
loc_40b273:	jne    loc_40b265 
loc_40b275:	mov    eax,edx
loc_40b277:	ret
loc_40b278:	mov    edx,0x0
loc_40b27d:	jmp    loc_40b275 
