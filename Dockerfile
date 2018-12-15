
FROM ubuntu:17.10
MAINTAINER Ege Balcı <ege.balci@invictuseurope.com>
USER root
RUN apt-get update -y
RUN apt-get install -y git golang nasm fonts-powerline
RUN mkdir /root/go
ENV GOPATH /root/go
RUN go get -v github.com/egebalci/amber

ENTRYPOINT ["/root/go/bin/amber"]
CMD ["--help"]
