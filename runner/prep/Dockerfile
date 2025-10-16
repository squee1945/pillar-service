FROM ubuntu:22.04

RUN apt-get update \
    && apt-get install -y \
      git \
    && rm -rf /var/lib/apt/lists/*

COPY run.sh /usr/local/bin/run.sh
RUN chmod +x /usr/local/bin/run.sh

ENTRYPOINT ["/usr/local/bin/run.sh"]
