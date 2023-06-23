[BITS 32]


resolve_imports:
loc_4042ff:	push   ebp
loc_404300:	push   edi
loc_404301:	push   esi
loc_404302:	push   ebx
loc_404303:	sub    esp,0x2c
loc_404306:	mov    ebp,DWORD [esp+0x40]
loc_40430a:	mov    eax,DWORD [ebp+0x3c]
loc_40430d:	mov    edx,DWORD [ebp+eax*1+0x80]
loc_404314:	mov    eax,0x0
loc_404319:	test   edx,edx
loc_40431b:	je     loc_4043d5 
loc_404321:	lea    eax,[ebp+edx*1+0x0]
loc_404325:	mov    DWORD [esp+0x1c],eax
loc_404329:	mov    eax,DWORD [eax+0xc]
loc_40432c:	test   eax,eax
loc_40432e:	jne    loc_4043a4 
loc_404330:	mov    eax,0x1
loc_404335:	jmp    loc_4043d5 
loc_40433a:	mov    DWORD [esp+0x8],eax
loc_40433e:	mov    DWORD [esp+0x4],0x0
loc_404346:	mov    DWORD [esp],edi
loc_404349:	call   get_proc_by_crc
loc_40434e:	test   eax,eax
loc_404350:	je     loc_404354 
loc_404352:	mov    DWORD [esi],eax
loc_404354:	add    ebx,0x4
loc_404357:	add    esi,0x4
loc_40435a:	mov    eax,DWORD [ebx]
loc_40435c:	test   eax,eax
loc_40435e:	je     loc_404394 
loc_404360:	test   eax,eax
loc_404362:	js     loc_40433a 
loc_404364:	mov    DWORD [esp+0x4],0x0
loc_40436c:	lea    eax,[ebp+eax*1+0x2]
loc_404370:	mov    DWORD [esp],eax
loc_404373:	call   calc_crc
loc_404378:	mov    DWORD [esp+0x8],0xffffffff
loc_404380:	mov    DWORD [esp+0x4],eax
loc_404384:	mov    DWORD [esp],edi
loc_404387:	call   get_proc_by_crc
loc_40438c:	test   eax,eax
loc_40438e:	je     loc_404354 
loc_404390:	mov    DWORD [esi],eax
loc_404392:	jmp    loc_404354 
loc_404394:	add    DWORD [esp+0x1c],0x14
loc_404399:	mov    eax,DWORD [esp+0x1c]
loc_40439d:	mov    eax,DWORD [eax+0xc]
loc_4043a0:	test   eax,eax
loc_4043a2:	je     loc_4043c9 
loc_4043a4:	add    eax,ebp
loc_4043a6:	mov    DWORD [esp],eax
loc_4043a9:	call   load_module
loc_4043ae:	mov    edi,eax
loc_4043b0:	test   eax,eax
loc_4043b2:	je     loc_4043d0 
loc_4043b4:	mov    eax,DWORD [esp+0x1c]
loc_4043b8:	mov    ebx,ebp
loc_4043ba:	add    ebx,DWORD [eax]
loc_4043bc:	mov    esi,ebp
loc_4043be:	add    esi,DWORD [eax+0x10]
loc_4043c1:	mov    eax,DWORD [ebx]
loc_4043c3:	test   eax,eax
loc_4043c5:	jne    loc_404360 
loc_4043c7:	jmp    loc_404394 
loc_4043c9:	mov    eax,0x1
loc_4043ce:	jmp    loc_4043d5 
loc_4043d0:	mov    eax,0x0
loc_4043d5:	add    esp,0x2c
loc_4043d8:	pop    ebx
loc_4043d9:	pop    esi
loc_4043da:	pop    edi
loc_4043db:	pop    ebp
loc_4043dc:	ret
