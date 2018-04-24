FROM ubuntu:latest
ENV WORKSPACE=/home/joly/workspace
RUN mkdir -p $WORKSPACE
COPY . $WORKSPACE
COPY config-container/config $WORKSPACE/config/
WORKDIR $WORKSPACE/bin/
