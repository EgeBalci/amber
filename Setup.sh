#!/bin/bash


echo "//   █████╗ ███╗   ███╗██████╗ ███████╗██████╗ "
echo "//  ██╔══██╗████╗ ████║██╔══██╗██╔════╝██╔══██╗"
echo "//  ███████║██╔████╔██║██████╔╝█████╗  ██████╔╝"
echo "//  ██╔══██║██║╚██╔╝██║██╔══██╗██╔══╝  ██╔══██╗"
echo "//  ██║  ██║██║ ╚═╝ ██║██████╔╝███████╗██║  ██║"
echo "//  ╚═╝  ╚═╝╚═╝     ╚═╝╚═════╝ ╚══════╝╚═╝  ╚═╝"
echo "//  POC Crypter For ReplaceProcess              "     
echo " "
echo "Author: Ege Balcı"
echo "Source: github.com/EgeBalci/Amber"

echo " "
echo " "
echo " "
echo "[*] Installing dependencies..."

sudo apt-get update
sudo apt-get upgrade

sudo apt-get install -y golang nasm wine mingw-w64-i686-dev mingw-w64-tools mingw-w64-x86-64-dev mingw-w64-common mingw-w64 mingw-ocaml gcc-multilib g++-multilib

echo "[*] Cloning git tools..."

git clone https://github.com/EgeBalci/MapPE.git

export AMBERPATH=$(pwd)
cd lib
export GOPATH=$(pwd)
cd ..

echo "[*] AMBERPATH=$AMBERPATH"
echo "[*] GOPATH=$GOPATH"

mv MapPE/MapPE.exe $AMBERPATH

go build -ldflags "-s -w" amber.go
go build -ldflags "-s -w" handler.go

echo "#!/bin/bash" >> /tmp/amber
echo "cd $AMBERPATH" >> /tmp/amber
echo "./amber \$@" >> /tmp/amber
sudo mv /tmp/amber /usr/local/bin/
sudo chmod 777 /usr/local/bin/amber


echo "[+] Setup completed !"

