#!/bin/bash

sudo rm -r /usr/local/bin/amber
sudo rm -r MapPE
sudo rm -r MapPE.exe
sudo rm ReplaceProcess/peb/Mem.map
sudo rm ReplaceProcess/iat/Mem.map
sudo rm -r amber
sudo rm -r handler
sudo rm Stub.o

echo 0 > Stub/PAYLOAD.h
echo 0 > Stub/KEY.h

echo " "
echo "[+] Uninstallation complete."


