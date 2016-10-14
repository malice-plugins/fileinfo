malice-fileinfo
===============

[![Circle CI](https://circleci.com/gh/maliceio/malice-fileinfo.png?style=shield)](https://circleci.com/gh/maliceio/malice-fileinfo) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org) [![Docker Stars](https://img.shields.io/docker/stars/malice/fileinfo.svg)](https://hub.docker.com/r/malice/fileinfo/) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/fileinfo.svg)](https://hub.docker.com/r/malice/fileinfo/) [![Docker Image](https://img.shields.io/badge/docker image-456.2 MB-blue.svg)](https://hub.docker.com/r/malice/fileinfo/)

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
  "file": {
    "name": "cat",
    "path": "/bin/cat",
    "size": "51.86 kB",
    "md5": "8fee23b7db38f7fa439ab2a71a5dfef4",
    "sha1": "ea010d65968bdccf995ec7777eb7ef73b8460285",
    "sha256": "41033b3bcc5805c072498bce21d328dae238626e513d5e16bc9f928864a8936e",
    "sha512": "a7cdb673c4d076144a329700e6cc7b2f187a8e1b9ee2142b21b651cbe600af50b44d8c3c9c8fa800c941916568983d205fa5e3205680473943974711efee25dc",
    "mime": "application/x-executable",
    "magic": "ELF 64-bit LSB executable, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib64/ld-linux-x86-64.so.2, for GNU/Linux 2.6.26, BuildID[sha1]=8ffd894e500a9f125b32fa8a3f700f0f710961de, stripped"
  },
  "ssdeep": "768:C7tsNKQhyl96U9eJqaZ2e5ofMolkcksNmisf4BB5iqboecL027:DkXe1UHfM4N3sfezcL0",
  "trid": [
    "50.1% (.) ELF Executable and Linkable format (Linux) (4025/14)",
    "49.8% (.O) ELF Executable and Linkable format (generic) (4000/1)"
  ],
  "exiftool": {
    "CPUArchitecture": "64 bit",
    "CPUByteOrder": "Little endian",
    "CPUType": "AMD x86-64",
    "ExifToolVersionNumber": "10.25",
    "FileSize": "51 kB",
    "FileType": "ELF executable",
    "FileTypeExtension": "",
    "MIMEType": "application/octet-stream",
    "ObjectFileType": "Executable file"
  }
}
```

### Sample Output STDOUT (Markdown Table):

---

#### File
| Field  | Value                                                                                                                                                                                                    |
| ------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Name   | cat                                                                                                                                                                                                      |
| Path   | /bin/cat                                                                                                                                                                                                 |
| Size   | 51.86 kB                                                                                                                                                                                                 |
| MD5    | 8fee23b7db38f7fa439ab2a71a5dfef4                                                                                                                                                                         |
| SHA1   | ea010d65968bdccf995ec7777eb7ef73b8460285                                                                                                                                                                 |
| SHA256 | 41033b3bcc5805c072498bce21d328dae238626e513d5e16bc9f928864a8936e                                                                                                                                         |
| Mime   | application/x-executable                                                                                                                                                                                 |
| Magic  | ELF 64-bit LSB executable, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib64/ld-linux-x86-64.so.2, for GNU/Linux 2.6.26, BuildID[sha1]=8ffd894e500a9f125b32fa8a3f700f0f710961de, stripped |

#### SSDeep
768:C7tsNKQhyl96U9eJqaZ2e5ofMolkcksNmisf4BB5iqboecL027:DkXe1UHfM4N3sfezcL0

#### TRiD
 -  50.1% (.) ELF Executable and Linkable format (Linux) (4025/14)
 -  49.8% (.O) ELF Executable and Linkable format (generic) (4000/1)

#### Exiftool
| Field                 | Value                    |
| --------------------- | ------------------------ |
| CPUByteOrder          | Little endian            |
| ObjectFileType        | Executable file          |
| CPUType               | AMD x86-64               |
| ExifToolVersionNumber | 10.25                    |
| FileType              | ELF executable           |
| MIMEType              | application/octet-stream |
| FileSize              | 51 kB                    |
| FileTypeExtension     |                          |
| CPUArchitecture       | 64 bit                   |

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
