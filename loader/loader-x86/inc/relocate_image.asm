[BITS 32]


relocate_image:
loc_403055:	push   ebp
loc_403056:	push   edi
loc_403057:	push   esi
loc_403058:	push   ebx
loc_403059:	mov    ebp,DWORD [esp+0x14]
loc_40305d:	mov    ebx,ebp
loc_40305f:	add    ebx,DWORD [ebp+0x3c]
loc_403062:	mov    edx,DWORD [ebx+0xa0]
loc_403068:	mov    eax,0x0
loc_40306d:	test   edx,edx
loc_40306f:	je     loc_4030d3 
loc_403071:	add    edx,ebp
loc_403073:	mov    esi,ebp
loc_403075:	sub    esi,DWORD [ebx+0x34]
loc_403078:	cmp    DWORD [edx],0x0
loc_40307b:	jne    loc_4030be 
loc_40307d:	mov    eax,0x1
loc_403082:	jmp    loc_4030d3 
loc_403084:	movzx  ecx,WORD [eax]
loc_403087:	and    ecx,0xfff
loc_40308d:	add    ecx,DWORD [edx]
loc_40308f:	add    DWORD [ebp+ecx*1+0x0],esi
loc_403093:	add    eax,0x2
loc_403096:	mov    ecx,edx
loc_403098:	add    ecx,DWORD [edx+0x4]
loc_40309b:	cmp    eax,ecx
loc_40309d:	je     loc_4030b7 
loc_40309f:	movzx  ecx,BYTE [eax+0x1]
loc_4030a3:	mov    edi,ecx
loc_4030a5:	and    edi,0xfffffff0
loc_4030a8:	mov    ebx,edi
loc_4030aa:	cmp    bl,0x30
loc_4030ad:	je     loc_403084 
loc_4030af:	cmp    cl,0xf
loc_4030b2:	jbe    loc_403093 
loc_4030b4:	int3
loc_4030b5:	jmp    loc_403093 
loc_4030b7:	mov    edx,eax
loc_4030b9:	cmp    DWORD [edx],0x0
loc_4030bc:	je     loc_4030ce 
loc_4030be:	lea    eax,[edx+0x8]
loc_4030c1:	mov    ecx,edx
loc_4030c3:	add    ecx,DWORD [edx+0x4]
loc_4030c6:	cmp    eax,ecx
loc_4030c8:	jne    loc_40309f 
loc_4030ca:	mov    edx,ecx
loc_4030cc:	jmp    loc_4030b9 
loc_4030ce:	mov    eax,0x1
loc_4030d3:	pop    ebx
loc_4030d4:	pop    esi
loc_4030d5:	pop    edi
loc_4030d6:	pop    ebp
loc_4030d7:	ret
