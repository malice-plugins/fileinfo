FROM gliderlabs/alpine:edge

MAINTAINER blacktop, https://github.com/blacktop

ENV SSDEEP ssdeep-2.13

RUN apk-install python file exiftool
RUN apk-install -t build-deps build-base autoconf automake python-dev py-pip flex libc-dev libffi-dev libtool curl ca-certificates \
  && set -x \
  && echo "Installing ssdeep..." \
  && curl -Ls https://downloads.sourceforge.net/project/ssdeep/$SSDEEP/$SSDEEP.tar.gz > /tmp/$SSDEEP.tar.gz \
  && cd /tmp \
  && tar zxvf $SSDEEP.tar.gz \
  && cd $SSDEEP \
  && ./configure \
  && make \
  && make install \
  && pip install envoy \
  && rm -rf /tmp/* /root/.cache \
  && apk del --purge build-deps

COPY . /opt/fileinfo

VOLUME ["/malware"]

WORKDIR /malware

ENTRYPOINT ["/opt/fileinfo/scan"]
