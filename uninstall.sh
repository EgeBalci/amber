#!/bin/bash

sudo rm -r /usr/local/bin/amber
sudo rm -r BitBender
sudo rm -r bitbender
sudo rm -r MapPE
sudo rm -r MapPE.exe
sudo rm -r src
sudo rm -r pkg
sudo rm -r amber
sudo rm -r handler

echo 0 > Stub/PAYLOAD.h
echo 0 > Stub/KEY.h

