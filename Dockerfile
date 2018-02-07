
FROM ubuntu:17.10
MAINTAINER Ege Balcı <ege.balci@invictuseurope.com>


USER root

RUN apt-get update -y
RUN apt-get install -y git golang nasm mingw-w64-i686-dev mingw-w64-tools mingw-w64-x86-64-dev mingw-w64-common mingw-w64 mingw-ocaml xxd gcc-multilib g++-multilib
RUN git clone https://github.com/egebalci/Amber.git /usr/share/Amber

WORKDIR /usr/share/Amber/src
ENV GOPATH /usr/share/Amber/lib
RUN go build -o /usr/share/Amber/amber

WORKDIR /usr/share/Amber
ENV TERM xterm-256color
ENTRYPOINT ["/usr/share/Amber/amber"]
CMD ["--help"]
