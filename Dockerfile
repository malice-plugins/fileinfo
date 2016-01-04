FROM debian:wheezy

MAINTAINER blacktop, https://github.com/blacktop

ENV SSDEEP ssdeep-2.13

COPY . /opt/fileinfo

RUN buildDeps='build-essential \
               python-dev \
               python-pip \
               curl' \
  && set -x \
  && echo 'deb http://ftp.us.debian.org/debian/ testing main contrib non-free' >> /etc/apt/sources.list \
  && apt-get update -qq \
  && apt-get install -t testing libc6 \
  && apt-get install -yq $buildDeps \
                          libimage-exiftool-perl \
                          python --no-install-recommends \
  && set -x \
  && echo "Update TRiD Definitions..." \
  && python /opt/fileinfo/trid/tridupdate.py \
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

VOLUME ["/malware"]

WORKDIR /malware

ENTRYPOINT ["/opt/fileinfo/scan"]
