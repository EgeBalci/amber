[BITS 32]

protect_sections:
loc_4057e8:	push   ebp
loc_4057e9:	push   edi
loc_4057ea:	push   esi
loc_4057eb:	push   ebx
loc_4057ec:	sub    esp,0x4c
loc_4057ef:	mov    ebp,DWORD [esp+0x60]
loc_4057f3:	mov    edi,ebp
loc_4057f5:	add    edi,DWORD [ebp+0x3c]
loc_4057f8:	movzx  eax,WORD [edi+0x14]
loc_4057fc:	lea    ebx,[edi+eax*1+0x18]
loc_405800:	mov    DWORD [esp+0x3c],0x0
loc_405808:	cmp    WORD [edi+0x6],0x0
loc_40580d:	je     loc_405840 
loc_40580f:	mov    esi,0x0
loc_405814:	mov    ebp,edi
loc_405816:	jmp    loc_40591e 
loc_40581b:	mov    BYTE [esp+0x2f],0x1
loc_405820:	jmp    loc_405863 
loc_405822:	mov    edx,0x20
loc_405827:	jmp    loc_4058ba 
loc_40582c:	mov    eax,0x0
loc_405831:	jmp    loc_405838 
loc_405833:	mov    eax,0x1
loc_405838:	add    esp,0x4c
loc_40583b:	pop    ebx
loc_40583c:	pop    esi
loc_40583d:	pop    edi
loc_40583e:	pop    ebp
loc_40583f:	ret
loc_405840:	mov    eax,0x1
loc_405845:	jmp    loc_405838 
loc_405847:	mov    edi,eax
loc_405849:	shr    edi,0x1f
loc_40584c:	mov    BYTE [esp+0x2f],0x0
loc_405851:	test   eax,0x20000000
loc_405856:	je     loc_405892 
loc_405858:	mov    edx,0x10
loc_40585d:	mov    ecx,edi
loc_40585f:	test   cl,cl
loc_405861:	je     loc_405897 
loc_405863:	and    eax,0x60000000
loc_405868:	mov    edi,0x1
loc_40586d:	cmp    eax,0x60000000
loc_405872:	mov    edx,0x20
loc_405877:	mov    eax,0x80
loc_40587c:	cmovne edx,eax
loc_40587f:	jmp    loc_4058a7 
loc_405881:	test   eax,0x20000000
loc_405886:	jne    loc_40581b 
loc_405888:	mov    BYTE [esp+0x2f],0x1
loc_40588d:	mov    edx,0x4
loc_405892:	mov    edi,0x0
loc_405897:	and    eax,0x60000000
loc_40589c:	cmp    eax,0x60000000
loc_4058a1:	je     loc_405822 
loc_4058a7:	cmp    BYTE [esp+0x2f],0x0
loc_4058ac:	je     loc_4058ba 
loc_4058ae:	mov    eax,edi
loc_4058b0:	test   al,al
loc_4058b2:	mov    eax,0x40
loc_4058b7:	cmovne edx,eax
loc_4058ba:	mov    eax,DWORD [esp+0x60]
loc_4058be:	mov    ecx,DWORD [esp+0x28]
loc_4058c2:	add    eax,DWORD [ecx+0xc]
loc_4058c5:	mov    DWORD [esp+0x3c],eax
loc_4058c9:	mov    eax,DWORD [ecx+0x10]
loc_4058cc:	mov    DWORD [esp+0x34],eax
loc_4058d0:	mov    DWORD [esp+0x38],0x0
loc_4058d8:	lea    eax,[esp+0x38]
loc_4058dc:	mov    DWORD [esp+0x10],eax
loc_4058e0:	mov    DWORD [esp+0xc],edx
loc_4058e4:	lea    eax,[esp+0x34]
loc_4058e8:	mov    DWORD [esp+0x8],eax
loc_4058ec:	lea    eax,[esp+0x3c]
loc_4058f0:	mov    DWORD [esp+0x4],eax
loc_4058f4:	mov    DWORD [esp],0xffffffff
loc_555555: push   0x6EDE4D41
loc_666666: call   api_call
loc_zzzzzz: add    esp,4
loc_4058fb:	call   eax
loc_405901:	sub    esp,0x14
loc_405904:	test   eax,eax
loc_405906:	js     loc_40582c 
loc_40590c:	add    esi,0x1
loc_40590f:	add    ebx,0x28
loc_405912:	movzx  eax,WORD [ebp+0x6]
loc_405916:	cmp    eax,esi
loc_405918:	jle    loc_405833 
loc_40591e:	mov    DWORD [esp+0x28],ebx
loc_405922:	mov    eax,DWORD [ebx+0x24]
loc_405925:	test   eax,eax
loc_405927:	je     loc_40590c 
loc_405929:	cdq
loc_40592a:	and    edx,0xffffffc8
loc_40592d:	add    edx,0x40
loc_405930:	test   eax,0x40000000
loc_405935:	je     loc_405847 
loc_40593b:	test   eax,eax
loc_40593d:	js     loc_405881 
loc_405943:	mov    edi,0x0
loc_405948:	mov    BYTE [esp+0x2f],0x1
loc_40594d:	mov    edx,0x2
loc_405952:	jmp    loc_405851 
