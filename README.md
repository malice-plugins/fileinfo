malice-fileinfo
===============

[![Circle CI](https://circleci.com/gh/maliceio/malice-fileinfo.png?style=shield)](https://circleci.com/gh/maliceio/malice-fileinfo) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org) [![Docker Stars](https://img.shields.io/docker/stars/malice/fileinfo.svg)](https://hub.docker.com/r/malice/fileinfo/) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/fileinfo.svg)](https://hub.docker.com/r/malice/fileinfo/) [![Docker Image](https://img.shields.io/badge/docker image-209.6 MB-blue.svg)](https://hub.docker.com/r/malice/fileinfo/)

Malice File Info Plugin (exiftool, TRiD and ssdeep)

This repository contains a **Dockerfile** of the FileInfo malice plugin **malice/fileinfo**.

### Dependencies

-	[debian:wheezy (*85 MB*\)](https://index.docker.io/_/debian/)

### Installation

1.	Install [Docker](https://www.docker.io/).
2.	Download [trusted build](https://hub.docker.com/r/malice/fileinfo/) from public [Docker Registry](https://index.docker.io/): `docker pull malice/fileinfo`

### Usage

```bash
$ docker run -v /path/to/malware:/malware malice/fileinfo FILE

Usage: fileinfo [OPTIONS] COMMAND [arg...]

Malice File Info Plugin - ssdeep/TRiD/exiftool

Version: v0.1.0, BuildTime: 20160114

Author:
  blacktop - <https://github.com/blacktop>

Options:                                                                                                                  
  --verbose, -V         verbose output                                                                                    
  --table, -t           output as Markdown table                                                                          
  --post, -p            POST results to Malice webhook [$MALICE_ENDPOINT]                                                 
  --proxy, -x           proxy settings for Malice webhook endpoint [$MALICE_PROXY]                                        
  --elasitcsearch value elasitcsearch address for Malice to store results [$MALICE_ELASTICSEARCH]                         
  --help, -h            show help                                                                                         
  --version, -v         print the version                                                                                 
                                                                                                                          
Commands:                                                                                                                 
  help  Shows a list of commands or help for one command                                                                  
                                                                                                                          
Run 'fileinfo COMMAND --help' for more information on a command.
```

This will output to stdout and POST to malice results API webhook endpoint.

### Sample Output JSON:

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

### Sample Output STDOUT (Markdown Table):

---

#### SSDeep

768:C7tsNKQhyl96U9eJqaZ2e5ofMolkcksNmisf4BB5iqboecL027:DkXe1UHfM4N3sfezcL0

#### TRiD

-	50.1% (.) ELF Executable and Linkable format (Linux) (4025/14)
-	49.8% (.O) ELF Executable and Linkable format (generic) (4000/1)

#### Exiftool

| Field                   | Value                    |
|-------------------------|--------------------------|
| File Type               | ELF executable           |
| CPU Architecture        | 64 bit                   |
| CPU Type                | AMD x86-64               |
| File Size               | 51 kB                    |
| ExifTool Version Number | 10.09                    |
| File Type Extension     |                          |
| MIME Type               | application/octet-stream |
| CPU Byte Order          | Little endian            |
| Object File Type        | Executable file          |

---

### To write results to [ElasticSearch](https://www.elastic.co/products/elasticsearch)

```bash
$ docker volume create --name malice
$ docker run -d -p 9200:9200 -v malice:/data --name elastic elasticsearch
$ docker run --rm --link elastic malice/fileinfo FILE
```

### Documentation

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-fileinfo/issues/new) and I'll get right on it.

### License

MIT Copyright (c) 2016 **blacktop**
