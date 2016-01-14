# malice-fileinfo

[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)
[![Docker Stars](https://img.shields.io/docker/stars/malice/fileinfo.svg)][hub]
[![Docker Pulls](https://img.shields.io/docker/pulls/malice/fileinfo.svg)][hub]
[![Image Size](https://img.shields.io/imagelayers/image-size/malice/fileinfo/latest.svg)](https://imagelayers.io/?images=malice/fileinfo:latest)
[![Image Layers](https://img.shields.io/imagelayers/layers/malice/fileinfo/latest.svg)](https://imagelayers.io/?images=malice/fileinfo:latest)

Malice File Info Plugin

This repository contains a **Dockerfile** of the **Malice File Info Plugin** for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/malice/fileinfo/) published to the public [Docker Registry](https://index.docker.io/).

### Dependencies

* [debian:wheezy (*85 MB*)](https://index.docker.io/_/debian/)

### Installation

1. Install [Docker](https://www.docker.io/).

2. Download [trusted build](https://index.docker.io/u/malice/fileinfo/) from public [Docker Registry](https://index.docker.io/): `docker pull malice/fileinfo`

#### Alternatively, build an image from Dockerfile
```bash
$ docker build -t malice/fileinfo github.com/maliceio/malice-fileinfo
```
### Usage
```bash
$ docker run -v /path/to/file:/malware malice/fileinfo FILE
```
This will output to stdout and POST to malice results API webhook endpoint.

### Output JSON:
```json
{
  "ssdeep": "768:C7tsNKQhyl96U9eJqaZ2e5ofMolkcksNmisf4BB5iqboecL027:DkXe1UHfM4N3sfezcL0",
  "trid": [
    "50.1% (.) ELF Executable and Linkable format (Linux) (4025/14)",
    "49.8% (.O) ELF Executable and Linkable format (generic) (4000/1)"
  ],
  "exiftool": {
    "CPU Architecture": "64 bit",
    "CPU Byte Order": "Little endian",
    "CPU Type": "AMD x86-64",
    "ExifTool Version Number": "10.08",
    "File Size": "51 kB",
    "File Type": "ELF executable",
    "File Type Extension": "",
    "MIME Type": "application/octet-stream",
    "Object File Type": "Executable file"
  }
}
```
### Output STDOUT (Markdown Table):
---
#### SSDeep
768:C7tsNKQhyl96U9eJqaZ2e5ofMolkcksNmisf4BB5iqboecL027:DkXe1UHfM4N3sfezcL0
#### Exiftool
| Field                   | Value                    |
| ----------------------- | ------------------------ |
| File Type               | ELF executable           |
| CPU Architecture        | 64 bit                   |
| CPU Type                | AMD x86-64               |
| File Size               | 51 kB                    |
| ExifTool Version Number | 10.09                    |
| File Type Extension     |                          |
| MIME Type               | application/octet-stream |
| CPU Byte Order          | Little endian            |
| Object File Type        | Executable file          |
#### TRiD
 -  50.1% (.) ELF Executable and Linkable format (Linux) (4025/14)
 -  49.8% (.O) ELF Executable and Linkable format (generic) (4000/1)
---
### To Run on OSX
 - Install [Homebrew](http://brew.sh)

```bash
$ brew install caskroom/cask/brew-cask
$ brew cask install virtualbox
$ brew install docker
$ brew install docker-machine
$ docker-machine create --driver virtualbox malice
$ eval $(docker-machine env malice)
```

### Documentation

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-fileinfo/issues/new) and I'll get right on it.

### Credits

### License
MIT Copyright (c) 2016 **blacktop**

[hub]: https://hub.docker.com/r/malice/fileinfo/
