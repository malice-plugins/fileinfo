FROM alpine:edge

MAINTAINER blacktop, https://github.com/blacktop

RUN buildDeps='build-base \
               python-dev \
               py-pip \
               curl \
               ca-certificates' \
  && set -x \
  && apk --update add python file exiftool $buildDeps \
  && pip install --upgrade pip setuptools wheel \
  && echo "Installing ssdeep..."
  && curl -Ls https://downloads.sourceforge.net/project/ssdeep/ssdeep-2.13/ssdeep-2.13.tar.gz > /tmp/ssdeep.tar.gz && \
  && cd /tmp \
  && tar zxvf ssdeep.tar.gz \
  && cd ssdeep \
  && ./configure \
  && make \
  && make install \
  && pip install ssdeep \
                 envoy \
  && apk del --purge $buildDeps \
  && rm -rf /tmp/* /root/.cache /var/cache/apk/*

# Add exiftool stuff
ADD exif/ /opt/fileinfo/exif
# Add TRiD stuff
ADD trid/ /opt/fileinfo/trid
# Add scan
ADD scan /opt/fileinfo/scan

VOLUME ["/malware"]

WORKDIR /malware

ENTRYPOINT ["/opt/fileinfo/scan"]
