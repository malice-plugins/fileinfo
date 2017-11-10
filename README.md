malice-fileinfo
===============

[![Circle CI](https://circleci.com/gh/malice-plugins/fileinfo.png?style=shield)](https://circleci.com/gh/malice-plugins/fileinfo) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org) [![Docker Stars](https://img.shields.io/docker/stars/malice/fileinfo.svg)](https://hub.docker.com/r/malice/fileinfo/) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/fileinfo.svg)](https://hub.docker.com/r/malice/fileinfo/) [![Docker Image](https://img.shields.io/badge/docker%20image-166MB-blue.svg)](https://hub.docker.com/r/malice/fileinfo/)

Malice File Info Plugin (exiftool, TRiD and ssdeep)

This repository contains a **Dockerfile** of the FileInfo malice plugin **malice/fileinfo**.

### Dependencies

-	[ubuntu:xenial (*122 MB*\)](https://index.docker.io/_/debian/)

### Installation

1.	Install [Docker](https://www.docker.io/).
2.	Download [trusted build](https://hub.docker.com/r/malice/fileinfo/) from public [Docker Registry](https://index.docker.io/): `docker pull malice/fileinfo`

### Usage

```bash
$ docker run -v /path/to/malware:/malware malice/fileinfo FILE

Usage: fileinfo [OPTIONS] COMMAND [arg...]

Malice File Info Plugin - ssdeep/TRiD/exiftool

Version: v0.1.0, BuildTime: 20171110

Author:
  blacktop - <https://github.com/blacktop>

Options:
  --verbose, -V         verbose output
  --table, -t           output as Markdown table
  --mime, -m		    output only mimetype
  --callback, -c	    POST results to Malice webhook [$MALICE_ENDPOINT]
  --proxy, -x           proxy settings for Malice webhook endpoint [$MALICE_PROXY]
  --timeout value       malice plugin timeout (in seconds) (default: 10) [$MALICE_TIMEOUT]
  --elasitcsearch value elasitcsearch address for Malice to store results [$MALICE_ELASTICSEARCH]
  --help, -h            show help
  --version, -v         print the version

Commands:
  web       Create a File Info scan web service  
  help		Shows a list of commands or help for one command

Run 'fileinfo COMMAND --help' for more information on a command.
```

Sample Output
-------------

### JSON:

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

Or click here [results.json](https://github.com/maliceio/malice-fileinfo/blob/master/docs/results.json)

### Markdown:

Click here [SAMPLE.md](https://github.com/maliceio/malice-fileinfo/blob/master/docs/SAMPLE.md)

---

Documentation
-------------

-	[To write results to ElasticSearch](https://github.com/maliceio/malice-fileinfo/blob/master/docs/elasticsearch.md)
-	[To create a File Info micro-service](https://github.com/maliceio/malice-fileinfo/blob/master/docs/web.md)
-	[To post results to a webhook](https://github.com/maliceio/malice-fileinfo/blob/master/docs/callback.md)

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice-fileinfo/issues/new)

### CHANGELOG

See [`CHANGELOG.md`](https://github.com/maliceio/malice-fileinfo/blob/master/CHANGELOG.md)

### Contributing

[See all contributors on GitHub](https://github.com/maliceio/malice-fileinfo/graphs/contributors).

Please update the [CHANGELOG.md](https://github.com/maliceio/malice-fileinfo/blob/master/CHANGELOG.md) and submit a [Pull Request on GitHub](https://help.github.com/articles/using-pull-requests/).

### License

MIT Copyright (c) 2016-2017 **blacktop**
