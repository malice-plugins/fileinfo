FROM debian:wheezy

LABEL maintainer "https://github.com/blacktop"

LABEL malice.plugin.repository = "https://github.com/malice-plugins/fileinfo.git"
LABEL malice.plugin.category="metadata"
LABEL malice.plugin.mime="*"
LABEL malice.plugin.docker.engine="*"

# Create a malice user and group first so the IDs get set the same way, even as
# the rest of this may change over time.
RUN groupadd -r malice && useradd -r -g malice malice

ENV GOSU_VERSION 1.10
ENV TINI_VERSION v0.9.0

RUN set -x \
  && apt-get update -qq \
  && apt-get install -y ca-certificates wget \
  && echo "Grab gosu for easy step-down from root..." \
  && dpkgArch="$(dpkg --print-architecture | awk -F- '{ print $NF }')" \
  && wget -O /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch" \
  && wget -O /usr/local/bin/gosu.asc "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch.asc" \
  && export GNUPGHOME="$(mktemp -d)" \
  && gpg --keyserver ha.pool.sks-keyservers.net --recv-keys B42F6819007F00F88E364FD4036A9C25BF357DD4 \
  && gpg --batch --verify /usr/local/bin/gosu.asc /usr/local/bin/gosu \
  && rm -r "$GNUPGHOME" /usr/local/bin/gosu.asc \
  && chmod +x /usr/local/bin/gosu \
  && gosu nobody true \
  && echo "Grab tini for signal processing and zombie killing..." \
  && wget -O /usr/local/bin/tini "https://github.com/krallin/tini/releases/download/$TINI_VERSION/tini" \
  && wget -O /usr/local/bin/tini.asc "https://github.com/krallin/tini/releases/download/$TINI_VERSION/tini.asc" \
  && export GNUPGHOME="$(mktemp -d)" \
  && gpg --keyserver ha.pool.sks-keyservers.net --recv-keys 6380DC428747F6C393FEACA59A84159D7001A4E5 \
  && gpg --batch --verify /usr/local/bin/tini.asc /usr/local/bin/tini \
  && rm -r "$GNUPGHOME" /usr/local/bin/tini.asc \
  && chmod +x /usr/local/bin/tini \
  && tini -h \
  && echo "Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove ca-certificates wget \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENV SSDEEP ssdeep-2.13

RUN buildDeps='ca-certificates \
               build-essential \
               openssl \
               unzip \
               curl' \
  && set -x \
  && echo 'deb http://ftp.us.debian.org/debian/ testing main contrib non-free' >> /etc/apt/sources.list \
  && apt-get update -qq \
  && apt-get install -t testing libc6 \
  && apt-get install -yq --no-install-recommends $buildDeps libmagic-dev libimage-exiftool-perl \
  && mkdir /malware \
  && chown -R malice:malice /malware \
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
  && curl -Ls https://downloads.sourceforge.net/project/ssdeep/$SSDEEP/$SSDEEP.tar.gz > \
    /tmp/$SSDEEP.tar.gz \
  && cd /tmp \
  && tar zxvf $SSDEEP.tar.gz \
  && cd $SSDEEP \
  && ./configure \
  && make \
  && make install \
  && echo "Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /root/.gnupg

ENV GOLANG_VERSION 1.8.3
ENV GOLANG_DOWNLOAD_SHA256 1862f4c3d3907e59b04a757cfda0ea7aa9ef39274af99a784f5be843c80c6772

COPY . /go/src/github.com/maliceio/malice-fileinfo
RUN buildDeps='ca-certificates \
               build-essential \
               mercurial \
               git-core \
               openssl \
               gnupg \
               curl' \
  && set -x \
  && apt-get update -qq \
  && apt-get install -yq --no-install-recommends $buildDeps \
  && echo "Install Go..." \
  && cd /tmp \
  && ARCH="$(dpkg --print-architecture)" \
  && curl -Ls https://storage.googleapis.com/golang/go$GOLANG_VERSION.linux-$ARCH.tar.gz > /tmp/golang.tar.gz \
  && echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
  && tar -C /usr/local -xzf /tmp/golang.tar.gz \
  && export PATH=$PATH:/usr/local/go/bin \
  && echo "Building info Go binary..." \
  && cd /go/src/github.com/maliceio/malice-fileinfo \
  && export GOPATH=/go \
  && go version \
  && go get \
  && go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/info \
  && echo "Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /go /usr/local/go /root/.gnupg

VOLUME ["/malware"]
WORKDIR /malware

ENTRYPOINT ["gosu","malice","tini","--","info"]
CMD ["--help"]
