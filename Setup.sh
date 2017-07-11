#!/bin/bash


tput setaf 1;echo "//   █████╗ ███╗   ███╗██████╗ ███████╗██████╗ "
tput setaf 1;echo "//  ██╔══██╗████╗ ████║██╔══██╗██╔════╝██╔══██╗"
tput setaf 1;echo "//  ███████║██╔████╔██║██████╔╝█████╗  ██████╔╝"
tput setaf 1;echo "//  ██╔══██║██║╚██╔╝██║██╔══██╗██╔══╝  ██╔══██╗"
tput setaf 1;echo "//  ██║  ██║██║ ╚═╝ ██║██████╔╝███████╗██║  ██║"
tput setaf 1;echo "//  ╚═╝  ╚═╝╚═╝     ╚═╝╚═════╝ ╚══════╝╚═╝  ╚═╝"
tput setaf 1;echo "//  POC Packer For Ophio              "     
tput setaf 1;echo " "
tput setaf 2;echo "Author: Ege Balcı"
tput setaf 4;echo "Source: github.com/EgeBalci/Amber"

echo " "
echo " "
echo " "
tput setaf 3;echo "[*] Installing dependencies..."

sudo apt-get update
sudo apt-get upgrade

sudo apt-get install -y golang nasm wine mingw-w64-i686-dev mingw-w64-tools mingw-w64-x86-64-dev mingw-w64-common mingw-w64 mingw-ocaml gcc-multilib g++-multilib

tput setaf 3;echo "[*] Cloning git tools..."

git clone https://github.com/EgeBalci/MapPE.git

export AMBERPATH=$(pwd)
cd lib
export GOPATH=$(pwd)
cd ..

tput setaf 3;echo "[*] AMBERPATH=$AMBERPATH"
tput setaf 3;echo "[*] GOPATH=$GOPATH"

mv MapPE/MapPE.exe $AMBERPATH

go build -ldflags "-s -w" src/amber.go
go build -ldflags "-s -w" src/handler.go

echo "#!/bin/bash" > /tmp/amber
echo "cd $AMBERPATH" >> /tmp/amber
echo "./amber \$@" >> /tmp/amber
sudo mv /tmp/amber /usr/local/bin/
sudo chmod 777 /usr/local/bin/amber


tput setaf 4;echo "[+] Setup completed !"

