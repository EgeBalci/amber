[BITS 32]


get_module_by_crc:
loc_4077e7:	push   edi
loc_4077e8:	push   esi
loc_4077e9:	push   ebx
loc_4077ea:	sub    esp,0x8
loc_4077ed:	mov    edi,DWORD [esp+0x18]
loc_4077f1:	mov    eax,fs:0x30
loc_4077f7:	mov    eax,DWORD [eax+0xc]
loc_4077fa:	lea    esi,[eax+0x14]
loc_4077fd:	mov    ebx,DWORD [eax+0x14]
loc_407800:	cmp    esi,ebx
loc_407802:	je     loc_407832 
loc_407804:	movzx  eax,WORD [ebx+0x24]
loc_407808:	mov    DWORD [esp+0x4],eax
loc_40780c:	mov    eax,DWORD [ebx+0x28]
loc_40780f:	mov    DWORD [esp],eax
loc_407812:	call   calc_crc
loc_407817:	cmp    eax,edi
loc_407819:	je     loc_407828 
loc_40781b:	mov    ebx,DWORD [ebx]
loc_40781d:	cmp    esi,ebx
loc_40781f:	jne    loc_407804 
loc_407821:	mov    eax,0x0
loc_407826:	jmp    loc_40782b 
loc_407828:	mov    eax,DWORD [ebx+0x10]
loc_40782b:	add    esp,0x8
loc_40782e:	pop    ebx
loc_40782f:	pop    esi
loc_407830:	pop    edi
loc_407831:	ret
loc_407832:	mov    eax,0x0
loc_407837:	jmp    loc_40782b 
