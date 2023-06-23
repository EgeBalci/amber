[BITS 32]


get_proc_by_crc:
loc_408e97:	push   ebp
loc_408e98:	mov    ebp,esp
loc_408e9a:	push   edi
loc_408e9b:	push   esi
loc_408e9c:	push   ebx
loc_408e9d:	and    esp,0xfffffff0
loc_408ea0:	sub    esp,0x240
loc_408ea6:	mov    ebx,DWORD [ebp+0x8]
loc_408ea9:	mov    eax,ebx
loc_408eab:	add    eax,DWORD [ebx+0x3c]
loc_408eae:	mov    edx,ebx
loc_408eb0:	add    edx,DWORD [eax+0x78]
loc_408eb3:	mov    eax,DWORD [eax+0x7c]
loc_408eb6:	mov    DWORD [esp+0x1c],eax
loc_408eba:	mov    edi,DWORD [edx+0x20]
loc_408ebd:	mov    eax,DWORD [edx+0x1c]
loc_408ec0:	mov    DWORD [esp+0x2c],eax
loc_408ec4:	mov    eax,DWORD [edx+0x24]
loc_408ec7:	mov    DWORD [esp+0x28],eax
loc_408ecb:	mov    ecx,DWORD [edx+0x18]
loc_408ece:	test   ecx,ecx
loc_408ed0:	je     loc_409041 
loc_408ed6:	mov    esi,0x0
loc_408edb:	add    edi,ebx
loc_408edd:	mov    DWORD [esp+0x24],edx
loc_408ee1:	mov    DWORD [esp+0x20],ecx
loc_408ee5:	mov    DWORD [esp+0x4],0x0
loc_408eed:	mov    eax,ebx
loc_408eef:	add    eax,DWORD [edi+esi*4]
loc_408ef2:	mov    DWORD [esp],eax
loc_408ef5:	call   calc_crc
loc_408efa:	cmp    DWORD [ebp+0x10],esi
loc_408efd:	je     loc_408f19 
loc_408eff:	cmp    eax,DWORD [ebp+0xc]
loc_408f02:	je     loc_408f19 
loc_408f04:	add    esi,0x1
loc_408f07:	mov    eax,DWORD [esp+0x20]
loc_408f0b:	cmp    esi,eax
loc_408f0d:	jne    loc_408ee5 
loc_408f0f:	mov    eax,0x0
loc_408f14:	jmp    loc_409023 
loc_408f19:	mov    edx,DWORD [esp+0x24]
loc_408f1d:	lea    eax,[ebx+esi*2]
loc_408f20:	mov    ecx,DWORD [esp+0x28]
loc_408f24:	movzx  eax,WORD [eax+ecx*1]
loc_408f28:	lea    eax,[ebx+eax*4]
loc_408f2b:	mov    ecx,DWORD [esp+0x2c]
loc_408f2f:	add    ebx,DWORD [eax+ecx*1]
loc_408f32:	cmp    ebx,edx
loc_408f34:	jb     loc_409021 
loc_408f3a:	mov    eax,DWORD [esp+0x1c]
loc_408f3e:	add    edx,eax
loc_408f40:	cmp    ebx,edx
loc_408f42:	jae    loc_409021 
loc_408f48:	vpxor  xmm0,xmm0,xmm0
loc_408f4c:	vmovdqu [esp+0x38],xmm0
loc_408f52:	lea    edi,[esp+0x48]
loc_408f56:	mov    eax,0x0
loc_408f5b:	mov    ecx,0x3d
loc_408f60:	rep stosd
loc_408f62:	vmovdqu [esp+0x13c],xmm0
loc_408f6b:	lea    edi,[esp+0x14c]
loc_408f72:	mov    ecx,0x3d
loc_408f77:	rep stosd
loc_408f79:	cmp    BYTE [ebx],0x2e
loc_408f7c:	je     loc_40902b 
loc_408f82:	mov    esi,0x0
loc_408f87:	add    esi,0x1
loc_408f8a:	mov    eax,esi
loc_408f8c:	cmp    BYTE [ebx+esi*1],0x2e
loc_408f90:	jne    loc_408f87 
loc_408f92:	lea    edx,[esp+0x38]
loc_408f96:	mov    DWORD [esp+0x8],eax
loc_408f9a:	mov    DWORD [esp+0x4],ebx
loc_408f9e:	mov    DWORD [esp],edx
loc_408fa1:	call   memcpy
loc_408fa6:	lea    ecx,[ebx+esi*1+0x1]
loc_408faa:	cmp    BYTE [ecx],0x0
loc_408fad:	je     loc_40903a 
loc_408fb3:	mov    eax,0x0
loc_408fb8:	add    esi,ebx
loc_408fba:	add    eax,0x1
loc_408fbd:	mov    edx,eax
loc_408fbf:	cmp    BYTE [esi+eax*1+0x1],0x0
loc_408fc4:	jne    loc_408fba 
loc_408fc6:	lea    eax,[esp+0x13c]
loc_408fcd:	mov    DWORD [esp+0x8],edx
loc_408fd1:	mov    DWORD [esp+0x4],ecx
loc_408fd5:	mov    DWORD [esp],eax
loc_408fd8:	call   memcpy
loc_408fdd:	lea    eax,[esp+0x38]
loc_408fe1:	mov    DWORD [esp],eax
loc_408fe4:	call   load_module
loc_408fe9:	mov    ebx,eax
loc_408feb:	mov    eax,0x0
loc_408ff0:	test   ebx,ebx
loc_408ff2:	je     loc_409023 
loc_408ff4:	mov    DWORD [esp+0x4],0x0
loc_408ffc:	lea    eax,[esp+0x13c]
loc_409003:	mov    DWORD [esp],eax
loc_409006:	call   calc_crc
loc_40900b:	mov    DWORD [esp+0x8],0xffffffff
loc_409013:	mov    DWORD [esp+0x4],eax
loc_409017:	mov    DWORD [esp],ebx
loc_40901a:	call   get_proc_by_crc
loc_40901f:	mov    ebx,eax
loc_409021:	mov    eax,ebx
loc_409023:	lea    esp,[ebp-0xc]
loc_409026:	pop    ebx
loc_409027:	pop    esi
loc_409028:	pop    edi
loc_409029:	pop    ebp
loc_40902a:	ret
loc_40902b:	mov    esi,0x0
loc_409030:	mov    eax,0x0
loc_409035:	jmp    loc_408f92 
loc_40903a:	mov    edx,0x0
loc_40903f:	jmp    loc_408fc6 
loc_409041:	mov    eax,0x0
loc_409046:	jmp    loc_409023 
