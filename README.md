

[![Banner](https://github.com/EgeBalci/Amber/raw/master/Banner.png)](https://github.com/egebalci/Amber)

[![Version](https://img.shields.io/badge/version-1.3.0-green.svg)](https://github.com/egebalci/Amber) [![License](https://img.shields.io/packagist/l/doctrine/orm.svg)](https://raw.githubusercontent.com/EgeBalci/Amber/master/LICENSE) [![Golang](https://img.shields.io/badge/Golang-1.9-blue.svg)](https://golang.org) [![Twitter](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/egeblc)


Amber is a proof of concept packer for stealthy payload deployment and anti-virus evasion purposes, it can pack regularly compiled PE files into reflective PE files that can be used as multi stage infection payloads. If you want to learn the packing methodology used inside the Amber check out below. For more detail about usage, installation  and how to decrease detection rate check out [WIKI](https://github.com/EgeBalci/Amber/wiki).


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

<br><br><br><br><br><br><br><br>

# INSTALLATION

SUPPORTED PLATFORMS:
<table>
    <tr>
        <th>Operating system</th>
        <th>Tested Version</th>
    </tr>
    <tr>
        <td>Ubuntu</td>
        <td>16.04\16.10\17.04\17.08</td>
    </tr>
    <tr>
        <td>Kali linux</td>
        <td>2018.1\2017.1</td>
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
        <td>9.2</td>
    </tr>
</table>

<strong>BUILD</strong>

For compiling from source running the setup file will be enough.

        git clone https://github.com/egebalci/Amber.git
		cd Amber/setup/
        ./setup.sh

<strong>DOCKER</strong>

		docker pull egee/amber
		docker run -it egee/amber

# USAGE

        USAGE: 
        amber [options] file.exe
        OPTIONS:
        -k, -keysize                Size of the encryption key in bytes (Max:255/Min:8)
        -r, -reflective             Generated a reflective payload
        -i, -iat                    Uses import address table entries instead of export address table
        -s, -scrape                 Scrape the PE header info (May break some files)
        -no-resource                Don't add any resource data
        -ignore-mapping-size        Ignore mapping size mismatch errors
        -ignore-section-alignment   Ignore broken section alignment errors
        -v, -verbose                Verbose output mode
        -h, -H                      Show this massage
        EXAMPLE:
        (Default settings if no option parameter passed)
        amber -k 8 file.exe
<strong>On Docker</strong><br>
		`docker run -it -v /tmp/:/tmp/ amber /tmp/file.exe`

# EXAMPLE USAGE
<br>
<a href="https://www.youtube.com/watch?v=JVv_spX6D4U" target="_blank">
	<img src="http://img.youtube.com/vi/JVv_spX6D4U/0.jpg" alt="DEMO1" width="400" height="300" align="right"/>
</a>

<a href="https://www.youtube.com/watch?v=3en0ftnjEpE" target="_blank">
	<img src="https://pentest.blog/wp-content/uploads/Screenshot-at-2018-02-23-22-42-18-2-1024x704.png" alt="DEMO1" width="400" height="300" align="left"/>
</a><br><br><br>


