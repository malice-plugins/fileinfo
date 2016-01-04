FROM debian:wheezy

MAINTAINER blacktop, https://github.com/blacktop

ENV SSDEEP ssdeep-2.13

RUN buildDeps='build-essential \
               python-dev \
               python-pip \
               adduser \
               curl' \
  && set -x \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps \
                          python \
  && set -x \
  && echo "Installing ssdeep..." \
  && curl -Ls https://downloads.sourceforge.net/project/ssdeep/$SSDEEP/$SSDEEP.tar.gz > /tmp/$SSDEEP.tar.gz \
  && cd /tmp \
  && tar zxvf $SSDEEP.tar.gz \
  && cd $SSDEEP \
  && ./configure \
  && make \
  && make install \
  && pip install envoy pydeep \
  && echo "Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY . /opt/fileinfo

VOLUME ["/malware"]

WORKDIR /malware

ENTRYPOINT ["/opt/fileinfo/scan"]
