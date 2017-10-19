# AMBER ![Version](https://img.shields.io/badge/version-1.0-brightgreen.svg) [![License](https://img.shields.io/packagist/l/doctrine/orm.svg)](https://raw.githubusercontent.com/EgeBalci/Amber/master/LICENSE) [![Golang](https://img.shields.io/badge/Golang-1.9-blue.svg)](https://golang.org) [![Twitter](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/egeblc)	

[![Banner](https://github.com/EgeBalci/Amber/raw/master/Banner.png)](https://github.com/egebalci/Amber)


Amber is a proof of concept packer, it can pack regularly compiled PE files into reflective PE files that can be used as multi stage infection payloads. If you want to learn the packing methodology used inside the Amber check out below. 

PS: This is not a complete tool some things may break so take it easy on the issues :sweat_smile: and feel free to contribute.

# REFLECTIVE PE PACKING WITH AMBER

[BLOG POST](https://pentest.blog/packing-reflective-pe-files-with-amber)
<br>
[PAPER](https://raw.githubusercontent.com/EgeBalci/Amber/master/PAPER.pdf)

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

		# Version: 1.0.0
		# Source: github.com/egebalci/Amber


		USAGE: 
		  amber file.exe [options]


		OPTIONS:
		  
		  -k, --key       [string]        Custom cipher key
		  -ks,--keysize   <length>        Size of the encryption key in bytes (Max:100/Min:4)
		  --staged                        Generated a staged payload
		  --iat                           Uses import address table entries instead of hash api
		  --no-resource                   Don't add any resource
		  -v, --verbose                   Verbose output mode
		  -h, --help                      Show this massage

		EXAMPLE:
		  (Default settings if no option parameter passed)
		  amber file.exe -ks 8


<div align="center">
	<a href="http://www.youtube.com/watch?feature=player_embedded&v=ZeauXofZw-g" target="_blank">
		<img src="http://img.youtube.com/vi/ZeauXofZw-g/0.jpg" alt="DEMO1" width="500" height="400" border="10" />
	</a>
</div>


# DETECTION




# TODO

- [ ] Add x64 support
- [ ] .NET file support
- [ ] Add a IAT parser shellcode to stub
- [ ] Add yara rules to repo
- [ ] Write a unpacker for Amber payloads
- [ ] Automate IAT index address finding on --iat option
- [ ] Add assembly encoder & anti debug features
- [ ] Add more integrity checks to the file mapping function