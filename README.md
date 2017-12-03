

[![Banner](https://github.com/EgeBalci/Amber/raw/master/Banner.png)](https://github.com/egebalci/Amber)

[![Version](https://img.shields.io/badge/version-1.1.0-green.svg)](https://github.com/egebalci/Amber) [![License](https://img.shields.io/packagist/l/doctrine/orm.svg)](https://raw.githubusercontent.com/EgeBalci/Amber/master/LICENSE) [![Golang](https://img.shields.io/badge/Golang-1.9-blue.svg)](https://golang.org) [![Twitter](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/egeblc)


Amber is a proof of concept packer, it can pack regularly compiled PE files into reflective PE files that can be used as multi stage infection payloads. If you want to learn the packing methodology used inside the Amber check out below. 

PS: This is not a complete tool some things may break so take it easy on the issues :sweat_smile: and feel free to contribute.


Developed By Ege Balcı from [INVICTUS](https://invictuseurope.com)/[PRODAFT](https://prodaft.com).

# REFLECTIVE PE PACKING WITH AMBER

<br>

<a href="https://pentest.blog/introducing-new-packing-method-first-reflective-pe-packer" target="_blank">
		<img height="250" align="left" src="https://pentest.blog/wp-content/uploads/68747470733a2f2f696d6167652e6962622e636f2f66426e51566d2f70656e746573745f626c6f67332e6a7067.jpeg" alt="DEMO1"  />
</a>
<a href="https://raw.githubusercontent.com/EgeBalci/Amber/master/PAPER.pdf"></a>
<a href="https://github.com/EgeBalci/Amber/raw/master/PAPER.pdf">
	<img align="right" src="https://pentest.blog/wp-content/uploads/pdf2.png"/>
</a>

<br>
<br>
<br>
<br>
<br>
<br>
<br>
<br>


# INSTALLATION

SUPPORTED PLATFORMS:
<table>
    <tr>
        <th>Operative system</th>
        <th> Version </th>
    </tr>
    <tr>
        <td>Ubuntu</td>
        <td> 16.04\16.10\17.04 </td>
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
        <td>Debian</td>
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
		//  POC Reflective PE Packer                                             

		# Version: 1.1.0
		# Source: github.com/egebalci/Amber


		USAGE: 
		  amber [options] file.exe


		OPTIONS:
		  
		  -k, --key       [string]        Custom cipher key
		  -ks,--keysize   <length>        Size of the encryption key in bytes (Max:255/Min:8)
		  --staged                        Generated a staged payload
		  --iat                           Uses import address table entries instead of hash api
		  --no-resource                   Don't add any resource
		  -v, --verbose                   Verbose output mode
		  -h, --help                      Show this massage

		EXAMPLE:
		  (Default settings if no option parameter passed)
		  amber -ks 8 file.exe


<strong>Fileless ransomware deployment with powershell</strong>

<div align="center">
	<a href="https://www.youtube.com/watch?v=JVv_spX6D4U" target="_blank">
		<img src="http://img.youtube.com/vi/JVv_spX6D4U/0.jpg" alt="DEMO1" width="500" height="400"/>
	</a>
</div>

<strong>Multi Stage EXE deployment with metasploit stagers</strong>

<div align="center">
	<a href="https://www.youtube.com/watch?v=3en0ftnjEpE" target="_blank">
		<img src="http://img.youtube.com/vi/3en0ftnjEpE/0.jpg" alt="DEMO1" width="500" height="400"/>
	</a>
</div>


# DETECTION
Current detection rate (19.10.2017) of the POC packer is pretty satisfying but since this is going to be a public project current detection score will rise inevitably :)

When no extra parameters passed (only the file name) packer generates a multi stage payload and performs an basic XOR cipher with a multi byte random key then compiles it into a EXE file with adding few extra anti detection functions. Generated EXE file executes the stage payload like a regular shellcode after deciphering the payload and making the required environmental checks. This particular sample is the mimikats.exe (sha256 - 9369b34df04a2795de083401dda4201a2da2784d1384a6ada2d773b3a81f8dad) file packed with a 12 byte XOR key (./amber mimikats.exe -ks 12).  The detection rate of the mimikats.exe file before packing is 51/66 on VirusTotal. In this particular example packer uses the default way to find the windows API addresses witch is using the hash API, avoiding the usage of hash API will decrease the detection rate. Currently packer supports the usage of fixed addresses of  IAT offsets also next versions will include IAT parser shellcodes for more alternative API address finding methods.

<strong>VirusTotal</strong> (5/65)

[![VirusTotal](https://pentest.blog/wp-content/uploads/VirusTotal-1.png)](https://www.virustotal.com/#/file/3330d02404c56c1793f19f5d18fd5865cadfc4bd015af2e38ed0671f5e737d8a/detection)

<strong>VirusCheckmate</strong> (0/36)

[![VirusCheckmate](https://pentest.blog/wp-content/uploads/VirusCheckmate.png)](http://viruscheckmate.com/id/1ikb99sNVrOM)

<strong>NoDistribute</strong> (0/36)

[![NoDistribute](https://NoDistribute.com/result/image/7uMa96SNOY13rtmTpW5ckBqzAv.png)](https://NoDistribute.com/result/image/7uMa96SNOY13rtmTpW5ckBqzAv.png)



# TODO

- [ ] Add x64 support
- [ ] Add MacOS support
- [ ] Add.NET file support
- [ ] Add a IAT parser shellcode to stub
- [x] Add yara rules to repo
- [ ] Write a unpacker for Amber payloads
- [x] Add RC4 encryption to payloads
- [ ] Automate IAT index address finding on --iat option
- [ ] Add more integrity checks to the file mapping function
- [x] Better installation mechanism
