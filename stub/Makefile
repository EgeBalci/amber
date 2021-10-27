normal:
	i686-w64-mingw32-windres Resource.rc -o res.o
	i686-w64-mingw32-g++-win32 -c stub.c -o stub.o
	i686-w64-mingw32-g++-win32 -Wl,--subsystem,windows stub.o res.o -o stub32.exe
	i686-w64-mingw32-strip stub32.exe
	x86_64-w64-mingw32-windres Resource.rc -o res.o
	x86_64-w64-mingw32-g++-win32 -c stub.c -o stub.o
	x86_64-w64-mingw32-g++-win32 -Wl,--subsystem,windows stub.o res.o -o stub.exe
	x86_64-w64-mingw32-strip stub.exe
	rm *.o
