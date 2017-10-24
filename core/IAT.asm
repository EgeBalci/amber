;
; This file contains the import address table index addresses.
;
; Author: Ege Balcı	<ege.balci@invictuseurope.com>
;


[BITS 32]
[ORG 0]

%define LLA 0x00000000			; &(LoadLibraryA())
%define GPA 0x00000000			; &(GetProcAddress())
%define VA 0x00000000			; &(VirtualAlloc())
%define VP 0x00000000			; &(VirtualProtect())
%define CT 0x00000000			; &(CreateThread())
