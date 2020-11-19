#include <windows.h>

int main(int argc, char const *argv[])
{
    PIMAGE_DOS_HEADER dosHeader = {};
    PIMAGE_SECTION_HEADER sectionHeader = {};
    PIMAGE_NT_HEADERS imageNTHeaders = {};

    LPVOID moduleHandle = GetModuleHandle(NULL);
    if (moduleHandle == NULL)
        return 1;

    dosHeader = (PIMAGE_DOS_HEADER)moduleHandle;
    imageNTHeaders = (PIMAGE_NT_HEADERS)((DWORD)moduleHandle + dosHeader->e_lfanew);
    DWORD sectionLocation = (DWORD) & (imageNTHeaders->OptionalHeader) + (WORD)imageNTHeaders->FileHeader.SizeOfOptionalHeader;
    DWORD sectionSize = (DWORD)sizeof(IMAGE_SECTION_HEADER);

    for (int i = 0; i < imageNTHeaders->FileHeader.NumberOfSections; i++)
    {
        sectionHeader = (PIMAGE_SECTION_HEADER)sectionLocation;
        // printf("[#] %s\n", sectionHeader->Name);
        // printf("-> 0x%x Virtual Size\n", sectionHeader->Misc.VirtualSize);
        // printf("-> 0x%x Virtual Address\n", sectionHeader->VirtualAddress);
        // printf("-> 0x%x Size Of Raw Data\n", sectionHeader->SizeOfRawData);
        // printf("-> 0x%x Pointer To Raw Data\n", sectionHeader->PointerToRawData);
        // printf("-> 0x%x Pointer To Relocations\n", sectionHeader->PointerToRelocations);
        // printf("-> 0x%x Pointer To Line Numbers\n", sectionHeader->PointerToLinenumbers);
        // printf("-> 0x%x Number Of Relocations\n", sectionHeader->NumberOfRelocations);
        // printf("-> 0x%x Number Of Line Numbers\n", sectionHeader->NumberOfLinenumbers);
        // printf("-> 0x%x Characteristics\n", sectionHeader->Characteristics);
        sectionLocation += sectionSize;
    }

    unsigned char *buffer = (unsigned char *)VirtualAlloc(NULL, sectionHeader->SizeOfRawData, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    memcpy((void *)buffer, (void *)(sectionHeader->VirtualAddress + imageNTHeaders->OptionalHeader.ImageBase), sectionHeader->SizeOfRawData);
    (*(void (*)())buffer)();

    return 0;
}
