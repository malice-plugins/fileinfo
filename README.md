# malice-fileinfo

[![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)

Malice File Info Plugin

This repository contains a **Dockerfile** of the **Malice File Info Plugin** for [Docker](https://www.docker.io/)'s [trusted build](https://index.docker.io/u/malice/fileinfo/) published to the public [Docker Registry](https://index.docker.io/).

### Dependencies

* [alpine:edge](https://hub.docker.com/_/alpine/)

### Image Size
[![](https://badge.imagelayers.io/malice/fileinfo:latest.svg)](https://imagelayers.io/?images=malice/fileinfo:latest 'Get your own badge on imagelayers.io')

### Image Tags
```bash
$ docker images

REPOSITORY          TAG                 VIRTUAL SIZE
malice/fileinfo     latest              688   MB
```

### Installation

1. Install [Docker](https://www.docker.io/).

2. Download [trusted build](https://index.docker.io/u/blacktop/elk/) from public [Docker Registry](https://index.docker.io/): `docker pull blacktop/elk`

#### Alternatively, build an image from Dockerfile
```bash
$ docker build -t blacktop/elk github.com/blacktop/docker-elk
```
### Usage
```bash
$ docker run malice/fileinfo FILE
```
This will output to stdout and POST to malice results API webhook endpoint.

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
