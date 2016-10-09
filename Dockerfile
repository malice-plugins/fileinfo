FROM debian:wheezy

MAINTAINER blacktop, https://github.com/blacktop

# This is the release of https://github.com/hashicorp/docker-base to pull in order
# to provide HashiCorp-built versions of basic utilities like dumb-init and gosu.
ENV DOCKER_BASE_VERSION 0.0.4
ENV DBASE_URL releases.hashicorp.com/docker-base

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

COPY . /go/src/github.com/maliceio/malice-fileinfo
RUN buildDeps='ca-certificates \
               build-essential \
               golang-go \
               mercurial \
               git-core \
               openssl \
               unzip \
               gnupg \
               curl' \
  && set -x \
  && echo 'deb http://ftp.us.debian.org/debian/ testing main contrib non-free' >> /etc/apt/sources.list \
  && apt-get update -qq \
  && apt-get install -t testing libc6 \
  && apt-get install -yq $buildDeps libimage-exiftool-perl --no-install-recommends \
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
  && curl -Ls https://downloads.sourceforge.net/project/ssdeep/$SSDEEP/$SSDEEP.tar.gz > /tmp/$SSDEEP.tar.gz \
  && cd /tmp \
  && tar zxvf $SSDEEP.tar.gz \
  && cd $SSDEEP \
  && ./configure \
  && make \
  && make install \
  && echo "Building info Go binary..." \
  && cd /go/src/github.com/maliceio/malice-fileinfo \
  && mv docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh \
  && export GOPATH=/go \
  && go get \
  && go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/info \
  && echo "Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /go /root/.gnupg

VOLUME ["/malware"]

WORKDIR /malware

# ENTRYPOINT ["docker-entrypoint.sh"]
ENTRYPOINT ["gosu","malice","tini","--","info"]

CMD ["--help"]
