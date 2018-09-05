####################################################
# GOSU BUILDER
####################################################
FROM ubuntu:xenial as gosu_builder

ENV GOSU_VERSION 1.10
RUN set -ex; \
  \
  fetchDeps=' \
  ca-certificates \
  dirmngr \
  wget \
  '; \
  apt-get update; \
  apt-get install -y --no-install-recommends $fetchDeps; \
  rm -rf /var/lib/apt/lists/*; \
  \
  dpkgArch="$(dpkg --print-architecture | awk -F- '{ print $NF }')"; \
  wget -O /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch"; \
  wget -O /usr/local/bin/gosu.asc "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch.asc"; \
  \
  # verify the signature
  export GNUPGHOME="$(mktemp -d)"; \
  export GPG_KEYS=B42F6819007F00F88E364FD4036A9C25BF357DD4; \
  ( gpg --keyserver ha.pool.sks-keyservers.net --recv-keys "$GPG_KEYS" \
  || gpg --keyserver pgp.mit.edu --recv-keys "$GPG_KEYS" \
  || gpg --keyserver keyserver.pgp.com --recv-keys "$GPG_KEYS" ); \
  gpg --batch --verify /usr/local/bin/gosu.asc /usr/local/bin/gosu; \
  rm -r "$GNUPGHOME" /usr/local/bin/gosu.asc; \
  \
  chmod +x /usr/local/bin/gosu; \
  # verify that the binary works
  gosu nobody true

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
FROM ubuntu:xenial

LABEL maintainer "https://github.com/blacktop"

LABEL malice.plugin.repository = "https://github.com/malice-plugins/fileinfo.git"
LABEL malice.plugin.category="metadata"
LABEL malice.plugin.mime="*"
LABEL malice.plugin.docker.engine="*"

ENV SSDEEP 2.14.1
ENV EXIFTOOL 11.10

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

RUN apt-get update -qq && apt-get install -yq --no-install-recommends ca-certificates \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=gosu_builder /usr/local/bin/gosu /usr/local/bin/gosu
COPY --from=go_builder /bin/info /bin/info

# Create a malice user and group first so the IDs get set the same way, even as
# the rest of this may change over time.
RUN groupadd -r malice \
  && useradd --no-log-init -r -g malice malice \
  && mkdir /malware \
  && chown -R malice:malice /malware

WORKDIR /malware

ENTRYPOINT ["gosu","malice","info"]
CMD ["--help"]

####################################################
####################################################