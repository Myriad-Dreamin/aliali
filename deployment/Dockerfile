FROM ubuntu:20.04

RUN apt-get update
RUN apt-get install ca-certificates -y
RUN update-ca-certificates

COPY bin/notifier /usr/bin/bilibili-notifier

WORKDIR /uploader-workdir
