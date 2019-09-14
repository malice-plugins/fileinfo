# fileinfo

[![Circle CI](https://circleci.com/gh/malice-plugins/fileinfo.png?style=shield)](https://circleci.com/gh/malice-plugins/fileinfo) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org) [![Docker Stars](https://img.shields.io/docker/stars/malice/fileinfo.svg)](https://hub.docker.com/r/malice/fileinfo/) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/fileinfo.svg)](https://hub.docker.com/r/malice/fileinfo/) [![Docker Image](https://img.shields.io/badge/docker%20image-122MB-blue.svg)](https://hub.docker.com/r/malice/fileinfo/)

Malice File Info Plugin ([exiftool](https://www.sno.phy.queensu.ca/~phil/exiftool/), [TRiD](http://mark0.net/soft-trid-e.html) and [ssdeep](https://github.com/ssdeep-project/ssdeep))

> This repository contains a **Dockerfile** of the FileInfo malice plugin **malice/fileinfo**.

---

### Dependencies

- [ubuntu:bionic (_84.1 MB_\)](https://hub.docker.com/_/ubuntu/)

## Installation

1. Install [Docker](https://www.docker.io/).
2. Download [trusted build](https://hub.docker.com/r/malice/fileinfo/) from public [Docker Registry](https://index.docker.io/): `docker pull malice/fileinfo`

## Usage

```bash
$ docker run -v /path/to/malware:/malware malice/fileinfo FILE

Usage: fileinfo [OPTIONS] COMMAND [arg...]

Malice File Info Plugin - ssdeep/exiftool/TRiD

Version: , BuildTime: 20180902

Author:
  blacktop - <https://github.com/blacktop>

Options:
  --verbose, -V          verbose output
  --table, -t            output as Markdown table
  --mime, -m             output only mimetype
  --callback, -c         POST results to Malice webhook [$MALICE_ENDPOINT]
  --proxy, -x            proxy settings for Malice webhook endpoint [$MALICE_PROXY]
  --elasticsearch value  elasticsearch url for Malice to store results [$MALICE_ELASTICSEARCH_URL]
  --timeout value        malice plugin timeout (in seconds) (default: 10) [$MALICE_TIMEOUT]
  --help, -h             show help
  --version, -v          print the version

Commands:
  web   Create a File Info web service
  help  Shows a list of commands or help for one command

Run 'fileinfo COMMAND --help' for more information on a command.
```

## Sample Output

### [JSON](https://github.com/malice-plugins/fileinfo/blob/master/docs/results.json)

```json
{
  "magic": {
    "mime": "application/x-executable",
    "description": "ELF 64-bit LSB executable, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib64/ld-linux-x86-64.so.2, for GNU/Linux 2.6.26, BuildID[sha1]=8ffd894e500a9f125b32fa8a3f700f0f710961de, stripped"
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

### [Markdown](https://github.com/malice-plugins/fileinfo/blob/master/docs/SAMPLE.md)

#### Magic

| Field       | Value                                             |
| ----------- | ------------------------------------------------- |
| Mime        | application/x-dosexec                             |
| Description | PE32 executable (GUI) Intel 80386, for MS Windows |

#### SSDeep

- `768:15jQ4nVHQaeO379u4XckKVCsknBN9A4hUnDxDiNZ957ZpK0IUUiM95Zdz:15jQ4nVHQaeO9uwckKuBN9A4UnDxcbFi`

#### TRiD

- 30.4% (.EXE) Win32 Executable MS Visual C&#43;&#43; (generic) (31206/45/13)
- 26.9% (.EXE) Win64 Executable (generic) (27625/18/4)
- 25.9% (.EXE) Win32 EXE Yoda&#39;s Crypter (26569/9/4)
- 6.4% (.DLL) Win32 Dynamic Link Library (generic) (6578/25/2)
- 4.3% (.EXE) Win32 Executable (generic) (4508/7/1)

#### Exiftool

| Field                 | Value                 |
| --------------------- | --------------------- |
| CharacterSet          | Unicode               |
| CodeSize              | 20480                 |
| Comments              |                       |
| CompanyName           | Microsoft Corporation |
| EntryPoint            | 0x5a46                |
| ExifToolVersionNumber | 11.06                 |

...`SNIP`...

---

## Documentation

- [To write results to ElasticSearch](https://github.com/malice-plugins/fileinfo/blob/master/docs/elasticsearch.md)
- [To create a File Info micro-service](https://github.com/malice-plugins/fileinfo/blob/master/docs/web.md)
- [To post results to a webhook](https://github.com/malice-plugins/fileinfo/blob/master/docs/callback.md)

## Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/malice-plugins/fileinfo/issues/new)

## CHANGELOG

See [`CHANGELOG.md`](https://github.com/malice-plugins/fileinfo/blob/master/CHANGELOG.md)

## Contributing

[See all contributors on GitHub](https://github.com/malice-plugins/fileinfo/graphs/contributors).

Please update the [CHANGELOG.md](https://github.com/malice-plugins/fileinfo/blob/master/CHANGELOG.md) and submit a [Pull Request on GitHub](https://help.github.com/articles/using-pull-requests/).

## License

MIT Copyright (c) 2016-2018 **blacktop**
