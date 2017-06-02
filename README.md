# Amber [![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://raw.githubusercontent.com/EgeBalci/Amber/master/LICENSE) [![Golang](https://img.shields.io/badge/Go-1.6-blue.svg)](https://golang.org)	
POC crypter for ReplaceProcess method.

![](https://github.com/EgeBalci/Amber/raw/master/amber.ico)

# INSTALLATION

SUPPORTED PLATFORMS:
<table>
    <tr>
        <th>Operative system</th>
        <th> Version </th>
    </tr>
    <tr>
        <td>Ubuntu</td>
        <td> * </td>
    </tr>
    <tr>
        <td>Kali linux</td>
        <td> * </td>
    </tr>
    <tr>
        <td>Manjaro</td>
        <td> * </td>
    </tr>
    <tr>
        <td>Arch Linux</td>
        <td> * </td>
    </tr>
    <tr>
        <td>Black Arch</td>
        <td> * </td>
    </tr>
    <tr>
        <td>Parrot OS</td>
        <td> * </td>
    </tr>
</table>


		sudo chmod +x Setup.sh
		sudo ./Setup.sh
# USAGE


		//   █████╗ ███╗   ███╗██████╗ ███████╗██████╗ 
		//  ██╔══██╗████╗ ████║██╔══██╗██╔════╝██╔══██╗
		//  ███████║██╔████╔██║██████╔╝█████╗  ██████╔╝
		//  ██╔══██║██║╚██╔╝██║██╔══██╗██╔══╝  ██╔══██╗
		//  ██║  ██║██║ ╚═╝ ██║██████╔╝███████╗██║  ██║
		//  ╚═╝  ╚═╝╚═╝     ╚═╝╚═════╝ ╚══════╝╚═╝  ╚═╝
		//  POC Crypter For ReplaceProcess                                             

		# Version: 1.0.0
		# Source: github.com/EgeBalci/Amber


		USAGE: 
		  amber file.exe [options]


		OPTIONS:
		  
		  -k, --key       [string]        Custom cipher key
		  -ks,--keysize   <length>        Size of the encryption key in bytes (Max:100/Min:4)
		  -v, --verbose                   Verbose output mode
		  -h, --help                      Show this massage

		EXAMPLE:
		  (Default settings if no option parameter passed)
		  amber file.exe -ks 8 -o crypted.exe
