

[![Banner](https://github.com/EgeBalci/amber/raw/master/banner.png)](https://github.com/egebalci/amber)

[![Version](https://img.shields.io/badge/version-2.0.0-green.svg)](https://github.com/egebalci/amber) [![License](https://img.shields.io/packagist/l/doctrine/orm.svg)](https://raw.githubusercontent.com/EgeBalci/amber/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/egebalci/amber)](https://goreportcard.com/report/github.com/egebalci/amber) [![Twitter](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/egeblc)


amber is a reflective PE packer for bypassing security products and mitigations. It can pack regularly compiled PE files into reflective payloads that can load and execute itself like a shellcode. It enables stealthy in-memory payload deployment that can be used to bypass anti-virus, firewall, IDS, IPS products and application white-listing mitigations.  If you want to learn more about the packing methodology used inside amber check out below. For more detail about usage, installation and how to decrease detection rate check out [WIKI](https://github.com/egebalci/amber/wiki).


Developed By Ege Balcı from [INVICTUS](https://invictuseurope.com)/[PRODAFT](https://prodaft.com).

# REFLECTIVE PE PACKING WITH AMBER

<br>

<a href="https://pentest.blog/introducing-new-packing-method-first-reflective-pe-packer" target="_blank">
		<img height="250" align="left" src="https://pentest.blog/wp-content/uploads/68747470733a2f2f696d6167652e6962622e636f2f66426e51566d2f70656e746573745f626c6f67332e6a7067.jpeg" alt="DEMO1"  />
</a>
<a href="https://raw.githubusercontent.com/EgeBalci/amber/master/PAPER.pdf"></a>
<a href="https://github.com/EgeBalci/amber/raw/master/PAPER.pdf">
	<img align="right" src="https://pentest.blog/wp-content/uploads/pdf2.png"/>
</a>

<br><br><br><br><br><br><br><br>

# INSTALLATION


***DEPENDENCIES***

- [go](https://golang.org/dl/)
- [NASM](https://www.nasm.us/)

On *nix systems both of the dependencies can be installed with OS packet managers. (APT/PACMAN/YUM)


Get one of the pre-build release [here](https://github.com/egebalci/amber/releases). Or get it with following alternatives.

***GO (suggested)***
```
go get github.com/egebalci/amber
```

***BLACKARCH INSTALL***     
```
sudo pacman -S amber
```

***DOCKER INSTALL***

[![Docker](http://dockeri.co/image/egee/amber)](https://hub.docker.com/r/egee/amber/)

```
docker pull egee/amber
docker run -it egee/amber
```

# USAGE
```
USAGE: 
  amber [options] file.exe
OPTIONS:
  -k, -keysize                Size of the encryption key in bytes (Max:255/Min:8)
  -r, -reflective             Generated a reflective payload
  -a, -anti-analysis          Add anti-analysis measures
  -i, -iat                    Use import address table entries instead of export address table
  -s, -scrape                 Scrape the PE header info (May break some files)
  -no-resource                Don't add any resource data (removes icon)
  -ignore-integrity           Ignore integrity check errors
  -v, -verbose                Verbose output mode
  -h, -H                      Show this massage
EXAMPLE:
  (Default settings if no option parameter passed)
  amber -k 8 file.exe
```

***Docker Usage***
```
docker run -it -v /tmp/:/tmp/ amber /tmp/file.exe
```

# EXAMPLE USAGE

- <strong>NOPcon 2018 [DEMO](https://www.youtube.com/watch?v=lCPdKSH6RMc)</strong>

<br><br>

<a href="https://www.youtube.com/watch?v=JVv_spX6D4U" target="_blank">
	<img src="http://img.youtube.com/vi/JVv_spX6D4U/0.jpg" alt="DEMO1" width="400" height="300" align="right"/>
</a>

<a href="https://www.youtube.com/watch?v=3en0ftnjEpE" target="_blank">
	<img src="https://pentest.blog/wp-content/uploads/Screenshot-at-2018-02-23-22-42-18-2-1024x704.png" alt="DEMO1" width="400" height="300" align="left"/>
</a><br><br><br>
<br>