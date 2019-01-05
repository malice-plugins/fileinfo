####################################################
# GOLANG BUILDER
####################################################
FROM golang:1.11 as go_builder
RUN apt-get update && apt-get install -y libmagic-dev libc6
COPY . /go/src/github.com/maliceio/malice-fileinfo
WORKDIR /go/src/github.com/maliceio/malice-fileinfo
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN go build -ldflags "-s -w -X main.Version=v$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/info

####################################################
# FILEINFO BUILDER
####################################################
FROM ubuntu:bionic

LABEL maintainer "https://github.com/blacktop"

LABEL malice.plugin.repository = "https://github.com/malice-plugins/fileinfo.git"
LABEL malice.plugin.category="metadata"
LABEL malice.plugin.mime="*"
LABEL malice.plugin.docker.engine="*"

# Create a malice user and group first so the IDs get set the same way, even as
# the rest of this may change over time.
RUN groupadd -r malice \
  && useradd --no-log-init -r -g malice malice \
  && mkdir /malware \
  && chown -R malice:malice /malware

ENV SSDEEP 2.14.1
ENV EXIFTOOL 11.11

RUN buildDeps='ca-certificates \
  build-essential \
  openssl \
  unzip \
  curl' \
  && set -x \
  && apt-get update -qq \
  && apt-get install -yq --no-install-recommends $buildDeps libmagic-dev libc6 \
  && echo "Downloading TRiD and Database..." \
  && curl -Ls http://mark0.net/download/trid_linux_64.zip > /tmp/trid_linux_64.zip \
  && curl -Ls http://mark0.net/download/triddefs.zip > /tmp/triddefs.zip \
  && cd /tmp \
  && unzip trid_linux_64.zip \
  && unzip triddefs.zip \
  && chmod +x trid \
  && mv trid /usr/bin/ \
  && mv triddefs.trd /usr/bin/ \
  && echo "Installing ssdeep..." \
  && curl -Ls https://github.com/ssdeep-project/ssdeep/releases/download/release-$SSDEEP/ssdeep-$SSDEEP.tar.gz > \
  /tmp/ssdeep-$SSDEEP.tar.gz \
  && cd /tmp \
  && tar xzf ssdeep-$SSDEEP.tar.gz \
  && cd ssdeep-$SSDEEP \
  && ./configure \
  && make \
  && make install \
  && echo "Installing exiftool..." \
  && curl -Ls https://www.sno.phy.queensu.ca/~phil/exiftool/Image-ExifTool-$EXIFTOOL.tar.gz > \
  /tmp/exiftool.tar.gz \
  && cd /tmp \
  && tar xzf exiftool.tar.gz \
  && cd Image-ExifTool-$EXIFTOOL \
  && perl Makefile.PL \
  && make test \
  && make install \
  && echo "Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /root/.gnupg

RUN apt-get update -qq && apt-get install -yq --no-install-recommends ca-certificates gosu \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=go_builder /bin/info /bin/info

WORKDIR /malware

ENTRYPOINT ["gosu","malice","info"]
CMD ["--help"]

####################################################
####################################################