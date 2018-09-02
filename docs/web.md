# Create a File Info scan micro-service

```bash
$ docker run -d -p 3993:3993 malice/fileinfo web

INFO[0000] web service listening on port :3993
```

## Now you can perform scans like so

```bash
$ http -f localhost:3993/scan malware@/path/to/evil/malware
```

> **NOTE:** I am using **httpie** to POST to the malice micro-service

```bash
HTTP/1.1 200 OK
Content-Length: 124
Content-Type: application/json; charset=UTF-8
Date: Sat, 21 Jan 2017 05:39:29 GMT

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
