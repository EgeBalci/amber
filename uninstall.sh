#!/bin/bash

sudo rm -r /usr/local/bin/amber
sudo rm core/NonASLR/Payload
sudo rm core/NonASLR/iat/Payload
sudo rm core/NonASLR/Mem.map
sudo rm core/NonASLR/iat/Mem.map
sudo rm core/ASLR/Payload
sudo rm core/ASLR/iat/Payload
sudo rm core/ASLR/Mem.map
sudo rm core/ASLR/iat/Mem.map
sudo rm -r amber
sudo rm -r handler

echo 0 > Stub/PAYLOAD.h
echo 0 > Stub/KEY.h

echo " "
echo "[+] Uninstallation complete."


