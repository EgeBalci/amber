/*
 * =====================================================================================
 *
 *       Filename:  test.cpp
 *
 *    Description:  
 *
 *        Version:  1.0
 *        Created:  22-06-2018 15:30:52
 *       Revision:  none
 *       Compiler:  gcc
 *
 *         Author:  YOUR NAME (), 
 *   Organization:  
 *
 * =====================================================================================
 */
#include <windows.h>

unsigned char buf[] = { 0x90,0x90,0x90,0x90,0xc3 };

int main(){
	void * mem = (void *)VirtualAlloc(NULL,sizeof(buf),MEM_COMMIT,0x40);
	CreateThread(NULL,0,LPTHREAD_START_ROUTINE(buf),NULL,0,NULL);

	

	HANDLE dll = LoadLibraryA("user32");
	void * proc = (void*)GetProcAddress(HMODULE(dll),"MessageBoxA");



}

