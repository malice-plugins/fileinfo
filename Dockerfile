FROM alpine:edge

MAINTAINER blacktop, https://github.com/blacktop

RUN buildDeps='autoconf \
              automake \
              file-dev \
              flex \
              gcc \
              git \
              jansson-dev \
              libc-dev \
              libtool \
              make \
              openssl-dev \
              python-dev' \
  && set -x \
  && apk --update add python openssl file exiftool $buildDeps \
  && apk del --purge $buildDeps \
  && rm -rf /tmp/* /root/.cache /var/cache/apk/*

# Add scan
ADD scan /malware/scan

VOLUME ["/malware"]
# VOLUME ["/rules"]

WORKDIR /malware

ENTRYPOINT ["scan"]
