#!/bin/bash


echo "//   █████╗ ███╗   ███╗██████╗ ███████╗██████╗ "
echo "//  ██╔══██╗████╗ ████║██╔══██╗██╔════╝██╔══██╗"
echo "//  ███████║██╔████╔██║██████╔╝█████╗  ██████╔╝"
echo "//  ██╔══██║██║╚██╔╝██║██╔══██╗██╔══╝  ██╔══██╗"
echo "//  ██║  ██║██║ ╚═╝ ██║██████╔╝███████╗██║  ██║"
echo "//  ╚═╝  ╚═╝╚═╝     ╚═╝╚═════╝ ╚══════╝╚═╝  ╚═╝"
echo "//  POC Crypter For ReplaceProcess              "     

echo " "
echo " "
echo " "
echo "[*] Installing dependencies..."

sudo apt-get update
sudo apt-get upgrade

sudo apt-get install -y golang nasm wine mingw-w64-i686-dev mingw-w64-tools mingw-w64-x86-64-dev mingw-w64-common mingw-w64 mingw-ocaml gcc-multilib g++-multilib

echo "[*] Cloning git tools..."

git clone https://github.com/EgeBalci/MapPE.git
git clone https://github.com/EgeBalci/BitBender.git

export AMBERPATH=$(pwd)
export GOPATH=$AMBERPATH

echo "[*] Downloading golang packages..."
go get github.com/fatih/color
go get gopkg.in/cheggaaa/pb.v1

cd BitBender
export GOPATH=$(pwd)
echo "[*] GOPATH=$GOPATH"
go get github.com/fatih/color
echo "[*] Building BitBender..."
go build -ldflags "-s -w" BitBender.go
cd ..
export GOPATH=$AMBERPATH
echo "[*] GOPATH=$GOPATH"

mv MapPE/MapPE.exe $AMBERPATH
mv BitBender/BitBender $AMBERPATH/bitbender

go build -ldflags "-s -w" amber.go
go build -ldflags "-s -w" handler.go

echo "#!/bin/bash" >> /tmp/amber
echo "cd $AMBERPATH" >> /tmp/amber
echo "./amber \$@" >> /tmp/amber
sudo mv /tmp/amber /usr/local/bin/
sudo chmod 777 /usr/local/bin/amber


echo "[+] Setup completed !"

