FROM debian:wheezy

MAINTAINER blacktop, https://github.com/blacktop

# This is the release of https://github.com/hashicorp/docker-base to pull in order
# to provide HashiCorp-built versions of basic utilities like dumb-init and gosu.
ENV DOCKER_BASE_VERSION 0.0.4
ENV DBASE_URL releases.hashicorp.com/docker-base

# Create a malice user and group first so the IDs get set the same way, even as
# the rest of this may change over time.
RUN groupadd -r malice && useradd -r -g malice malice

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
  && echo "Install hashicorp/docker-base..." \
  && gpg --recv-keys 91A6E7F85D05C65630BEF18951852D87348FFC4C \
  && mkdir -p /tmp/build \
  && cd /tmp/build \
  && wget https://${DBASE_URL}/${DOCKER_BASE_VERSION}/docker-base_${DOCKER_BASE_VERSION}_linux_amd64.zip \
  && wget https://${DBASE_URL}/${DOCKER_BASE_VERSION}/docker-base_${DOCKER_BASE_VERSION}_SHA256SUMS \
  && wget https://${DBASE_URL}/${DOCKER_BASE_VERSION}/docker-base_${DOCKER_BASE_VERSION}_SHA256SUMS.sig \
  && gpg --batch --verify docker-base_${DOCKER_BASE_VERSION}_SHA256SUMS.sig docker-base_${DOCKER_BASE_VERSION}_SHA256SUMS \
  && grep ${DOCKER_BASE_VERSION}_linux_amd64.zip docker-base_${DOCKER_BASE_VERSION}_SHA256SUMS | sha256sum -c \
  && unzip docker-base_${DOCKER_BASE_VERSION}_linux_amd64.zip  \
  && cp bin/gosu bin/dumb-init /bin \
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
ENTRYPOINT ["/bin/info"]

CMD ["--help"]
