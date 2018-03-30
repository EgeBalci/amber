#!/bin/bash

# Reset
Color_Off='\033[0m'       # Text Reset

# Regular Colors
Black='\033[0;30m'        # Black
Red='\033[0;31m'          # Red
Green='\033[0;32m'        # Green
Yellow='\033[0;33m'       # Yellow
Blue='\033[0;34m'         # Blue
Purple='\033[0;35m'       # Purple
Cyan='\033[0;36m'         # Cyan
White='\033[0;37m'        # White

# Bold
BBlack='\033[1;30m'       # Black
BRed='\033[1;31m'         # Red
BGreen='\033[1;32m'       # Green
BYellow='\033[1;33m'      # Yellow
BBlue='\033[1;34m'        # Blue
BPurple='\033[1;35m'      # Purple
BCyan='\033[1;36m'        # Cyan
BWhite='\033[1;37m'       # White

# Underline
UBlack='\033[4;30m'       # Black
URed='\033[4;31m'         # Red
UGreen='\033[4;32m'       # Green
UYellow='\033[4;33m'      # Yellow
UBlue='\033[4;34m'        # Blue
UPurple='\033[4;35m'      # Purple
UCyan='\033[4;36m'        # Cyan
UWhite='\033[4;37m'       # White

# Background
On_Black='\033[40m'       # Black
On_Red='\033[41m'         # Red
On_Green='\033[42m'       # Green
On_Yellow='\033[43m'      # Yellow



### Getting OS Information

if [[ $(uname) == "Darwin" ]]; then
	DIST="MacOS"
	DIST_VER=$(uname -r)
else
	if [ -f /etc/lsb-release ]; then
		. /etc/lsb-release
		DIST=$DISTRIB_ID
		DIST_VER=$DISTRIB_RELEASE
	else
		DIST="Unknown"
		DIST_VER="Unknown"
	fi	
fi


if [ $DIST == "Arch" ] || [ $DIST == "Manjaro" ]; then

	echo -e -n $BRed
	echo '//    _____      _____ _______________________________ '
	echo '//   /  _  \    /     \\______   \_   _____/\______   \'
	echo '//  /  /_\  \  /  \ /  \|    |  _/|    __)_  |       _/'
	echo '// /    |    \/    Y    \    |   \|        \ |    |   \'
	echo '// \____|__  /\____|__  /______  /_______  / |____|_  /'
	echo '//         \/         \/       \/        \/         \/ '
	echo '// POC Reflective PE Packer                       1.3.0'
else 

	#tput setaf 1;
	echo -e -n $BRed
	echo "//   █████╗ ███╗   ███╗██████╗ ███████╗██████╗ "
	echo "//  ██╔══██╗████╗ ████║██╔══██╗██╔════╝██╔══██╗"
	echo "//  ███████║██╔████╔██║██████╔╝█████╗  ██████╔╝"
	echo "//  ██╔══██║██║╚██╔╝██║██╔══██╗██╔══╝  ██╔══██╗"
	echo "//  ██║  ██║██║ ╚═╝ ██║██████╔╝███████╗██║  ██║"
	echo "//  ╚═╝  ╚═╝╚═╝     ╚═╝╚═════╝ ╚══════╝╚═╝  ╚═╝"
	echo "//  POC Reflective PE Packer  ☣           1.3.0"     
	#echo -e " "$Color_Off
fi


echo -e  $Blue
echo "Author: Ege Balcı"
echo -e -n $Green
echo "Source: github.com/egebalci/Amber"


echo -e $BYellow
echo "[*] OS Distro: "$DIST
echo "[*] Distro Version: "$DIST_VER
echo -e $Yellow
echo "[*] Installing dependencies..."
echo -e $Color_Off 

if [ $DIST == "Ubuntu" ] || [ $DIST == "Kali" ] || [ $DIST == "Mint" ] || [ $DIST == "Debian" ]
then
	sudo apt-get update
	sudo apt-get install -y golang nasm mingw-w64-i686-dev mingw-w64-tools mingw-w64-x86-64-dev mingw-w64-common mingw-w64 mingw-ocaml gcc-multilib g++-multilib
elif [ $DIST == "Arch" ] || [ $DIST == "Manjaro" ]
then
	pacman -S --noconfirm go nasm mingw-w64-binutils mingw-w64-crt mingw-w64-gcc mingw-w64-headers mingw-w64-headers-bootstrap mingw-w64-winpthreads gcc-multilib
elif [ $DIST == "MacOS" ]
then
	# MacOS install...
	brew install mingw-w64
	brew install go
	brew install gcc-multilib

elif [ $DIST == "Unknown" ]
then
	echo -e -n $BRed
	echo "[!] OS not supported :("
	echo -e $Color_Off
	return 1	
fi

echo -e $Yellow
echo "[*] Building source..."

cd ../lib
export GOPATH=$(pwd)
cd ..

echo "[*] GOPATH=$GOPATH"
echo -e -n $Color_Off
go build -ldflags "-s -w" handler.go
cd src
go build -ldflags "-s -w" -o ../amber
cd ..
echo "[*] Creating Amber directory..."
sudo mkdir /usr/share/Amber
sudo chown $(whoami):$(whoami) /usr/share/Amber
cp -r * /usr/share/Amber/

sudo ln -s /usr/share/Amber/amber /usr/local/bin/amber
sudo chown $(whoami):$(whoami) /usr/local/bin/amber
sudo ln -s /usr/share/Amber/handler /usr/local/bin/amber_handler
sudo chown $(whoami):$(whoami) /usr/local/bin/amber_handler

echo -e $BGreen
echo "[✔] Setup completed !"
echo -e $Color_Off

